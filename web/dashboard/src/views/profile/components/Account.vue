<template>
  <el-card style="padding: 20px">
    <el-form>
      <el-form-item label="Account">
        <el-input v-model.trim="tmp.userName"/>
      </el-form-item>
      <el-form-item label="Username">
        <el-input v-model.trim="tmp.nickName"/>
      </el-form-item>
      <el-form-item label="Email">
        <el-input v-model.trim="tmp.email"/>
      </el-form-item>
      <el-form-item label="Phone Number">
        <el-input v-model.trim="tmp.phone"/>
      </el-form-item>
      <p>Last Login Time: {{ tmp.lastLogin }}</p>
      <el-form-item>
        <el-button type="primary" @click="submit">Update</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script>
import {fetchUpdate} from "@/api/user";
import {setUserInfo} from "@/utils/auth";

export default {
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
    return {
      tmp: {}
    }
  },
  created() {
    this.tmp = this.user
  },
  methods: {
    submit() {
      const tmp = this.user
      fetchUpdate(tmp).then(() => {
        setUserInfo(this.tmp)
        this.$notify({
          message: 'Update successful!',
          type: 'success'
        })
      })
    }
  }
}
</script>
