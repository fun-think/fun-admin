import { useGet, usePost, usePut, useDelete } from '~/utils/request'

// 获取字典类型列表
export async function getDictTypes(params = {}) {
  return useGet('/api/admin/dictionaries', params)
}

// 获取字典类型详情
export async function getDictType(id) {
  return useGet(`/api/admin/dictionaries/${id}`)
}

// 创建字典类型
export async function createDictType(data) {
  return usePost('/api/admin/dictionaries', data)
}

// 更新字典类型
export async function updateDictType(id, data) {
  return usePut(`/api/admin/dictionaries/${id}`, data)
}

// 删除字典类型
export async function deleteDictType(id) {
  return useDelete(`/api/admin/dictionaries/${id}`)
}

// 获取字典数据列表
export async function getDictData(typeCode, params = {}) {
  return useGet(`/api/admin/dict/data/${typeCode}`, params)
}

// 获取字典数据详情
export async function getDictDataItem(typeCode, id) {
  return useGet(`/api/admin/dict/data/${typeCode}/${id}`)
}

// 创建字典数据
export async function createDictData(typeCode, data) {
  return usePost(`/api/admin/dict/data/${typeCode}`, data)
}

// 更新字典数据
export async function updateDictData(typeCode, id, data) {
  return usePut(`/api/admin/dict/data/${typeCode}/${id}`, data)
}

// 删除字典数据
export async function deleteDictData(typeCode, id) {
  return useDelete(`/api/admin/dict/data/${typeCode}/${id}`)
}

// 批量删除字典数据
export async function batchDeleteDictData(typeCode, ids) {
  return useDelete(`/api/admin/dict/data/${typeCode}`, { ids })
}

// 获取字典数据（通过代码）
export async function getDictByCode(code) {
  return useGet(`/api/admin/dict/code/${code}`)
}

// 刷新字典缓存
export async function refreshDictCache() {
  return usePost('/api/admin/dict/refresh-cache')
}

// 导出字典数据
export async function exportDictData(typeCode, format = 'excel') {
  return useGet(`/api/admin/dict/export/${typeCode}`, { format })
}

// 导入字典数据
export async function importDictData(typeCode, file) {
  const formData = new FormData()
  formData.append('file', file)
  
  return usePost(`/api/admin/dict/import/${typeCode}`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}