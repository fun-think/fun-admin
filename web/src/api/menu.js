export function getMenus() {
  return useGet('/api/admin/menus')
}
export function getAdminMenus() {
  return useGet('/api/admin/menus')
}
export function createMenu(params) {
  return usePost('/api/admin/menus', params)
}
export function getMenu(id) {
  return useGet(`/api/admin/menus/${id}`)
}
export function updateMenu(id, params) {
  return usePut(`/api/admin/menus/${id}`, params)
}
export function deleteMenu(id) {
  return useDelete(`/api/admin/menus/${id}`)
}
