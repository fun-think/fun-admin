<script setup>
import { ref, computed, h, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { LayoutSider, Menu, Avatar, Dropdown, Badge } from 'ant-design-vue'
import {
  DashboardOutlined,
  SettingOutlined,
  FileSearchOutlined,
  UserOutlined,
  LogoutOutlined,
  SafetyCertificateOutlined,
  BellOutlined,
  MonitorOutlined,
  DatabaseOutlined
} from '@ant-design/icons-vue'
import { getResourceList } from '~/api/admin/resources'
import { useUserStore } from '~/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 菜单展开状态
const collapsed = defineModel('collapsed')

// 当前选中的菜单项
const selectedKeys = computed(() => {
  const path = route.path
  if (path.startsWith('/admin/dashboard')) {
    return ['dashboard']
  } else if (path.startsWith('/admin/settings')) {
    return ['settings']
  } else if (path.startsWith('/admin/operation-logs')) {
    return ['operation-logs']
  } else if (path.startsWith('/admin/profile')) {
    return ['profile']
  } else if (path.startsWith('/admin/permissions')) {
    return ['permissions']
  } else if (path.startsWith('/admin/notifications')) {
    return ['notifications']
  } else if (path.startsWith('/admin/system-monitors')) {
    return ['system-monitors']
  } else if (path.startsWith('/admin/system-settings')) {
    return ['system-settings']
  } else {
    // 动态资源路由
    const matches = path.match(/^\/admin\/([^\/]+)/)
    if (matches && matches[1]) {
      return [matches[1]]
    }
    return []
  }
})

// 静态菜单项（固定）
const staticItems = ref([
  { key: 'dashboard', icon: () => h(DashboardOutlined), label: '仪表板', path: '/admin/dashboard' },
  { key: 'settings', icon: () => h(SettingOutlined), label: '系统设置', path: '/admin/settings' },
  { key: 'operation-logs', icon: () => h(FileSearchOutlined), label: '操作日志', path: '/admin/operation-logs' },
  { key: 'permissions', icon: () => h(SafetyCertificateOutlined), label: '权限管理', path: '/admin/permissions' },
  { key: 'notifications', icon: () => h(BellOutlined), label: '通知中心', path: '/admin/notifications' },
  { key: 'system-monitors', icon: () => h(MonitorOutlined), label: '系统监控', path: '/admin/system-monitors' },
  { key: 'system-settings', icon: () => h(DatabaseOutlined), label: '系统配置', path: '/admin/system-settings' },
  { key: 'profile', icon: () => h(UserOutlined), label: '个人资料', path: '/admin/profile' },
])

// 动态资源菜单（来自后端）
const resourceGroups = ref([]) // [{ group: '系统管理', items: [{key,label,icon,badge}] }]

const loadResources = async () => {
  const res = await getResourceList({ language: localStorage.getItem('admin-language') || 'zh-CN' })
  if (res.code !== 0) return
  const list = res.data?.resources || []
  // 分组
  const groups = new Map()
  list.forEach(r => {
    const group = r.nav_group || '资源'
    if (!groups.has(group)) groups.set(group, [])
    groups.get(group).push(r)
  })
  // 构建并排序
  resourceGroups.value = Array.from(groups.entries()).map(([group, items]) => ({
    group,
    items: items
      .filter(i => i.slug)
      .sort((a, b) => (a.nav_sort || 0) - (b.nav_sort || 0))
      .map(i => ({
        key: i.slug,
        label: i.title,
        path: `/admin/${i.slug}`,
        icon: i.nav_icon,
        badge: i.nav_badge_count,
      })),
  }))
}

onMounted(() => {
  loadResources()
})

// 处理菜单点击
const handleMenuClick = ({ key }) => {
  // 先查静态项
  if (key === 'global-search') {
    router.push('/admin/global-search')
    return
  }
  const staticItem = staticItems.value.find(item => item.key === key)
  if (staticItem) {
    router.push(staticItem.path)
    return
  }
  // 动态项
  for (const grp of resourceGroups.value) {
    const item = grp.items.find(i => i.key === key)
    if (item) {
      router.push(item.path)
      return
    }
  }
}

// 用户退出登录
const handleLogout = () => {
  useUserStore().logout()
  router.push('/login')
}
</script>

<template>
  <a-layout-sider 
    v-model:collapsed="collapsed" 
    :trigger="null" 
    collapsible
    width="256"
    class="sidebar"
  >
    <div class="logo">
      <img src="/logo.png" alt="Logo" v-if="!collapsed" />
      <span v-if="!collapsed">管理系统</span>
    </div>

    <a-menu
      v-model:selectedKeys="selectedKeys"
      mode="inline"
      theme="dark"
      @click="handleMenuClick"
    >
      <a-menu-item key="global-search">
        <span>全局搜索</span>
      </a-menu-item>
      <a-menu-item v-for="item in staticItems" :key="item.key">
        <component :is="item.icon" />
        <span>{{ item.label }}</span>
      </a-menu-item>

      <template v-for="group in resourceGroups" :key="group.group">
        <a-menu-item-group :title="group.group">
          <a-menu-item v-for="it in group.items" :key="it.key">
            <span>{{ it.label }}</span>
            <template #icon>
              <span v-if="it.icon" class="menu-icon">{{ it.icon }}</span>
            </template>
            <template #title>
              <span>{{ it.label }}</span>
              <a-badge v-if="it.badge != null" :count="it.badge" :overflow-count="99" offset="[8,0]" />
            </template>
          </a-menu-item>
        </a-menu-item-group>
      </template>
    </a-menu>

    <div class="user-info">
      <a-dropdown>
        <div class="user-dropdown">
          <a-avatar :size="32" icon="user" />
          <span v-if="!collapsed" class="username">{{ userStore.user?.nickname || userStore.user?.username }}</span>
        </div>

        <template #overlay>
          <a-menu>
            <a-menu-item key="profile" @click="() => router.push('/admin/profile')">
              <UserOutlined />
              个人资料
            </a-menu-item>
            <a-menu-item key="logout" @click="handleLogout">
              <LogoutOutlined />
              退出登录
            </a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </div>
  </a-layout-sider>
</template>

<style scoped>
.sidebar {
  position: relative;
  height: 100vh;
  overflow-y: auto;
}
.logo { display:flex; align-items:center; justify-content:center; height:64px; background: rgba(255,255,255,0.05); color:#fff; font-size:18px; font-weight:bold; }
.logo img{ height:32px; margin-right:12px; }
.user-info { position:absolute; bottom:0; width:100%; padding:16px; border-top:1px solid rgba(255,255,255,0.05); }
.user-dropdown { display:flex; align-items:center; cursor:pointer; }
.username { margin-left:12px; color:white; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }
.menu-icon { margin-right: 8px; opacity: 0.85; }
</style>