export function getMenusApi() {
  return useGet('/api/admin/menus')
}
export function getAdminMenusApi() {
  return useGet('/api/admin/menus')
}
export function createMenuApi(params) {
  return usePost('/api/admin/menus',params)
}
export function updateMenuApi(params) {
  return usePut('/api/admin/menus',params)
}
export function deleteMenusApi(params) {
  return useDelete('/api/admin/menus',params)
}