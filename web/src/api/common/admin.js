export function getRolesApi(params) {
  return useGet('/api/admin/roles',params)
}
export function createRoleApi(params) {
  return usePost('/api/admin/roles',params)
}
export function updateRoleApi(params) {
  return usePut('/api/admin/roles',params)
}
export function deleteRoleApi(params) {
  return useDelete('/api/admin/roles',params)
}
export function getUserPermissionsApi(params) {
  return useGet('/api/admin/permissions/user',params)
}
export function getRolePermissionsApi(params) {
  return useGet('/api/admin/permissions/role',params)
}
export function updateRolePermissionsApi(params) {
  return usePut('/api/admin/permissions/role',params)
}

export function getApiApi(params) {
  return useGet('/api/admin/apis',params)
}
export function createApiApi(params) {
  return usePost('/api/admin/apis',params)
}
export function updateApiApi(params) {
  return usePut('/api/admin/apis',params)
}
export function deleteApiApi(params) {
  return useDelete('/api/admin/apis',params)
}