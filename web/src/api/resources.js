import { useGet, usePost, usePut, useDelete } from '~/utils/request'

// 获取所有资源列表（支持 language 参数）
export async function getResourceList(params = {}) {
  return useGet('/api/admin/resources', params)
}

// 获取指定资源的配置
export async function getResourceConfig(slug, params = {}) {
  return useGet(`/api/admin/resources/${slug}`, params)
}

// 获取资源数据列表（分页/筛选/搜索/排序）
export async function getResourceData(slug, params = {}) {
  return useGet(`/api/admin/resource-crud/${slug}`, params)
}

// 获取资源单条记录详情
export async function getResourceRecord(slug, id, params = {}) {
  return useGet(`/api/admin/resource-crud/${slug}/${id}`, params)
}

// 创建资源记录
export async function createResourceRecord(slug, data, params = {}) {
  try {
    const response = await usePost(`/api/admin/resource-crud/${slug}`, data, params)
    return response
  } catch (error) {
    // 处理验证错误
    if (error.response && error.response.status === 422) {
      return {
        code: 422,
        message: '数据验证失败',
        errors: error.response.data.errors
      }
    }
    throw error
  }
}

// 更新资源记录
export async function updateResourceRecord(slug, id, data, params = {}) {
  try {
    const response = await usePut(`/api/admin/resource-crud/${slug}/${id}`, data, params)
    return response
  } catch (error) {
    // 处理验证错误
    if (error.response && error.response.status === 422) {
      return {
        code: 422,
        message: '数据验证失败',
        errors: error.response.data.errors
      }
    }
    throw error
  }
}

// 删除资源记录
export async function deleteResourceRecord(slug, id, params = {}) {
  return useDelete(`/api/admin/resource-crud/${slug}/${id}`, params)
}

// 批量删除资源记录
export async function deleteResourceRecords(slug, ids, params = {}) {
  return useDelete(`/api/admin/resource-crud/${slug}`, { ids }, params)
}

// 运行资源动作（支持批量）
export async function runResourceAction(slug, action, payload = {}, params = {}) {
  return usePost(`/api/admin/resource-crud/${slug}/actions/${action}`, payload, params)
}

// 全局搜索（关键字 + 可选资源 slugs + 每资源返回数限制）
export async function globalSearch(keyword, slugs = [], limit = 5) {
  return usePost('/api/admin/resources/search', { keyword, slugs, limit })
}

// 导出数据
export async function exportData(slug, params = {}) {
  return useGet(`/api/admin/export/${slug}`, params)
}