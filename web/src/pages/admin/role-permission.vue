<template>
  <div class="role-permission-container">
    <a-card>
      <a-page-header
        style="border: 1px solid rgb(235, 237, 240)"
        title="角色权限管理"
        :subtitle="`当前角色: ${currentRole.name}`"
        @back="() => router.push('/admin/roles')"
      >
        <template #extra>
          <a-button type="primary" @click="handleSave" :loading="saveLoading">保存权限</a-button>
        </template>
      </a-page-header>
      
      <a-spin :spinning="loading">
        <a-alert
          message="提示"
          description="请为角色分配相应的API权限，分配后角色下的用户将拥有对应权限。"
          type="info"
          show-icon
          style="margin-bottom: 20px"
        />
        
        <a-tabs v-model:activeKey="activeKey">
          <a-tab-pane key="1" tab="按API分配">
            <a-tree
              v-model:checkedKeys="checkedApiKeys"
              checkable
              :tree-data="apiTreeData"
              :expanded-keys="expandedApiKeys"
              @expand="onApiExpand"
              style="margin-top: 20px"
            >
              <template #title="{ title, method }">
                <span v-if="method">
                  <a-tag :color="getMethodTagColor(method)" style="margin-right: 8px;">
                    {{ method }}
                  </a-tag>
                  {{ title }}
                </span>
                <span v-else>{{ title }}</span>
              </template>
            </a-tree>
          </a-tab-pane>
          
          <a-tab-pane key="2" tab="按菜单分配" force-render>
            <a-tree
              v-model:checkedKeys="checkedMenuKeys"
              checkable
              :tree-data="menuTreeData"
              :expanded-keys="expandedMenuKeys"
              @expand="onMenuExpand"
              style="margin-top: 20px"
            />
          </a-tab-pane>
        </a-tabs>
      </a-spin>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { getRole } from '~/api/common/role.js'
import { getApis } from '~/api/common/api.js'
import { getRolePermissions, updateRolePermissions } from '~/api/common/permission.js'

const route = useRoute()
const router = useRouter()

// 加载状态
const loading = ref(false)
const saveLoading = ref(false)

// 当前角色
const currentRole = reactive({
  id: null,
  name: '',
  keyword: '',
  description: ''
})

// Tab相关
const activeKey = ref('1')

// API权限树
const apiTreeData = ref([])
const checkedApiKeys = ref([])
const expandedApiKeys = ref([])

// 菜单权限树
const menuTreeData = ref([])
const checkedMenuKeys = ref([])
const expandedMenuKeys = ref([])

// 获取方法标签颜色
const getMethodTagColor = (method) => {
  switch (method) {
    case 'GET':
      return 'green'
    case 'POST':
      return 'blue'
    case 'PUT':
      return 'orange'
    case 'DELETE':
      return 'red'
    default:
      return 'default'
  }
}

// API树展开处理
const onApiExpand = (keys) => {
  expandedApiKeys.value = keys
}

// 菜单树展开处理
const onMenuExpand = (keys) => {
  expandedMenuKeys.value = keys
}

// 获取角色详情
const loadRoleDetail = async () => {
  try {
    loading.value = true
    const roleId = route.params.id
    const res = await getRole(roleId)
    currentRole.id = res.data.id
    currentRole.name = res.data.name
    currentRole.sid = res.data.sid
    currentRole.description = res.data.description
    
    // TODO: 加载角色已有的权限
    loadPermissions()
  } catch (err) {
    console.error('获取角色详情失败:', err)
    message.error('获取角色详情失败')
  } finally {
    loading.value = false
  }
}

// 加载权限数据
const loadPermissions = async () => {
  try {
    // 获取API列表并构建API树
    const res = await getApis({ page: 1, page_size: 1000 })
    
    // 按group分组API
    const groupedApis = {}
    res.data.items.forEach(api => {
      if (!groupedApis[api.group]) {
        groupedApis[api.group] = []
      }
      groupedApis[api.group].push(api)
    })
    
    // 构建API树
    const apiTree = []
    Object.keys(groupedApis).forEach(group => {
      const children = groupedApis[group].map(api => ({
        title: `${api.path}`,
        key: `${api.id}`,
        method: api.method
      }))
      
      apiTree.push({
        title: group,
        key: `group-${group}`,
        children: children
      })
    })
    
    apiTreeData.value = apiTree
    
    // 获取角色已有权限
    const permissionRes = await getRolePermissions(currentRole.sid)
    
    // 设置角色已有的权限
    checkedApiKeys.value = permissionRes.data.list.map(p => {
      // 从权限字符串中提取API ID
      const parts = p.split(':')
      return parts[1] || ''
    }).filter(id => id !== '')
    
    // 默认展开所有节点
    expandedApiKeys.value = apiTree.map(group => group.key)
  } catch (err) {
    console.error('加载权限数据失败:', err)
    message.error('加载权限数据失败')
  }
}

// 保存权限
const handleSave = async () => {
  try {
    saveLoading.value = true
    // 构造权限列表
    const permissions = checkedApiKeys.value.map(key => `api:${key}`)
    
    // 调用保存权限的API
    await updateRolePermissions(currentRole.sid, permissions)
    
    message.success('权限保存成功')
  } catch (err) {
    console.error('保存权限失败:', err)
    message.error('保存权限失败: ' + (err.message || ''))
  } finally {
    saveLoading.value = false
  }
}

// 初始化加载
onMounted(() => {
  loadRoleDetail()
})
</script>

<style scoped>
.role-permission-container {
  height: 100%;
}
</style>
</script>
</file>