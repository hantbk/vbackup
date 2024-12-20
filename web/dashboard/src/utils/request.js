import axios from 'axios'
import { Notification, Loading } from 'element-ui'
import store from '@/store'
import { getToken } from '@/utils/auth'
import router from '@/router'
import { fetchVersion } from '@/api/system'

// create an axios instance
const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_API, // url = base url + request url
  timeout: 300 * 1000 // request timeout
})

// request interceptor
service.interceptors.request.use(
  config => {
    if (!config.url.includes('/login')) {
      refreshToken()
      if (store.getters.token.token) {
        config.headers['Authorization'] = 'Bearer ' + getToken().token
      }
    }
    return config
  },
  error => {
    console.error(error) // for debug
    return Promise.reject(error)
  }
)

// Validate whether the current token has expired
const isTokenExpired = () => {
  // Get the expiration time of the token in seconds
  const expireTime = new Date(getToken().expiresAt).getTime() / 1000

  if (expireTime) {
    // Get the current time in seconds
    const nowTime = new Date().getTime() / 1000

    // If the token will expire in less than 20 minutes, fetch a new token
    return (expireTime - nowTime) < 1200 // 1200 seconds = 20 minutes
  }

  return false // Return false if the expireTime is not valid
}

let isRefreshing = false

const refreshToken = async() => {
  if (isTokenExpired() && !isRefreshing) {
    isRefreshing = true
    await store.dispatch('user/refreshToken').finally(() => {
      isRefreshing = false
    })
  }
}

let RepoLoading = 'normal'

let timei = null

let loading = null

const getLoadingText = (load) => {
  // Determine the loading text based on the load status
  switch (load) {
    case 'normal':
      return 'Normal' // Normal status
    case 'loading':
      return 'The repository is loading, please wait...' // Loading message
    case 'upgrading':
      return 'Upgrading... After the upgrade is successful, macOS needs to be restarted manually, while Linux users should wait for the automatic restart to complete.' // Upgrade message with instructions
    default:
      return 'Normal' // Default case
  }
}

const checkRepoLoading = () => {
  if (RepoLoading !== 'normal') {
    if (loading != null) {
      return
    }
    loading = Loading.service({
      lock: true,
      text: getLoadingText(RepoLoading),
      spinner: 'el-icon-loading',
      background: 'rgba(0, 0, 0, 0.7)'
    })
    timei = setInterval(() => {
      fetchVersion()
    }, 1000)
  } else {
    if (loading != null) {
      loading.close()
      location.reload()
    }
    if (timei != null) {
      clearInterval(timei)
    }
  }
}

// response interceptor
service.interceptors.response.use(
  /**
   * If you want to get http information such as headers or status
   * Please return  response => response
   */

  /**
   * Determine the request status by custom code
   * Here is just an example
   * You can also judge the status by HTTP Status Code
   */
  response => {
    const res = response.data
    if (!res.success) {
      if (res.code === 401) {
        store.dispatch('user/logout')
        router.push('/login')
      } else if (res.code === 403) {
        router.push('/403')
      } else {
        Notification({
          title: 'Error',
          message: res.message || 'Error',
          type: 'error'
        })
        return Promise.reject(new Error(res.message || 'Error'))
      }
    } else {
      RepoLoading = res.systemStatus
      checkRepoLoading()
      return res
    }
  },
  error => {
    Notification({
      title: 'Error',
      message: error.message,
      type: 'error'
    })
    return Promise.reject(error)
  }
)

export default service
