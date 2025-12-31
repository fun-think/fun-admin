import { usePost } from '~/utils/request'

export function loginApi(params) {
  return usePost('/api/admin/login', params, {
    // 设置为false的时候不会携带token
    token: false,
    // 开发模式下使用自定义的接口
    customDev: false,
    // 是否开启全局请求loading
    loading: true,
  })
}

export function sendSMSCodeApi(params) {
  return usePost('/api/send-sms-code', params, {
    // 设置为false的时候不会携带token
    token: false,
    // 开发模式下使用自自定义的接口
    customDev: false,
    // 是否开启全局请求loading
    loading: true,
  })
}