<template>
  <div class="app-container">
    <div class="handle-box">
      <el-button type="success" icon="el-icon-plus" class="mr5" @click="handleAdd">Create</el-button>
      <el-button type="primary" icon="el-icon-search" @click="getList">Search</el-button>
    </div>
    <el-table v-loading="listLoading" :data="list" border fit highlight-current-row style="width: 100%">
      <el-table-column prop="id" align="center" label="ID"/>
      <el-table-column prop="userName" align="center" label="Account"/>
      <el-table-column prop="nickName" align="center" label="Username"/>
      <el-table-column prop="email" align="center" label="Email" width="220"/>
      <el-table-column prop="phone" align="center" label="Phone Number"/>
      <el-table-column prop="lastLogin" :formatter="dateFormat" align="center" label="Last Login Time"/>

      <el-table-column align="center" label="Actions" width="200">
        <template slot-scope="{row}">
          <el-button-group>
            <el-button type="primary" size="small" icon="el-icon-edit" class="mr5"
                       @click="handleEdit(row)">
              Edit
            </el-button>
            <el-button type="danger" size="small" icon="el-icon-delete" @click="handleDel(row.id)">
              Delete
            </el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form
        ref="dataForm"
        :rules="rules"
        :model="temp"
        label-position="left"
        label-width="90px"
        style="width: 400px; margin-left:50px;"
      >
        <el-form-item label="Account" prop="userName">
          <el-input v-model="temp.userName" :disabled="dialogStatus === 'update'"/>
        </el-form-item>
        <el-form-item label="Username" prop="nickName">
          <el-input v-model="temp.nickName"/>
        </el-form-item>
        <el-form-item label="Password" prop="password" v-if="dialogStatus === 'create'">
          <el-input v-model="temp.password" type="password"/>
        </el-form-item>
        <el-form-item label="Confirm Password" prop="confirmPassword" v-if="dialogStatus === 'create'">
          <el-input v-model="temp.confirmPassword" type="password"/>
        </el-form-item>
        <el-form-item label="Email" prop="email">
          <el-input v-model="temp.email" type="email"/>
        </el-form-item>
        <el-form-item label="Phone Number" prop="phone">
          <el-input v-model="temp.phone"/>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          Cancel
        </el-button>
        <el-button
          type="primary"
          :loading="buttonLoading"
          @click="dialogStatus === 'create' ? createData() : updateData()"
        >
          Confirm
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {fetchCreate, fetchDel, fetchList, fetchUpdate} from '@/api/user'
import {dateFormat} from "@/utils";

export default {
  name: 'UserList',
  data() {
    const validatePassword = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('This field is required'))
      } else if (value !== this.temp.password) {
        callback(new Error('The two password entries do not match'))
      } else {
        callback()
      }
    }
    return {
      listLoading: false,
      list: [],
      textMap: {
        update: 'Update',
        create: 'Create'
      },
      dialogStatus: '',
      dialogFormVisible: false,
      buttonLoading: false,
      temp: {
        userName: "",
        nickName: "",
        password: "",
        confirmPassword: '',
        email: "",
        phone: "",
      },
      rules: {
        userName: [{required: true, message: 'This field is required', trigger: 'blur'}],
        nickName: [{required: true, message: 'This field is required', trigger: 'blur'}],
        confirmPassword: [{required: true, validator: validatePassword, trigger: 'blur'}],
        password: [{required: true, message: 'This field is required', trigger: 'blur'}],
      }
    }
  },
  created() {
    this.getList()
  },
  methods: {
    resetTemp() {
      this.temp = {
        userName: "",
        nickName: "",
        password: "",
        confirmPassword: '',
        email: "",
        phone: "",
      }
    },
    dateFormat(row, column, cellValue, index) {
      return dateFormat(cellValue, 'yyyy-MM-dd hh:mm')
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
          fetchCreate(this.temp).then(() => {
            this.$notify.success('Create successful!')
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
          fetchUpdate(this.temp).then(() => {
            this.$notify.success('Update successful!')
            this.buttonLoading = false
            this.dialogFormVisible = false
            this.getList()
          }).catch(() => {
            this.buttonLoading = false
            this.dialogFormVisible = false
          })
        }
      })
    },
    handleDel(id) {
      this.$confirm('Are you sure you want to delete this user?', 'Delete', {
        type: 'warning'
      }).then(() => {
        this.listLoading = true
        fetchDel(id).then(() => {
          this.$notify.success('Delete successful!')
          this.getList()
        }).finally(() => {
          this.listLoading = false
        })
      }).catch(() => {
        this.$notify.info('Delete canceled')
      })
    },

    getList() {
      this.listLoading = true
      fetchList().then(response => {
        this.list = response.data
      }).finally(() => {
        this.listLoading = false
      })
    }
  }
}
</script>

<style scoped>

</style>
