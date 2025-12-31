import { useGet } from '@/utils/request.js'

// 获取仪表板数据
export function getDashboardData(params = {}) {
  return useGet('/api/admin/dashboard', params)
}
