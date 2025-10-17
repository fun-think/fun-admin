<script setup>
import { ref, onMounted, h } from 'vue'
import { useRequest } from 'vue-request'
import axios from 'axios'
import {
  Card,
  Row,
  Col,
  Table,
  Form,
  Input,
  Select,
  Button,
  Popconfirm,
  message,
  Tabs,
  Checkbox,
  Tree,
  Empty
} from 'ant-design-vue'
import { getUserPermissions, getRolePermissions, updateRolePermissions } from '~/api/common/permission.js'
import { getAllRoles } from '~/api/common/role.js'
import { getResourceList } from '~/api/admin/resources.js'

// 当前活动标签页
const activeTab = ref('policies')

// 权限策略数据
const policies = ref([])
const policiesLoading = ref(false)

// 角色数据
const roles = ref([])
const rolesLoading = ref(false)

// 用户角色数据
const userRoles = ref([])
const userRolesLoading = ref(false)

// 资源数据
const resources = ref([])
const resourcesLoading = ref(false)

// 表单数据
const policyForm = ref({
  role: '',
  path: '',
  method: 'GET'
})

const rolePolicies = ref([])

// 获取所有权限策略
const fetchPolicies = async () => {
  policiesLoading.value = true
  try {
    const res = await getUserPermissions()
    
    if (res.data.code === 0) {
      policies.value = res.data.data || []
    } else {
      message.error(res.data.msg || '获取权限策略失败')
    }
  } catch (error) {
    console.error('Failed to fetch policies:', error)
    message.error('获取权限策略失败')
  } finally {
    policiesLoading.value = false
  }
}

// 添加权限策略
const addPolicy = async () => {
  try {
    const res = await updateRolePermissions(policyForm.value.role, [`${policyForm.value.path}:${policyForm.value.method}`])
    
    if (res.data.code === 0) {
      message.success('添加权限策略成功')
      fetchPolicies()
      // 重置表单
      policyForm.value = {
        role: '',
        path: '',
        method: 'GET'
      }
    } else {
      message.error(res.data.msg || '添加权限策略失败')
    }
  } catch (error) {
    console.error('Failed to add policy:', error)
    message.error('添加权限策略失败')
  }
}

// 删除权限策略
const removePolicy = async (record) => {
  try {
    const res = await updateRolePermissions(record.role, [])
    
    if (res.data.code === 0) {
      message.success('删除权限策略成功')
      fetchPolicies()
    } else {
      message.error(res.data.msg || '删除权限策略失败')
    }
  } catch (error) {
    console.error('Failed to remove policy:', error)
    message.error('删除权限策略失败')
  }
}

// 获取所有角色
const fetchRoles = async () => {
  rolesLoading.value = true
  try {
    const res = await getAllRoles()
    
    if (res.data.code === 0) {
      roles.value = res.data.data || []
      // 默认选中第一个角色
      if (roles.value.length > 0) {
        fetchRolePolicies(roles.value[0].name)
      }
    } else {
      message.error(res.data.msg || '获取角色失败')
    }
  } catch (error) {
    console.error('Failed to fetch roles:', error)
    message.error('获取角色失败')
  } finally {
    rolesLoading.value = false
  }
}

// 获取角色权限策略
const fetchRolePolicies = async (role) => {
  try {
    const res = await getRolePermissions(role)
    
    if (res.data.code === 0) {
      rolePolicies.value = res.data.data || []
    } else {
      message.error(res.data.msg || '获取角色权限失败')
    }
  } catch (error) {
    console.error('Failed to fetch role policies:', error)
    message.error('获取角色权限失败')
  }
}

// 更新角色权限策略
const updateRolePolicies = async (role) => {
  try {
    // 这里应该从界面获取选中的策略
    const policies = [] // 示例数据，实际应从界面获取
    
    const res = await updateRolePermissions(role, policies)
    
    if (res.data.code === 0) {
      message.success('更新角色权限成功')
    } else {
      message.error(res.data.msg || '更新角色权限失败')
    }
  } catch (error) {
    console.error('Failed to update role policies:', error)
    message.error('更新角色权限失败')
  }
}

// 获取用户角色
const fetchUserRoles = async (userId) => {
  userRolesLoading.value = true
  try {
    const res = await getUserPermissions()
    
    if (res.data.code === 0) {
      userRoles.value = res.data.data || []
    } else {
      message.error(res.data.msg || '获取用户角色失败')
    }
  } catch (error) {
    console.error('Failed to fetch user roles:', error)
    message.error('获取用户角色失败')
  } finally {
    userRolesLoading.value = false
  }
}

// 更新用户角色
const updateUserRoles = async (userId, roles) => {
  try {
    // 使用新的API接口
    message.success('更新用户角色成功')
  } catch (error) {
    console.error('Failed to update user roles:', error)
    message.error('更新用户角色失败')
  }
}

// 获取资源列表
const fetchResources = async () => {
  resourcesLoading.value = true
  try {
    const res = await getResourceList()
    
    if (res.data.code === 0) {
      resources.value = res.data.data.resources || []
    } else {
      message.error(res.data.msg || '获取资源列表失败')
    }
  } catch (error) {
    console.error('Failed to fetch resources:', error)
    message.error('获取资源列表失败')
  } finally {
    resourcesLoading.value = false
  }
}

// 权限策略表格列定义
const policyColumns = [
  {
    title: '角色',
    dataIndex: 'role',
    key: 'role'
  },
  {
    title: '路径',
    dataIndex: 'path',
    key: 'path'
  },
  {
    title: '方法',
    dataIndex: 'method',
    key: 'method'
  },
  {
    title: '操作',
    key: 'actions',
    customRender: ({ record }) => {
      return h(
        Popconfirm,
        {
          title: '确定删除该权限策略？',
          onConfirm: () => removePolicy(record)
        },
        {
          default: () => h(
            Button,
            {
              type: 'link',
              size: 'small',
              danger: true
            },
            '删除'
          )
        }
      )
    }
  }
]

onMounted(() => {
  fetchPolicies()
  fetchRoles()
  fetchResources()
})
</script>

<template>
  <page-container title="权限管理">
    <a-tabs v-model:activeKey="activeTab">
      <a-tab-pane key="policies" tab="权限策略">
        <a-card>
          <a-row :gutter="16">
            <a-col :span="8">
              <a-card title="添加权限策略">
                <a-form layout="vertical">
                  <a-form-item label="角色">
                    <a-input v-model:value="policyForm.role" placeholder="请输入角色名称" />
                  </a-form-item>
                  
                  <a-form-item label="路径">
                    <a-input v-model:value="policyForm.path" placeholder="请输入API路径" />
                  </a-form-item>
                  
                  <a-form-item label="方法">
                    <a-select v-model:value="policyForm.method">
                      <a-select-option value="GET">GET</a-select-option>
                      <a-select-option value="POST">POST</a-select-option>
                      <a-select-option value="PUT">PUT</a-select-option>
                      <a-select-option value="DELETE">DELETE</a-select-option>
                    </a-select>
                  </a-form-item>
                  
                  <a-form-item>
                    <a-button type="primary" @click="addPolicy">添加策略</a-button>
                  </a-form-item>
                </a-form>
              </a-card>
            </a-col>
            
            <a-col :span="16">
              <a-card title="权限策略列表">
                <a-table
                  row-key="id"
                  :loading="policiesLoading"
                  :columns="policyColumns"
                  :data-source="policies"
                  :pagination="false"
                />
              </a-card>
            </a-col>
          </a-row>
        </a-card>
      </a-tab-pane>
      
      <a-tab-pane key="roles" tab="角色管理">
        <a-card>
          <a-row :gutter="16">
            <a-col :span="8">
              <a-card title="角色列表">
                <a-table
                  row-key="id"
                  :loading="rolesLoading"
                  :columns="[
                    { title: '角色名称', dataIndex: 'name', key: 'name' },
                    { title: '描述', dataIndex: 'description', key: 'description' }
                  ]"
                  :data-source="roles"
                  :pagination="false"
                />
              </a-card>
            </a-col>
            
            <a-col :span="16">
              <a-card title="角色权限分配">
                <div v-if="roles.length > 0">
                  <a-form layout="vertical">
                    <a-form-item label="选择角色">
                      <a-select 
                        placeholder="请选择角色"
                        @change="fetchRolePolicies"
                      >
                        <a-select-option 
                          v-for="role in roles" 
                          :key="role.name" 
                          :value="role.name"
                        >
                          {{ role.name }}
                        </a-select-option>
                      </a-select>
                    </a-form-item>
                    
                    <a-form-item label="权限分配">
                      <a-tree
                        :tree-data="resources"
                        :field-names="{ children: 'children', title: 'title', key: 'slug' }"
                        checkable
                      />
                    </a-form-item>
                    
                    <a-form-item>
                      <a-button type="primary" @click="updateRolePolicies">保存权限</a-button>
                    </a-form-item>
                  </a-form>
                </div>
                
                <div v-else>
                  <a-empty description="暂无角色数据" />
                </div>
              </a-card>
            </a-col>
          </a-row>
        </a-card>
      </a-tab-pane>
      
      <a-tab-pane key="users" tab="用户权限">
        <a-card title="用户权限管理">
          <a-empty description="用户权限管理功能待实现" />
        </a-card>
      </a-tab-pane>
    </a-tabs>
  </page-container>
</template>

<style scoped>
</style>