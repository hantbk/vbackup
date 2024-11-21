<template>
  <div class="app-container">
    <div class="handle-search">
      <el-form :model="listQuery" inline @submit.native.prevent>
        <el-form-item label="Name">
          <el-input v-model="listQuery.name" placeholder="name" style="width: 150px;" class="filter-item" clearable/>
        </el-form-item>
        <el-form-item :label="'type' | i18n">
          <el-select v-model="listQuery.type" class="handle-select mr5" placeholder="Please select">
            <el-option
              v-for="(item, index) in [{code: '', name: 'All'}].concat(typeList)"
              :key="index"
              :label="item.name"
              :value="item.code"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="el-icon-search" @click="getList">Search</el-button>
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

      <el-table-column prop="createdAt" align="center" :formatter="dateFormat" label="Creation Time"/>
      <el-table-column prop="endPoint" align="left" label="Server"/>
      <el-table-column class-name="status-col" label="Storage Type" width="110">
        <template slot-scope="{row}">
          {{ formatType(row.type).name }}
        </template>
      </el-table-column>
      <el-table-column class-name="status-col" label="Compression Mode" width="110">
        <template slot-scope="{row}">
          <el-tag :type="formatCompression(row.compression).color">
            {{ formatCompression(row.compression).name }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column class-name="status-col" label="Connection Status" width="110">
        <template slot-scope="{row}">
          <el-tooltip class="item" v-if="row.errmsg" effect="dark" :content="row.errmsg" placement="bottom">
            <el-tag :type="formatStatus(row.status).color">
              {{ formatStatus(row.status).name }}
            </el-tag>
          </el-tooltip>
          <el-tag v-else :type="formatStatus(row.status).color">
            {{ formatStatus(row.status).name }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column align="center" label="Actions">
        <template slot-scope="{row}">
          <el-dropdown trigger="click" hide-on-click @command="handleCmd">
            <span class="el-dropdown-link">
              Actions<i class="el-icon-arrow-down el-icon--right"></i>
            </span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item icon="el-icon-video-camera" :command="{cmd:'restore',data:row.id}">Restore
              </el-dropdown-item>
              <el-dropdown-item icon="el-icon-setting" :command="{cmd:'oper',data:row.id}">Maintenance
              </el-dropdown-item>
              <el-dropdown-item icon="el-icon-video-camera" :command="{cmd:'snap',data:row.id}">Snapshot
              </el-dropdown-item>
              <el-dropdown-item icon="el-icon-video-camera" :command="{cmd:'edit',data:row}">Edit
              </el-dropdown-item>
              <el-dropdown-item icon="el-icon-delete" :command="{cmd:'del',data:row.id}">Delete
              </el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible" top="5vh">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="220px">

        <el-form-item label="Name" prop="name">
          <el-input v-model="temp.name" clearable />
        </el-form-item>

      <el-form-item label="Storage Type" prop="type">
        <el-select v-model="temp.type" placeholder="Please select" @change="onTypeChange">
          <el-option v-for="item in typeList" :key="item.code" :label="item.name" :value="item.code" />
        </el-select>
      </el-form-item>

      <el-form-item v-if="temp.type === 4" label="Path" prop="path">
        <el-input v-model="temp.path" disabled>
          <el-button slot="append" @click="openDirSelect()">Select</el-button>
        </el-input>
      </el-form-item>

        <el-form-item v-if="temp.type===1||temp.type===2" label="Endpoint" prop="endPoint">
          <el-input v-model="temp.endPoint" :placeholder="endPointPlaceholder" clearable/>
        </el-form-item>

        <el-form-item v-if="temp.type===1||temp.type===2" label="Region" prop="region">
          <el-input v-model="temp.region" clearable/>
        </el-form-item>
        <el-form-item v-if="temp.type===1||temp.type===2||temp.type===6" label="Bucket" prop="bucket">
          <el-input v-model="temp.bucket" clearable/>
        </el-form-item>
        <el-form-item v-if="temp.type===6" label="Access Key" prop="keyId">
          <el-input v-model="temp.keyId" clearable/>
        </el-form-item>
        <el-form-item v-if="temp.type===6" label="Secret Key" prop="secret">
          <el-input v-model="temp.secret" type="password" show-password clearable/>
        </el-form-item>
        <el-form-item v-if="temp.type===1||temp.type===2" label="AWS_ACCESS_KEY_ID" prop="keyId">
          <el-input v-model="temp.keyId" clearable/>
        </el-form-item>
        <el-form-item v-if="temp.type===1||temp.type===2" label="AWS_SECRET_ACCESS_KEY" prop="secret">
          <el-input v-model="temp.secret" type="password" show-password clearable/>
        </el-form-item>
        <el-form-item v-if="temp.type===7" label="SecretID" prop="keyId">
          <el-input v-model="temp.keyId" clearable/>
        </el-form-item>
        <el-form-item v-if="temp.type===7" label="SecretKey" prop="secret">
          <el-input v-model="temp.secret" type="password" show-password clearable/>
        </el-form-item>
        <el-form-item v-if="temp.type===5" label="Account" prop="keyId">
          <el-input v-model="temp.keyId" clearable/>
        </el-form-item>
        <el-form-item v-if="temp.type===5" label="Password" prop="secret">
          <el-input v-model="temp.secret" type="password" show-password clearable/>
        </el-form-item>
        <el-form-item v-if="dialogStatus === 'create'" label="Repository Password" prop="password">
          <el-input v-model="temp.password" show-password clearable type="password"/>
        </el-form-item>
        <el-form-item v-if="dialogStatus === 'create'" label="Confirm Password" prop="confirmPassword">
          <el-input v-model="temp.confirmPassword" show-password clearable type="password"/>
        </el-form-item>
        <el-form-item v-if="dialogStatus === 'create'" label="Compression Mode" prop="type">
          <el-select v-model="temp.compression" placeholder="Please select">
            <el-option v-for="item in compressionList" :key="item.code" :label="item.name" :value="item.code"/>
          </el-select>
        </el-form-item>
        <el-form-item label="PackSize" prop="PackSize">
          <el-input v-model="temp.packSize" clearable>
            <template slot="append">MiB</template>
          </el-input>
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
import {fetchCreate, fetchDel, fetchList, fetchUpdate} from '@/api/repository'
import {dateFormat} from '@/utils'
import {repoStatusList, repoTypeList, compressionList} from "@/consts"
import {fetchLs} from "@/api/system";

// const os = require('os');

// const username = os.userInfo().username;
// const userHome = os.platform() === 'darwin' ? `/Users/${username}` : `/home/${username}`;

// console.log('userHome', userHome);

export default {
  name: 'RepositoryList',
  data() {
    const validatePassword = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('This item is required'))
      } else if (value !== this.temp.password) {
        callback(new Error('The two passwords do not match'))
      } else {
        callback()
      }
    }
    return {
      dialogDirVisible: false,
      dirList: [],
      dirCur: userHome,
      typeList: repoTypeList,
      statusList: repoStatusList,
      compressionList: compressionList,
      list: [],
      listLoading: false,
      listQuery: {
        name: '',
        type: '',
        pageNum: 1,
        pageSize: 10,
      },
      textMap: {
        update: 'Modify Repository',
        create: 'Create Repository'
      },
      dialogStatus: '',
      dialogFormVisible: false,
      dialogVisible: false,
      dialogdata: '',
      buttonLoading: false,
      endPointPlaceholder: '',
      temp: {
        id: 0,
        name: '',
        type: 4,
        endPoint: '',
        region: '',
        bucket: '',
        keyId: '',
        secret: '',
        projectId: '',
        accountName: '',
        accountKey: '',
        accountId: '',
        password: '',
        confirmPassword: '',
        compression: 0,
        packSize: 16,
        path: '/'
      },
      rules: {
        name: [{required: true, message: 'This field is required', trigger: 'blur'}],
        type: [{required: true, message: 'This field is required', trigger: 'change'}],
        endPoint: [{required: true, message: 'This field is required', trigger: 'blur'}],
        region: [{required: false, message: 'This field is required', trigger: 'blur'}],
        bucket: [{required: true, message: 'This field is required', trigger: 'blur'}],
        keyId: [{required: true, message: 'This field is required', trigger: 'blur'}],
        secret: [{required: true, message: 'This field is required', trigger: 'blur'}],
        password: [{required: true, message: 'This field is required', trigger: 'blur'}],
        confirmPassword: [{required: true, validator: validatePassword, trigger: 'blur'}],
      }
    }
  },
  created() {
    this.getList()
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
    resetTemp() {
      this.temp = {
        id: 0,
        name: '',
        type: 4,
        endPoint: '',
        path: '/',
        region: '',
        bucket: '',
        keyId: '',
        secret: '',
        projectId: '',
        accountName: '',
        accountKey: '',
        accountId: '',
        password: '',
        confirmPassword: '',
        compression: 0,
        packSize: 16
      }
      this.endPointPlaceholder = ''
    },
    onTypeChange(val) {
      this.endPointPlaceholder = ''
      if (val === 4) {
        this.temp.path = '';
      }
      switch (val) {
        case 1:
          this.endPointPlaceholder = 'http(s)://s3host:port'
          break
        case 3:
          this.endPointPlaceholder = 'user@host:/data/my_backup_repo'
          break
        case 4:
          this.endPointPlaceholder = '/data/my_backup_repo'
          break
        case 5:
          this.endPointPlaceholder = 'http(s)://host:8000/my_backup_repo/'
          break
      }
    },
    handleCmd(datas) {
      const cmd = datas.cmd
      const data = datas.data
      switch (cmd) {
        case 'snap':
          this.$router.push('/repository/snapshot/' + data)
          break
        case 'restore':
          this.$router.push('/repository/restore/' + data)
          break
        case 'del':
          this.handleDel(data)
          break
        case 'edit':
          this.handleEdit(data)
          break
        case 'oper':
          this.$router.push('/repository/operation/' + data)
          break
      }
    },
    handleAdd() {
      this.resetTemp()
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          this.buttonLoading = true
          this.temp.packSize = Number(this.temp.packSize)
          fetchCreate(this.temp).then(() => {
            this.$notify.success('Created successfully!')
            this.buttonLoading = false
            this.dialogFormVisible = false
            this.getList()
          }).catch(() => {
            this.buttonLoading = false
          })
        }
      })
    },
    handleEdit(row) {
      this.temp = Object.assign({}, row)
      this.dialogStatus = 'update'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    updateData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          this.buttonLoading = true
          this.temp.packSize = Number(this.temp.packSize)
          fetchUpdate(this.temp).then(() => {
            this.$notify.success('Modification succeeded!')
            this.buttonLoading = false
            this.dialogFormVisible = false
            this.getList()
          }).catch(() => {
            this.buttonLoading = false
          })
        }
      })
    },
    handleDel(id) {
      this.$confirm('Are you sure you want to delete this repository? ', 'Delete', {
        type: 'warning'
      }).then(() => {
        this.listLoading = true
        fetchDel(id).then(() => {
          this.$notify.success('Delete successfully!')
          this.getList()
        }).finally(() => {
          this.listLoading = false
        })
      }).catch(() => {
        this.$notify.info('Cancel deletion')
      })
    },
    formatType(code) {
      return this.typeList.find(item => item.code === code)
    },
    formatCompression(code) {
      return this.compressionList.find(item => item.code === code)
    },
    formatStatus(code) {
      return this.statusList.find(item => item.code === code)
    },
    getList() {
      this.listLoading = true
      fetchList(this.listQuery).then(response => {
        this.list = response.data
      }).finally(() => {
        this.listLoading = false
      })
    },
    openDirSelect() {
      this.dialogDirVisible = true
      this.dirCur = this.temp.path || userHome
      this.dirList = []
      this.lsDir(this.dirCur)
    },
    confirmDirSelect(path) {
      if (path) {
        this.dirCur = path;
      }
      this.temp.path = this.dirCur;

      if (this.temp.type === 4) {
        this.temp.endPoint = this.temp.path;
      }

      if (!this.temp.name) {
        this.temp.name = this.dirCur;
      }
      this.dialogDirVisible = false;
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

.el-dropdown-link {
  cursor: pointer;
  color: #409EFF;
}

.el-icon-arrow-down {
  font-size: 12px;
}

.repo-type-tips {
  margin-left: 10px;
  color: red;
}

</style>
