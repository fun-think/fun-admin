import { toArray } from '@v-c/utils'

export function useAccess() {
  const userStore = useUserStore()
  const roles = computed(() => userStore.roles)
  const permissions = computed(() => userStore.permissions || [])
  const hasAccess = (roles2) => {
    const accessRoles = userStore.roles
    const roleArr = toArray(roles2).flat(1)
    return roleArr.some(role => accessRoles?.includes(role))
  }
  const hasPermission = (perm) => {
    if (!perm) return true
    const list = toArray(perm).flat(1)
    return list.some(p => permissions.value?.includes(p))
  }
  return {
    hasAccess,
    hasPermission,
    roles,
    permissions,
  }
}
