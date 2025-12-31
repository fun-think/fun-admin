export function getUserInfoApi() {
  return useGet('/api/admin/profile')
}
export function getUsersApi(params) {
  return useGet('/api/admin/users',params)
}
export function createUserApi(params) {
  return usePost('/api/admin/users',params)
}
export function getUserApi(id) {
  return useGet(`/api/admin/users/${id}`)
}
export function updateUserApi(id, params) {
  return usePut(`/api/admin/users/${id}`, params)
}
export function deleteUserApi(id) {
  return useDelete(`/api/admin/users/${id}`)
}