import { usePost, useGet } from '~/utils/request'

// 数据导入 - Excel
export async function importExcel(file, options = {}) {
  const formData = new FormData()
  formData.append('file', file)
  
  if (options.resource) {
    formData.append('resource', options.resource)
  }
  
  if (options.mapping) {
    formData.append('mapping', JSON.stringify(options.mapping))
  }
  
  if (options.skipHeader) {
    formData.append('skip_header', options.skipHeader ? '1' : '0')
  }
  
  return usePost(`/api/admin/import/${options.resource}`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    ...options
  })
}

// 数据导入 - CSV
export async function importCSV(file, options = {}) {
  const formData = new FormData()
  formData.append('file', file)
  
  if (options.resource) {
    formData.append('resource', options.resource)
  }
  
  if (options.delimiter) {
    formData.append('delimiter', options.delimiter)
  }
  
  if (options.mapping) {
    formData.append('mapping', JSON.stringify(options.mapping))
  }
  
  if (options.skipHeader) {
    formData.append('skip_header', options.skipHeader ? '1' : '0')
  }
  
  return usePost(`/api/admin/import/${options.resource}`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    ...options
  })
}

// 数据导入 - JSON
export async function importJSON(file, options = {}) {
  const formData = new FormData()
  formData.append('file', file)
  
  if (options.resource) {
    formData.append('resource', options.resource)
  }
  
  return usePost(`/api/admin/import/${options.resource}`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    ...options
  })
}

// 获取导入模板
export async function getImportTemplate(resource, type = 'excel') {
  return useGet(`/api/admin/import/template/${resource}`, { type })
}

// 获取导入历史
export async function getImportHistory(params = {}) {
  return useGet('/api/admin/import/history', params)
}

// 获取导入任务详情
export async function getImportTask(taskId) {
  return useGet(`/api/admin/import/task/${taskId}`)
}

// 取消导入任务
export async function cancelImportTask(taskId) {
  return usePost(`/api/admin/import/task/${taskId}/cancel`)
}

// 重新执行导入任务
export async function retryImportTask(taskId) {
  return usePost(`/api/admin/import/task/${taskId}/retry`)
}

// 获取字段映射建议
export async function getFieldMapping(resource, fileHeaders) {
  return usePost('/api/admin/import/mapping', {
    resource,
    file_headers: fileHeaders
  })
}