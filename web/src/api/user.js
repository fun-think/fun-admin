export function getUsers(params) {
  return useGet('/api/admin/v1/users', params)
}
export function createUser(params) {
  return usePost('/api/admin/v1/users', params)
}
export function getUser(id) {
  return useGet(`/api/admin/v1/users/${id}`)
}
export function updateUser(id, params) {
  return usePut(`/api/admin/v1/users/${id}`, params)
}
export function deleteUser(id) {
  return useDelete(`/api/admin/v1/users/${id}`)
}
