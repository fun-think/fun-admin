import { useDelete, useGet, usePost, usePut } from '@/utils/request.js'

// 获取字典类型列表
export async function getDictionaryTypes(params = {}) {
  return useGet('/api/admin/dictionaries', params)
}

// 获取字典类型详情
export async function getDictionaryType(id) {
  return useGet(`/api/admin/dictionaries/${id}`)
}

// 创建字典类型
export async function createDictionaryType(data) {
  return usePost('/api/admin/dictionaries', data)
}

// 更新字典类型
export async function updateDictionaryType(id, data) {
  return usePut(`/api/admin/dictionaries/${id}`, data)
}

// 删除字典类型
export async function deleteDictionaryType(id) {
  return useDelete(`/api/admin/dictionaries/${id}`)
}

// 获取字典数据列表
export async function getDictionaryData(typeCode, params = {}) {
  return useGet(`/api/admin/dictionaries/${typeCode}/data`, params)
}

// 获取字典数据详情
export async function getDictionaryDataItem(typeCode, id) {
  return useGet(`/api/admin/dictionaries/${typeCode}/data/${id}`)
}

// 创建字典数据
export async function createDictionaryData(typeCode, data) {
  return usePost(`/api/admin/dictionaries/${typeCode}/data`, data)
}

// 更新字典数据
export async function updateDictionaryData(typeCode, id, data) {
  return usePut(`/api/admin/dictionaries/${typeCode}/data/${id}`, data)
}

// 删除字典数据
export async function deleteDictionaryData(typeCode, id) {
  return useDelete(`/api/admin/dictionaries/${typeCode}/data/${id}`)
}

// 批量删除字典数据
export async function batchDeleteDictionaryData(typeCode, ids) {
  return useDelete(`/api/admin/dictionaries/${typeCode}/data`, { ids })
}

// 获取字典数据（通过代码）
export async function getDictionaryByCode(code) {
  return useGet(`/api/admin/dictionaries/code/${code}`)
}

// 刷新字典缓存
export async function refreshDictionaryCache() {
  return usePost('/api/admin/dictionaries/refresh-cache')
}

// 导出字典数据
export async function exportDictionaryData(typeCode, format = 'excel') {
  return useGet(`/api/admin/dictionaries/${typeCode}/export`, { format })
}

// 导入字典数据
export async function importDictionaryData(typeCode, file) {
  const formData = new FormData()
  formData.append('file', file)

  return usePost(`/api/admin/dictionaries/${typeCode}/import`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}
