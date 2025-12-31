import { useDelete, useGet } from '@/utils/request.js'

// 获取操作日志列表
export function getOperationLogs(params = {}) {
  return useGet('/api/admin/operation-logs', params)
}

// 删除单个操作日志
export function deleteOperationLog(id, params = {}) {
  return useDelete(`/api/admin/operation-logs/${id}`, params)
}

// 批量删除操作日志
export function deleteOperationLogs(data, params = {}) {
  return useDelete('/api/admin/operation-logs', data, params)
}
