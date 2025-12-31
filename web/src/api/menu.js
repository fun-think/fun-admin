export function getMenusApi() {
  return useGet('/api/admin/menus')
}
export function getAdminMenusApi() {
  return useGet('/api/admin/menus')
}
export function createMenuApi(params) {
  return usePost('/api/admin/menus',params)
}
export function getMenuApi(id) {
  return useGet(`/api/admin/menus/${id}`)
}
export function updateMenuApi(id, params) {
  return usePut(`/api/admin/menus/${id}`, params)
}
export function deleteMenuApi(id) {
  return useDelete(`/api/admin/menus/${id}`)
}