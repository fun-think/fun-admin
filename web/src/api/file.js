import { useDelete, useGet, usePost } from '@/utils/request.js'

// 文件上传
export async function uploadFile(file, options = {}) {
  const formData = new FormData()
  formData.append('file', file)

  if (options.storageType) {
    formData.append('storage_type', options.storageType)
  }

  if (options.path) {
    formData.append('path', options.path)
  }

  return usePost('/api/admin/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    ...options,
  })
}

// 获取文件信息
export async function getFileInfo(fileId) {
  return useGet(`/api/admin/files/${fileId}`)
}

// 删除文件
export async function deleteFile(fileId) {
  return useDelete(`/api/admin/files/${fileId}`)
}

// 获取文件列表
export async function getFileList(params = {}) {
  return useGet('/api/admin/files', params)
}

// 获取存储配置
export async function getStorageConfig() {
  return useGet('/api/admin/configs/storage')
}

// 更新存储配置
export async function updateStorageConfig(config) {
  return usePut('/api/admin/configs/storage', config)
}

// 测试存储连接
export async function testStorageConnection(config) {
  return usePost('/api/admin/config/storage/test', config)
}

// 获取支持的存储类型
export async function getStorageTypes() {
  return useGet('/api/admin/config/storage/types')
}
