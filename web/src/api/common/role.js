import { useGet, usePost, usePut, useDelete } from '~/utils/request'

// 获取角色列表
export async function getRoles(params = {}) {
  return useGet('/api/admin/roles', params)
}

// 创建角色
export async function createRole(data) {
  return usePost('/api/admin/roles', data)
}

// 获取角色详情 (注意：后端未提供单独获取角色详情的接口，这里暂时使用获取列表接口)
export async function getRole(id) {
  return useGet(`/api/admin/roles/${id}`)
}

// 更新角色
export async function updateRole(id, data) {
  return usePut(`/api/admin/roles/${id}`, data)
}

// 删除角色
export async function deleteRole(id) {
  return useDelete(`/api/admin/roles/${id}`)
}

// 获取所有角色
export async function getAllRoles() {
  return useGet('/api/admin/roles/all')
}