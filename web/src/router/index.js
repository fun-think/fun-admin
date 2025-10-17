import { createRouter, createWebHistory } from 'vue-router'
import { useTitle } from '@vueuse/core'
import { useUserStore } from '~/stores/user'
import { useAuthorization } from '~/composables/authorization'

// 静态路由配置
const Layout = () => import('~/layouts/index.vue')

const staticRoutes = [
  {
    path: '/',
    redirect: '/admin/dashboard'
  },
  {
    path: '/login',
    component: () => import('~/pages/common/login.vue'),
    meta: {
      title: '登录',
      auth: false
    },
  },
  {
    path: '/401',
    component: () => import('~/pages/exception/401.vue'),
    meta: {
      title: '授权已过期',
      auth: false
    },
  },
  {
    path: '/admin',
    component: Layout,
    redirect: '/admin/dashboard',
    meta: {
      auth: true
    },
    children: [
      {
        path: 'dashboard',
        component: () => import('~/pages/admin/dashboard.vue'),
        meta: {
          title: '仪表板',
          icon: 'dashboard'
        }
      },
      {
        path: ':slug',
        component: () => import('~/pages/admin/dynamic-list.vue'),
        meta: {
          title: '列表'
        }
      },
      {
        path: ':slug/create',
        component: () => import('~/pages/admin/dynamic-form.vue'),
        meta: {
          title: '新增'
        }
      },
      {
        path: ':slug/edit/:id',
        component: () => import('~/pages/admin/dynamic-form.vue'),
        meta: {
          title: '编辑'
        }
      },
      {
        path: 'profile',
        component: () => import('~/pages/admin/profile.vue'),
        meta: {
          title: '个人资料'
        }
      },
      {
        path: 'apis',
        component: () => import('~/pages/admin/api-list.vue'),
        meta: {
          title: 'API管理'
        }
      },
      {
        path: 'roles',
        component: () => import('~/pages/admin/role-list.vue'),
        meta: {
          title: '角色管理'
        }
      },
      {
        path: 'role-permission/:id',
        component: () => import('~/pages/admin/role-permission.vue'),
        meta: {
          title: '角色权限'
        }
      },
      {
        path: 'menus',
        component: () => import('~/pages/admin/menu-list.vue'),
        meta: {
          title: '菜单管理'
        }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    component: () => import('~/pages/exception/404.vue'),
    meta: {
      title: '页面未找到',
      auth: false
    }
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(import.meta.env.VITE_APP_BASE),
  routes: staticRoutes
})

// 设置路由守卫
router.beforeEach(async (to, from, next) => {
  // 更新页面标题
  const title = useTitle()
  title.value = to.meta?.title || 'FunAdmin'
  
  // 检查是否需要认证
  const needAuth = to.meta?.auth !== false
  const token = useAuthorization()
  
  // 不需要认证的页面直接访问
  if (!needAuth) {
    next()
    return
  }
  
  // 需要认证但没有token，跳转到登录页
  if (!token.value) {
    next({ 
      path: '/login',
      query: to.path !== '/' ? { redirect: to.fullPath } : {}
    })
    return
  }
  
  // 已登录用户访问登录页时重定向到仪表板
  if (to.path === '/login') {
    next({ path: '/admin/dashboard' })
    return
  }
  
  // 已登录但用户信息不存在时获取用户信息
  const userStore = useUserStore()
  if (!userStore.user || !userStore.user.value) {
    try {
      await userStore.getUserInfo()
    } catch (error) {
      // 获取用户信息失败，清除登录状态并重定向到登录页
      userStore.logout()
      next({ 
        path: '/login',
        query: { redirect: to.fullPath } 
      })
      return
    }
  }
  
  next()
})

router.afterEach(() => {
  // 滚动到顶部
  window.scrollTo(0, 0)
})

export default router