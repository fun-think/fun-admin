import { computed, reactive } from 'vue'
import { useUserStore } from '@/stores/user'

/**
 * 权限类型定义
 */
export const PERMISSION_TYPES = {
  VIEW: 'view', // 查看权限
  CREATE: 'create', // 创建权限
  UPDATE: 'update', // 更新权限
  DELETE: 'delete', // 删除权限
  EXPORT: 'export', // 导出权限
  IMPORT: 'import', // 导入权限
  AUDIT: 'audit', // 审核权限
  MANAGE: 'manage', // 管理权限
}

/**
 * 路由权限状态
 */
const permissionState = reactive({
  permissions: [], // 用户权限列表
  roles: [], // 用户角色列表
  resources: {}, // 资源权限映射
  loading: false, // 加载状态
  error: null, // 错误信息
})

/**
 * 路由权限组合式API
 */
export function useRoutePermission() {
  const userStore = useUserStore()

  /**
   * 初始化权限数据
   */
  const initPermissions = async () => {
    permissionState.loading = true
    permissionState.error = null

    try {
      // 从用户store获取权限数据
      await userStore.fetchPermissions()

      permissionState.permissions = userStore.permissions || []
      permissionState.roles = userStore.roles || []
      permissionState.resources = userStore.resourcePermissions || {}
    }
    catch (error) {
      console.error('Failed to init permissions:', error)
      permissionState.error = error.message
    }
    finally {
      permissionState.loading = false
    }
  }

  /**
   * 检查是否有指定权限
   * @param {string | Array} permission - 权限代码
   * @returns {boolean} 是否有权限
   */
  const hasPermission = (permission) => {
    if (!permission)
      return true

    if (Array.isArray(permission)) {
      return permission.some(p => permissionState.permissions.includes(p))
    }

    return permissionState.permissions.includes(permission)
  }

  /**
   * 检查是否有指定角色
   * @param {string | Array} role - 角色代码
   * @returns {boolean} 是否有角色
   */
  const hasRole = (role) => {
    if (!role)
      return true

    if (Array.isArray(role)) {
      return role.some(r => permissionState.roles.includes(r))
    }

    return permissionState.roles.includes(role)
  }

  /**
   * 检查资源权限
   * @param {string} resource - 资源名称
   * @param {string} action - 操作类型
   * @returns {boolean} 是否有权限
   */
  const hasResourcePermission = (resource, action) => {
    const resourcePerms = permissionState.resources[resource]
    if (!resourcePerms)
      return false

    return resourcePerms.includes(action) || resourcePerms.includes('*')
  }

  /**
   * 检查路由权限
   * @param {object} route - 路由对象
   * @returns {boolean} 是否可以访问
   */
  const canAccessRoute = (route) => {
    const meta = route.meta || {}

    // 检查权限
    if (meta.permissions && !hasPermission(meta.permissions)) {
      return false
    }

    // 检查角色
    if (meta.roles && !hasRole(meta.roles)) {
      return false
    }

    // 检查资源权限
    if (meta.resource && meta.action) {
      return hasResourcePermission(meta.resource, meta.action)
    }

    return true
  }

  /**
   * 过滤有权限的路由
   * @param {Array} routes - 路由数组
   * @returns {Array} 过滤后的路由
   */
  const filterRoutesByPermission = (routes) => {
    return routes.filter((route) => {
      // 检查路由权限
      if (!canAccessRoute(route)) {
        return false
      }

      // 递归过滤子路由
      if (route.children && route.children.length > 0) {
        route.children = filterRoutesByPermission(route.children)
      }

      return true
    })
  }

  /**
   * 生成权限菜单
   * @param {Array} menuConfig - 菜单配置
   * @returns {Array} 过滤后的菜单
   */
  const generatePermissionMenu = (menuConfig) => {
    return menuConfig.filter((menu) => {
      // 检查菜单权限
      if (menu.permission && !hasPermission(menu.permission)) {
        return false
      }

      if (menu.role && !hasRole(menu.role)) {
        return false
      }

      // 递归处理子菜单
      if (menu.children && menu.children.length > 0) {
        menu.children = generatePermissionMenu(menu.children)
        // 如果所有子菜单都没权限，隐藏父菜单
        return menu.children.length > 0
      }

      return true
    })
  }

  /**
   * 检查按钮权限
   * @param {string} action - 操作类型
   * @param {string} resource - 资源名称（可选）
   * @returns {boolean} 是否有权限
   */
  const canPerformAction = (action, resource) => {
    if (resource) {
      return hasResourcePermission(resource, action)
    }

    return hasPermission(action)
  }

  return {
    // 状态
    permissionState: readonly(permissionState),

    // 计算属性
    isLoading: computed(() => permissionState.loading),
    hasError: computed(() => !!permissionState.error),

    // 方法
    initPermissions,
    hasPermission,
    hasRole,
    hasResourcePermission,
    canAccessRoute,
    filterRoutesByPermission,
    generatePermissionMenu,
    canPerformAction,
  }
}

/**
 * 权限指令
 */
export const permissionDirective = {
  mounted(el, binding) {
    const { value } = binding
    const { hasPermission } = useRoutePermission()

    if (!hasPermission(value)) {
      el.style.display = 'none'
    }
  },

  updated(el, binding) {
    const { value } = binding
    const { hasPermission } = useRoutePermission()

    if (!hasPermission(value)) {
      el.style.display = 'none'
    }
    else {
      el.style.display = ''
    }
  },
}

/**
 * 角色指令
 */
export const roleDirective = {
  mounted(el, binding) {
    const { value } = binding
    const { hasRole } = useRoutePermission()

    if (!hasRole(value)) {
      el.style.display = 'none'
    }
  },

  updated(el, binding) {
    const { value } = binding
    const { hasRole } = useRoutePermission()

    if (!hasRole(value)) {
      el.style.display = 'none'
    }
    else {
      el.style.display = ''
    }
  },
}

/**
 * 权限工具类
 */
export class PermissionChecker {
  constructor(permissions = [], roles = []) {
    this.permissions = new Set(permissions)
    this.roles = new Set(roles)
  }

  /**
   * 更新权限
   * @param {Array} permissions - 权限列表
   * @param {Array} roles - 角色列表
   */
  update(permissions = [], roles = []) {
    this.permissions = new Set(permissions)
    this.roles = new Set(roles)
  }

  /**
   * 检查权限
   * @param {string | Array} permission - 权限
   * @returns {boolean} 是否有权限
   */
  hasPermission(permission) {
    if (!permission)
      return true

    if (Array.isArray(permission)) {
      return permission.some(p => this.permissions.has(p))
    }

    return this.permissions.has(permission)
  }

  /**
   * 检查所有权限
   * @param {Array} permissions - 权限列表
   * @returns {boolean} 是否有所有权限
   */
  hasAllPermissions(permissions) {
    return permissions.every(p => this.permissions.has(p))
  }

  /**
   * 检查角色
   * @param {string | Array} role - 角色
   * @returns {boolean} 是否有角色
   */
  hasRole(role) {
    if (!role)
      return true

    if (Array.isArray(role)) {
      return role.some(r => this.roles.has(r))
    }

    return this.roles.has(role)
  }

  /**
   * 检查所有角色
   * @param {Array} roles - 角色列表
   * @returns {boolean} 是否有所有角色
   */
  hasAllRoles(roles) {
    return roles.every(r => this.roles.has(r))
  }
}

/**
 * 权限路由守卫
 * @param {object} to - 目标路由
 * @param {object} from - 来源路由
 * @param {Function} next - 下一步函数
 */
export function permissionGuard(to, from, next) {
  const { canAccessRoute } = useRoutePermission()

  if (!canAccessRoute(to)) {
    // 没有权限，跳转到403页面
    next({
      path: '/403',
      query: {
        redirect: to.fullPath,
        reason: 'permission_denied',
      },
    })
    return
  }

  next()
}

/**
 * 菜单权限过滤器
 * @param {Array} menus - 菜单列表
 * @returns {Array} 过滤后的菜单
 */
export function filterMenusByPermission(menus) {
  const { generatePermissionMenu } = useRoutePermission()
  return generatePermissionMenu(menus)
}

/**
 * 路由权限配置
 */
export const ROUTE_PERMISSIONS = {
  // 用户管理
  USER: {
    VIEW: 'user:view',
    CREATE: 'user:create',
    UPDATE: 'user:update',
    DELETE: 'user:delete',
    EXPORT: 'user:export',
  },

  // 角色管理
  ROLE: {
    VIEW: 'role:view',
    CREATE: 'role:create',
    UPDATE: 'role:update',
    DELETE: 'role:delete',
    ASSIGN: 'role:assign',
  },

  // 权限管理
  PERMISSION: {
    VIEW: 'permission:view',
    CREATE: 'permission:create',
    UPDATE: 'permission:update',
    DELETE: 'permission:delete',
  },

  // 系统管理
  SYSTEM: {
    MONITOR: 'system:monitor',
    SETTING: 'system:setting',
    LOG: 'system:logger',
  },
}
