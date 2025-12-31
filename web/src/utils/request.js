import axios from 'axios'
import { AxiosLoading } from './loading.js'
import { STORAGE_AUTHORIZE_KEY, useAuthorization } from '@/composables/authorization'
import { ContentTypeEnum, RequestEnum } from '~#/http-enum'
import router from '@/router'

const instance = axios.create({
  baseURL: import.meta.env.VITE_APP_BASE_API ?? '/',
  timeout: 6e4,
  headers: { 'Content-Type': ContentTypeEnum.JSON },
})
const axiosLoading = new AxiosLoading()
async function requestHandler(config) {
  if (import.meta.env.DEV && import.meta.env.VITE_APP_BASE_API_DEV && import.meta.env.VITE_APP_BASE_URL_DEV && config.customDev)
    config.baseURL = import.meta.env.VITE_APP_BASE_API_DEV

  const token = useAuthorization()
  if (token.value && config.token !== false)
    config.headers.set(STORAGE_AUTHORIZE_KEY, token.value)
  const { locale } = useI18nLocale()
  config.headers.set('Accept-Language', locale.value ?? 'zh-CN')
  if (config.loading)
    axiosLoading.addLoading()
  return config
}
function responseHandler(response) {
  return response.data
}
function errorHandler(error) {
  const token = useAuthorization()
  const notification = useNotification()

  if (error.response) {
    const { data, status, statusText } = error.response
    const errorMessage = data?.message || statusText || '未知错误'

    switch (status) {
      case 401:
        notification?.error({
          message: '未授权',
          description: errorMessage,
          duration: 3,
        })
        token.value = null
        router.push({
          path: '/login',
          query: {
            redirect: router.currentRoute.value.fullPath,
          },
        }).then(() => {})
        break
      case 403:
        notification?.error({
          message: '禁止访问',
          description: errorMessage,
          duration: 3,
        })
        break
      case 404:
        notification?.error({
          message: '资源不存在',
          description: errorMessage,
          duration: 3,
        })
        break
      case 422:
        notification?.error({
          message: '请求参数错误',
          description: errorMessage,
          duration: 3,
        })
        break
      case 429:
        notification?.error({
          message: '请求过于频繁',
          description: '请稍后再试',
          duration: 3,
        })
        break
      case 500:
        notification?.error({
          message: '服务器内部错误',
          description: errorMessage,
          duration: 3,
        })
        break
      case 502:
      case 503:
      case 504:
        notification?.error({
          message: '服务暂时不可用',
          description: '请稍后再试',
          duration: 3,
        })
        break
      default:
        notification?.error({
          message: '请求失败',
          description: errorMessage,
          duration: 3,
        })
    }
  }
  else if (error.request) {
    notification?.error({
      message: '网络错误',
      description: '请检查网络连接',
      duration: 3,
    })
  }
  else {
    notification?.error({
      message: '请求配置错误',
      description: error.message,
      duration: 3,
    })
  }

  return Promise.reject(error)
}
instance.interceptors.request.use(requestHandler)
instance.interceptors.response.use(responseHandler, errorHandler)
export default instance
function instancePromise(options) {
  const { loading, retryConfig = {} } = options
  const { retries = 0, retryDelay = 1000, retryCondition } = retryConfig

  const attemptRequest = (attemptNumber) => {
    return new Promise((resolve, reject) => {
      instance.request(options).then((res) => {
        resolve(res)
      }).catch((e) => {
        // 检查是否应该重试
        const shouldRetry = attemptNumber < retries
          && (!retryCondition || retryCondition(e))

        if (shouldRetry) {
          setTimeout(() => {
            attemptRequest(attemptNumber + 1).then(resolve).catch(reject)
          }, retryDelay * 2 ** attemptNumber) // 指数退避
        }
        else {
          reject(e)
        }
      }).finally(() => {
        if (loading && attemptNumber === retries) {
          axiosLoading.closeLoading()
        }
      })
    })
  }

  return attemptRequest(0)
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
