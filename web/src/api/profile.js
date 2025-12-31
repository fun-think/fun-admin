import { useGet, usePut } from '@/utils/request.js'

// 获取当前用户个人资料
export function getProfile(params = {}) {
  return useGet('/api/admin/profile', params)
}

// 更新个人资料
export function updateProfile(data, params = {}) {
  return usePut('/api/admin/profile', data, params)
}

// 更新密码
export function updatePassword(data, params = {}) {
  return usePut('/api/admin/profile/password', data, params)
}
