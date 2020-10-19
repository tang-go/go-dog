import axios from 'axios'
 import store from '@/store'
import storage from 'store'
import router from '@/router'
// import notification from 'ant-design-vue/es/notification'
import { VueAxios } from './axios'
import { ACCESS_TOKEN } from '@/store/mutation-types'

// 创建 axios 实例
const request = axios.create({
  // API 请求的默认前缀
  baseURL: process.env.VUE_APP_API_BASE_URL,
  timeout: 6000 // 请求超时时间
})

// 异常拦截处理器
const errorHandler = (error) => {
  console.log(error)
  return Promise.reject(error)
}

// request interceptor
request.interceptors.request.use(config => {
  const token = storage.get(ACCESS_TOKEN)
  // 如果 token 存在
  // 让每个请求携带自定义 token 请根据实际情况自行修改
  if (token) {
    config.headers['Token'] = token
  }
  config.headers['TimeOut'] = 6
  config.headers['IsTest'] = false
  config.headers['TraceID'] = '11321231'

  return config
}, errorHandler)

// response interceptor
request.interceptors.response.use((response) => {
  // token 失效 跳转登录页面
  if (response.data.code === 10004) {
    store.dispatch('Logout').then(() => {
      router.push({ name: 'login' })
    })
    return response.data
  }
  // 解析请求
  return response.data
}, errorHandler)

const installer = {
  vm: {},
  install (Vue) {
    Vue.use(VueAxios, request)
  }
}

export default request

export {
  installer as VueAxios,
  request as axios
}
