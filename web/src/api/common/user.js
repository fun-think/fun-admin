export function getUserInfoApi() {
  return useGet('/api/admin/profile')
}
export function getUsersApi(params) {
  return useGet('/api/admin/users',params)
}
export function createUserApi(params) {
  return usePost('/api/admin/users',params)
}
export function updateUserApi(params) {
  return usePut('/api/admin/users',params)
}
export function deleteUserApi(params) {
  return useDelete('/api/admin/users',params)
}