<template>
  <el-card>
    <div slot="header" class="clearfix">
      <span>Personal Center</span>
    </div>

    <div class="user-profile">
      <div class="box-center">
        <pan-thumb :height="'100px'" :width="'100px'" :hoverable="false">
          <div>Hello</div>
          {{ user.nickName }}
        </pan-thumb>
      </div>
      <div class="box-center">
        <div class="user-name text-center">{{ user.userName }}</div>
      </div>
      <div class="box-center">
        <el-button type="primary" @click="repwdHandler">Change Password</el-button>
      </div>
      <div class="box-center">
        <el-switch
          v-model="this.user.mfa"
          active-text="Two-factor Authentication" @change="mfaChangeHandler">
        </el-switch>
      </div>
      <div class="box-center" v-if="mfa">
        <p>Please use the <el-link type="success" target="_blank" href="">OTP application</el-link> to scan the QR code below to obtain a 6-digit verification code.
        </p>
        <img @click="getQrcode" v-if="mfaQrcode" :src="mfaQrcode" class="qrcode" alt="qrcode">
        <p class="secret">Key: {{ otpInfo.secret }}</p>
        <el-input placeholder="Please enter the verification code" v-model="otpInfo.code" class="input-with-select mfacode">
          <el-button slot="append" @click="bindOtp">Bind</el-button>
        </el-input>
      </div>
    </div>
    <el-dialog
      title="Change Password"
      :visible.sync="dialogFormVisible"
    >
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="150px">
        <el-form-item label="Old Password" prop="oldPassword">
          <el-input v-model="temp.oldPassword" type="password" clearable/>
        </el-form-item>
        <el-form-item label="New Password" prop="password">
          <el-input v-model="temp.password" type="password" clearable/>
        </el-form-item>
        <el-form-item label="Confirm Password" prop="confirmPassword">
          <el-input v-model="temp.confirmPassword" type="password" clearable/>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          Cancel
        </el-button>
        <el-button
          type="primary"
          @click="repwd"
        >
          Confirm
        </el-button>
      </div>
    </el-dialog>
  </el-card>
</template>

<script>
import PanThumb from '@/components/PanThumb'
import {fetchBindOtp, fetchDel, fetchDeleteOtp, fetchOtp, fetchRePwd} from "@/api/user";
import {setUserInfo} from "@/utils/auth";

export default {
  components: {PanThumb},
  props: {
    user: {
      type: Object,
      default: () => {
        return {
          id: '',
          userName: '',
          nickName: '',
          email: '',
          phone: '',
          lastLogin: '',
          mfa: false
        }
      }
    }
  },
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
      dialogFormVisible: false,
      mfa: false,
      mfaQrcode: '',
      otpInfo: {
        secret: '',
        code: '',
        interval: 0
      },
      temp: {
        oldPassword: "",
        password: "",
        confirmPassword: '',
      },
      rules: {
        oldPassword: [{required: true, message: 'This field is required', trigger: 'blur'}],
        confirmPassword: [{required: true, validator: validatePassword, trigger: 'blur'}],
        password: [{required: true, message: 'This field is required', trigger: 'blur'}],
      }
    }
  },
  methods: {
    resetTemp() {
      this.temp = {
        oldPassword: "",
        password: "",
        confirmPassword: '',
      }
    },
    repwdHandler() {
      this.resetTemp()
      this.dialogFormVisible = true
    },
    repwd() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          fetchRePwd(this.temp).then(res => {
            this.$notify.success({
              title: 'Notice',
              message: res.data
            })
          }).finally(() => {
            this.dialogFormVisible = false
          })
        }
      })
    },
    mfaChangeHandler(value) {
      if (value) {
        this.mfa = true
        this.getQrcode()
      } else {
        this.$confirm('Are you sure you want to turn off two-factor authentication?', 'Disable Authentication', {
          type: 'warning'
        }).then(() => {
          this.otpInfo = {
            secret: '',
            code: '',
            interval: 0
          }
          fetchDeleteOtp().then(res => {
            this.user.mfa = false
            setUserInfo(this.user)
            this.mfa = false
            this.$notify.success({
              title: 'Notice',
              message: res.data
            })
          })
        })
      }
    },
    getQrcode() {
      fetchOtp().then(res => {
        const data = res.data
        this.otpInfo.secret = data.secret
        this.otpInfo.interval = data.interval
        this.mfaQrcode = data.qrImg
      })
    },
    bindOtp() {
      if (this.otpInfo.code === '') {
        this.$notify.error({
          title: 'Error',
          message: 'Verification code cannot be empty'
        })
      }
      fetchBindOtp(this.otpInfo).then(res => {
        this.$notify.success({
          title: 'Notice',
          message: res.data
        })
        this.user.mfa = true
        setUserInfo(this.user)
        this.mfa = false
      })
    }
  }
}
</script>

<style lang="scss" scoped>
@import "../dashboard/src/styles/variables";

.box-center {
  margin: 0 auto;
  display: table;
}

.text-muted {
  color: #777;
}

.user-profile {
  .user-name {
    font-weight: bold;
  }

  .box-center {
    padding-top: 10px;
  }

  .qrcode {
    display: block;
    margin: auto;
    width: 200px;
    height: 200px;
  }

  .secret {
    display: block;
    margin: auto;
    font-size: 15px;
  }

  .mfacode {
    margin-top: 10px;
  }

  .user-role {
    padding-top: 10px;
    font-weight: 400;
    font-size: 14px;
  }

  .box-social {
    padding-top: 30px;

    .el-table {
      border-top: 1px solid #dfe6ec;
    }
  }

  .user-follow {
    padding-top: 20px;
  }
}

.user-bio {
  margin-top: 20px;
  color: #606266;

  span {
    padding-left: 4px;
  }

  .user-bio-section {
    font-size: 14px;
    padding: 15px 0;

    .user-bio-section-header {
      border-bottom: 1px solid #dfe6ec;
      padding-bottom: 10px;
      margin-bottom: 10px;
      font-weight: bold;
    }
  }
}
</style>
