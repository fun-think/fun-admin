import { useDelete, useGet, usePost, usePut } from '@/utils/request.js'

// 获取API列表
export async function getApis(params = {}) {
  return useGet('/api/admin/apis', params)
}

// 创建API
export async function createApi(data) {
  return usePost('/api/admin/apis', data)
}

// 获取API详情
export async function getApi(id) {
  return useGet(`/api/admin/apis/${id}`)
}

// 更新API
export async function updateApi(id, data) {
  return usePut(`/api/admin/apis/${id}`, data)
}

// 删除API
export async function deleteApi(id) {
  return useDelete(`/api/admin/apis/${id}`)
}
