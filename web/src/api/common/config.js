import { useGet, usePost, usePut } from '~/utils/request'

// 获取系统配置
export async function getSystemConfig() {
  return useGet('/api/admin/configs/system')
}

// 更新系统配置
export async function updateSystemConfig(config) {
  return usePut('/api/admin/configs/system', config)
}

// 获取邮件配置
export async function getEmailConfig() {
  return useGet('/api/admin/config/email')
}

// 更新邮件配置
export async function updateEmailConfig(config) {
  return usePut('/api/admin/config/email', config)
}

// 测试邮件配置
export async function testEmailConfig(config) {
  return usePost('/api/admin/config/email/test', config)
}

// 发送测试邮件
export async function sendTestEmail(to, subject, content) {
  return usePost('/api/admin/config/email/send-test', {
    to,
    subject,
    content
  })
}

// 获取存储配置
export async function getStorageConfig() {
  return useGet('/api/admin/config/storage')
}

// 更新存储配置
export async function updateStorageConfig(config) {
  return usePut('/api/admin/config/storage', config)
}

// 测试存储配置
export async function testStorageConfig(config) {
  return usePost('/api/admin/config/storage/test', config)
}

// 获取应用配置
export async function getAppConfig() {
  return useGet('/api/admin/config/app')
}

// 更新应用配置
export async function updateAppConfig(config) {
  return usePut('/api/admin/config/app', config)
}

// 获取安全配置
export async function getSecurityConfig() {
  return useGet('/api/admin/config/security')
}

// 更新安全配置
export async function updateSecurityConfig(config) {
  return usePut('/api/admin/config/security', config)
}

// 重置配置
export async function resetConfig(type) {
  return usePost(`/api/admin/configs/${type}/reset`)
}

// 导出配置
export async function exportConfig() {
  return useGet('/api/admin/configs/export')
}

// 导入配置
export async function importConfig(file) {
  const formData = new FormData()
  formData.append('file', file)
  
  return usePost('/api/admin/configs/import', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}