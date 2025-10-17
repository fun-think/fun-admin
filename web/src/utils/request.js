import axios from 'axios'
import { AxiosLoading } from './loading.js'
import { STORAGE_AUTHORIZE_KEY } from '~/composables/authorization'
import { ContentTypeEnum, RequestEnum } from '~#/http-enum'
import router from '~/router'
import { useUserStore } from '~/stores/user'
import { useAuthorization } from '~/composables/authorization'

const instance = axios.create({
  baseURL: import.meta.env.VITE_APP_BASE_API ?? '/',
  timeout: 6e4,
  headers: { 'Content-Type': ContentTypeEnum.JSON },
})
const axiosLoading = new AxiosLoading()
async function requestHandler(config) {
  if (import.meta.env.DEV && import.meta.env.VITE_APP_BASE_API_DEV && import.meta.env.VITE_APP_BASE_URL_DEV && config.customDev)
    config.baseURL = import.meta.env.VITE_APP_BASE_API_DEV

  // 使用useAuthorization获取token
  const token = useAuthorization()
  if (token.value && config.token !== false) {
    console.log('Setting Authorization header:', token.value) // 调试信息
    config.headers.Authorization = token.value
  }
  const { locale } = useI18nLocale()
  config.headers['Accept-Language'] = locale.value ?? 'zh-CN'
  if (config.loading)
    axiosLoading.addLoading()
  return config
}
function responseHandler(response) {
  return response.data
}
function errorHandler(error) {
  const userStore = useUserStore()
  const notification = useNotification()
  if (error.response) {
    const { data, status, statusText } = error.response
    const message = data?.message || statusText || '未知错误'
    const details = data?.details || ''
    
    if (status === 401) {
      notification?.error({
        message: '认证失败',
        description: message + (details ? `: ${details}` : ''),
        duration: 3,
      })
      userStore.logout()
      router.push({
        path: '/login',
        query: {
          redirect: router.currentRoute.value.fullPath,
        },
      }).then(() => {
      })
    }
    else if (status === 403) {
      notification?.error({
        message: '权限不足',
        description: message + (details ? `: ${details}` : ''),
        duration: 3,
      })
    }
    else if (status === 404) {
      notification?.error({
        message: '资源不存在',
        description: message + (details ? `: ${details}` : ''),
        duration: 3,
      })
    }
    else if (status === 500) {
      notification?.error({
        message: '服务器错误',
        description: message + (details ? `: ${details}` : ''),
        duration: 3,
      })
    }
    else {
      notification?.error({
        message: `请求错误 ${status}`,
        description: message + (details ? `: ${details}` : ''),
        duration: 3,
      })
    }
  } else if (error.request) {
    // 请求已发出但没有收到响应
    useNotification().error({
      message: '网络错误',
      description: '无法连接到服务器，请检查网络连接',
      duration: 3,
    })
  } else {
    // 其他错误
    useNotification().error({
      message: '请求错误',
      description: error.message || '未知错误',
      duration: 3,
    })
  }
  return Promise.reject(error)
}
instance.interceptors.request.use(requestHandler)
instance.interceptors.response.use(responseHandler, errorHandler)
export default instance
function instancePromise(options) {
  const { loading } = options
  return new Promise((resolve, reject) => {
    instance.request(options).then((res) => {
      resolve(res)
    }).catch((e) => {
      reject(e)
    }).finally(() => {
      if (loading)
        axiosLoading.closeLoading()
    })
  })
}
export function useGet(url, params, config) {
  const options = {
    url,
    params,
    method: RequestEnum.GET,
    ...config,
  }
  return instancePromise(options)
}
export function usePost(url, data, config) {
  const options = {
    url,
    data,
    method: RequestEnum.POST,
    ...config,
  }
  return instancePromise(options)
}
export function usePut(url, data, config) {
  const options = {
    url,
    data,
    method: RequestEnum.PUT,
    ...config,
  }
  return instancePromise(options)
}
export function useDelete(url, data, config) {
  const options = {
    url,
    data,
    method: RequestEnum.DELETE,
    ...config,
  }
  return instancePromise(options)
}