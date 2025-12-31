import { useDelete, useGet, usePost, usePut } from '@/utils/request.js'

// 获取用户权限
export async function getUserPermissions() {
  return useGet('/api/admin/permissions')
}

// 获取角色权限
export async function getRolePermissions(role) {
  return useGet('/api/admin/permissions/role', { role })
}

// 更新角色权限
export async function updateRolePermissions(role, permissions) {
  return usePut('/api/admin/permissions/role', { role, list: permissions })
}

// 获取用户所有权限
export async function getUserAllPermissions(user) {
  return useGet('/api/admin/permissions/user', { user })
}

// 为用户添加角色
export async function addUserRole(user, role) {
  return usePost('/api/admin/permissions/user', {}, { params: { user, role } })
}

// 删除用户角色
export async function deleteUserRole(user, role) {
  return useDelete('/api/admin/permissions/user', {}, { params: { user, role } })
}

// 获取所有角色
export async function getAllRoles() {
  return useGet('/api/admin/roles/all')
}
