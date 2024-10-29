<template>
  <div class="app-container">
    <div class="handle-search">
      <el-form :model="listQuery" inline @submit.native.prevent>
        <el-form-item label="Name">
          <el-input v-model="listQuery.name" placeholder="name" style="width: 150px;" class="filter-item" clearable/>
        </el-form-item>
        <el-form-item :label="'path' | i18n">
          <el-input v-model="listQuery.path" placeholder="path" style="width: 150px;" class="filter-item" clearable/>
        </el-form-item>
        <el-form-item label="Repository">
          <el-select v-model="listQuery.repositoryId" class="handle-select mr5" clearable placeholder="Please select">
            <el-option
              v-for="(item, index) in [{id: 0, name: 'All'}].concat(repositoryList)"
              :key="index"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="Status">
          <el-select v-model="listQuery.status" class="handle-select mr5" clearable placeholder="Please select">
            <el-option
              v-for="(item, index) in [{status: 0, name: 'All'}].concat(status)"
              :key="index"
              :label="item.name"
              :value="item.status"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="el-icon-search" @click="handleFilter">Search</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="handle-box">
      <el-button type="primary" icon="el-icon-plus" class="mr5" @click="handleAdd">Create</el-button>
    </div>
    <el-table v-loading="listLoading" :data="list" border fit highlight-current-row style="width: 100%" empty-text="No data available">
      <el-table-column align="center" :label="'ID' | i18n" width="80">
        <template slot-scope="scope">
          <span>{{ scope.row.id }}</span>
        </template>
      </el-table-column>

      <el-table-column prop="name" align="center" label="Name"/>

      <el-table-column prop="path" align="center" :label="'path' | i18n"/>

      <el-table-column prop="repositoryId" :formatter="filterRepo" align="center" :label="'repository' | i18n"/>

      <el-table-column prop="execTimeCron" align="center" label="Execution Time">
        <template slot-scope="{row}">
          <span>{{ row.execTimeCron + '   ' }}</span>
          <el-button circle type="text" icon="el-icon-view" @click="cronNext(row.execTimeCron)"/>
        </template>
      </el-table-column>

      <el-table-column class-name="status-col" label="Status">
        <template slot-scope="{row}">
          <el-tag :type="row.status === 1 ? 'success' : 'warning'">
            {{ row.status === 1 ? 'Running' : 'Stopped' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="createdAt" align="center" :label="'createdAt' | i18n" :formatter="dateFormat"/>
      <el-table-column align="center" label="Actions" width="200">
        <template slot-scope="{row}">
          <el-button-group>
            <el-button type="success" size="small" @click="backupHandler(row.id)"
                       v-loading.fullscreen.lock="fullscreenLoading"
                       element-loading-text="Executing, please do not close the page..."
                       element-loading-spinner="el-icon-loading">
              Run Now
            </el-button>
            <el-button type="primary" size="small" icon="el-icon-edit" class="mr5" @click="handleEdit(row)"/>
            <el-button type="danger" size="small" icon="el-icon-delete" @click="handleDel(row.id)"/>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>
    <pagination
      v-show="total > 0"
      :total="total"
      :page.sync="listQuery.pageNum"
      :limit.sync="listQuery.pageSize"
      @pagination="getList"
    />

    <el-dialog
      v-el-drag-dialog
      :title="textMap[dialogStatus]"
      :visible.sync="dialogFormVisible"
      @dragDialog="handleDrag"
    >
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="120px">
        <el-form-item label="Name" prop="name">
          <el-input v-model="temp.name"/>
        </el-form-item>
        <el-form-item :label="'path' | i18n" prop="path">
          <el-input v-model="temp.path" disabled>
            <el-button slot="append" @click="openDirSelect()">Select</el-button>
          </el-input>
        </el-form-item>
        <el-form-item :label="'repository' | i18n" prop="repositoryId">
          <el-select v-model="temp.repositoryId" placeholder="Please select">
            <el-option v-for="item in repositoryList" :key="item.id" :label="item.name" :value="item.id"/>
          </el-select>
        </el-form-item>
        <el-form-item label="Status" prop="status">
          <el-select v-model="temp.status" placeholder="Please select">
            <el-option v-for="item in status" :key="item.status" :label="item.name" :value="item.status"/>
          </el-select>
        </el-form-item>
        <el-form-item label="Cron Expression" prop="execTimeCron">
          <el-popover v-model="cronPopover">
            <cron @change="changeCron" @close="cronPopover=false"/>
            <el-input
              slot="reference"
              v-model="temp.execTimeCron"
              placeholder="Please enter scheduling strategy"
              clearable
              @click="cronPopover=true"
            />
          </el-popover>
          <el-button type="text" @click="cronNext(temp.execTimeCron)">Next Execution Time</el-button>
          <span style="margin-left: 20px;color: red">Note: The last position 'year' only supports every year (*), other values are invalid</span>
        </el-form-item>
        <el-form-item label="Read Concurrency" prop="ReadConcurrency">
          <el-input v-model="temp.readConcurrency" clearable>
            <template slot="append">Default 2</template>
          </el-input>
        </el-form-item>

        <el-form-item v-if="dialogStatus === 'create'" label="Backup Now">
          <el-switch v-model="temp.immediate"/>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          Cancel
        </el-button>
        <el-button
          type="primary"
          :loading="buttonLoading"
          @click=" dialogStatus === 'create' ? createData() : updateData()"
        >
          Confirm
        </el-button>
      </div>
    </el-dialog>
    <el-dialog
      title="Next Execution Time"
      :visible.sync="dialogVisible"
      width="20%"
    >
      <div class="nexttime">
        <p v-for="(item, index) in dialogdata" :key="index">{{ item }}</p>
      </div>

      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="dialogVisible = false">Confirm</el-button>
      </span>
    </el-dialog>
    <el-dialog
      title="Select Directory/Folder"
      :visible.sync="dialogDirVisible"
    >
      <div>
        <el-breadcrumb separator="/">
          <el-breadcrumb-item class="breadcrumb-item" v-for="(item, index) in getDirSpea()" :key="index">
            <span class="title" @click="lsDir(item.path,true)">{{ item.name }}</span>
          </el-breadcrumb-item>
        </el-breadcrumb>
        <div class="filenodes">
          <span class="custom-tree-node filenode" v-for="(item, index) in filteredDirList" :key="index"
                :class="{active : dirCur === item.path}"
                @dblclick.prevent="lsDir(item.path,item.isDir)"
                @click="selectDir(item.path,item.isDir)">
          <div class="file-title">
            <i v-if="item.isDir" class="el-icon-folder"/>
            <span style="margin-left: 5px;user-select: none;">{{ item.name }}</span>
          </div>
          <span>
            <el-button
              type="text"
              class="confirmbtn"
              size="mini"
              @click="confirmDirSelect(item.path)">
                    Confirm
                  </el-button>
          </span>
        </span>
        </div>
      </div>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogDirVisible = false">
          Cancel
        </el-button>
        <el-button
          type="primary"
          @click="confirmDirSelect()">
          Confirm
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {fetchCreate, fetchDel, fetchList, fetchNextTime, fetchUpdate} from '@/api/plan'
import {fetchList as repolist} from '@/api/repository'
import Pagination from '@/components/Pagination'
import {dateFormat} from '@/utils'
import elDragDialog from '@/directive/el-drag-dialog'
import {cron} from 'vue-cron'
import {fetchBackup} from '@/api/task'
import {fetchLs} from "@/api/system";

export default {
  name: 'Plan',
  directives: {elDragDialog},
  components: {Pagination, cron},
  data() {
    return {
      status: [
        {name: 'Running', status: 1},
        {name: 'Stopped', status: 2}
      ],
      dialogDirVisible: false,
      dirCur: "/",
      dirList: [],
      fullscreenLoading: false,
      repositoryList: [],
      list: [],
      total: 0,
      listLoading: false,
      listQuery: {
        name: '',
        type: '',
        repositoryId: 0,
        status: 0,
        pageNum: 1,
        pageSize: 10
      },
      textMap: {
        update: 'Update',
        create: 'Create'
      },
      cronPopover: false,
      dialogStatus: '',
      dialogFormVisible: false,
      buttonLoading: false,
      dialogVisible: false,
      dialogdata: '',
      temp: {
        name: '',
        path: '/',
        repositoryId: '',
        status: 2,
        immediate: false,
        execTimeCron: '',
        readConcurrency: 2
      },
      rules: {
        name: [{required: true, message: 'This field is required', trigger: 'blur'}],
        status: [{required: true, message: 'Please select a type', trigger: 'change'}],
        path: [{required: true, message: 'This field is required', trigger: 'blur'}],
        execTimeCron: [{required: true, message: 'This field is required', trigger: 'blur'}],
        repositoryId: [{required: true, message: 'This field is required', trigger: 'change'}]
      }
    }
  },
  created() {
    this.getList()
    repolist().then(response => {
      this.repositoryList = response.data
    })
  },
  computed: {
    filteredDirList() {
      return this.dirList.filter(item => item.isDir);
    }
  },
  methods: {
    dateFormat(row, column, cellValue, index) {
      return dateFormat(cellValue, 'yyyy-MM-dd hh:mm')
    },
    handleDrag() {
      this.$refs.select.blur()
    },
    changeCron(val) {
      this.temp.execTimeCron = val
    },
    openDirSelect() {
      this.dialogDirVisible = true
      this.dirCur = this.temp.path
      this.dirList = []
      this.lsDir(this.dirCur, true)
    },
    confirmDirSelect(path) {
      if (path) {
        this.dirCur = path
      }
      this.temp.path = this.dirCur
      if (!this.temp.name) {
        this.temp.name = this.dirCur
      }
      this.dialogDirVisible = false
    },
    getDirSpea() {
      var dirs = this.dirCur.split("/")
      dirs.shift()
      var res = []
      var path = ''
      dirs.forEach(n => {
        if (n === '') {
          return
        }
        path = path + '/' + n
        res.push({
          name: n,
          path: path
        })
      })
      res.unshift({
        name: 'Root',
        path: '/'
      })
      return res
    },
    selectDir(path, isdir) {
      if (!isdir) {
        return
      }
      this.dirCur = path
    },
    lsDir(path, isdir) {
      if (!isdir) {
        return
      }
      this.dirCur = path
      var q = {
        path: this.dirCur
      }
      fetchLs(q).then(res => {
        this.dirList = res.data
      })
    },
    cronNext(str) {
      this.dialogdata = []
      var q = {
        cron: str
      }
      fetchNextTime(q).then(res => {
        this.dialogVisible = true
        this.dialogdata = res.data
      })
    },
    filterRepo(row, column, cellValue, index) {
      let res = 'Repository deleted'
      this.repositoryList.forEach(value => {
        if (value.id === cellValue) {
          res = value.name
          return res
        }
      })
      return res
    },
    getList() {
      this.listLoading = true
      fetchList(this.listQuery).then(response => {
        this.list = response.data.items
        this.total = response.data.total
      }).finally(() => {
        this.listLoading = false
      })
    },
    handleFilter() {
      this.listQuery.pageNum = 1
      this.getList()
    },
    resetTemp() {
      this.temp = {
        name: '',
        path: '/',
        repositoryId: '',
        status: 2,
        immediate: false,
        execTimeCron: ''
      }
    },
    handleEdit(row) {
      this.temp = Object.assign({}, row)
      this.dialogStatus = 'update'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    handleAdd() {
      this.resetTemp()
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    backupHandler(planid) {
      this.fullscreenLoading = true
      fetchBackup(planid).then(() => {
        this.$notify.success({
          title: 'Backing up...',
          dangerouslyUseHTMLString: true,
          message: 'Please go to "<a style="color: #409EFF" href="/Task/index">Task Records</a>" to view'
        })
      }).finally(() => {
        this.fullscreenLoading = false
      })
    },
    handleDel(id) {
      this.$confirm('Are you sure you want to delete this plan?', 'Delete', {
        type: 'warning'
      }).then(() => {
        this.listLoading = true
        fetchDel(id).then(() => {
          this.$notify.success('Deleted successfully!')
          this.getList()
        }).finally(() => {
          this.listLoading = false
        })
      }).catch(() => {
        this.$notify.info('Delete canceled')
      })
    },
    updateData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          if (!this.temp.execTimeCron) {
            this.$notify.error('Please enter a cron expression for scheduled backup')
            return
          }
          this.buttonLoading = true
          if (this.temp.readConcurrency === '') {
            this.temp.readConcurrency = 2
          } else {
            this.temp.readConcurrency = Number(this.temp.readConcurrency)
          }
          fetchUpdate(this.temp).then(() => {
            this.$notify.success('Updated successfully!')
            this.buttonLoading = false
            this.dialogFormVisible = false
            this.getList()
          }).catch(() => {
            this.buttonLoading = false
          })
        }
      })
    },
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          if (!this.temp.execTimeCron) {
            this.$notify.error('Please enter a cron expression for scheduled backup')
            return
          }
          this.buttonLoading = true
          if (this.temp.readConcurrency === '') {
            this.temp.readConcurrency = 2
          } else {
            this.temp.readConcurrency = Number(this.temp.readConcurrency)
          }
          fetchCreate(this.temp).then(res => {
            var planid = res.data
            this.$notify.success('Created successfully!')
            this.getList()
            if (this.temp.immediate) {
              this.backupHandler(planid)
            }
          }).finally(() => {
            this.buttonLoading = false
            this.dialogFormVisible = false
          })
        }
      })
    }
  }
}
</script>

<style lang="scss" scoped>
@import "src/styles/variables";

.custom-tree-node {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 5px;

  .file-title {
    font-size: 14px;
    padding-right: 8px;
  }
}

.active {
  background: $light-blue;
  color: $menuHover;
}


.filenodes {
  margin-top: 10px;

  .confirmbtn {
    color: $menuHover;
  }
}

.filenode:hover {
  background-color: $light-blue;
  cursor: pointer;
  color: $menuHover;
}

.breadcrumb-item {
  padding: 5px;
}

.breadcrumb-item:hover {
  background-color: $light-blue;
  cursor: pointer;
  color: $menuHover;

  .title {
    color: $menuHover;
  }
}

.nexttime {
  text-align: center;
}

.bottom .value {
  margin-right: 20px !important;
}
</style>
