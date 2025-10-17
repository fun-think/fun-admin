<template>
  <div class="menu-list-container">
    <a-card>
      <template #title>
        <div class="card-header">
          <span>菜单管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon>
              <PlusOutlined />
            </template>
            新建菜单
          </a-button>
        </div>
      </template>
      
      <a-alert
        message="提示"
        description="菜单用于控制前端导航栏的显示，只有具有相应权限的用户才能看到对应菜单。"
        type="info"
        show-icon
        style="margin-bottom: 20px"
      />
      
      <a-table
        :columns="columns"
        :data-source="dataSource"
        :loading="loading"
        :pagination="false"
        bordered
        childrenColumnName="items"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'icon' && record.icon">
            <component :is="record.icon" />
          </template>
          <template v-else-if="column.dataIndex === 'hideInMenu'">
            <a-tag :color="record.hideInMenu ? 'red' : 'green'">
              {{ record.hideInMenu ? '隐藏' : '显示' }}
            </a-tag>
          </template>
          <template v-else-if="column.dataIndex === 'action'">
            <a-space>
              <a-button type="link" @click="handleEdit(record)">编辑</a-button>
              <a-button type="link" @click="handleDelete(record)" danger>删除</a-button>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>
    
    <!-- 菜单编辑/创建对话框 -->
    <a-modal
      v-model:open="modalVisible"
      :title="modalTitle"
      :confirm-loading="modalConfirmLoading"
      @ok="handleModalOk"
      @cancel="handleModalCancel"
      width="600px"
    >
      <a-form
        ref="menuFormRef"
        :model="menuForm"
        :rules="menuFormRules"
        layout="vertical"
      >
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="父级菜单" name="parentId">
              <a-tree-select
                v-model:value="menuForm.parentId"
                :tree-data="menuTreeData"
                placeholder="请选择父级菜单"
                tree-default-expand-all
                :field-names="{ children: 'children', label: 'title', value: 'id' }"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="排序权重" name="weight">
              <a-input-number v-model:value="menuForm.weight" :min="0" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
        
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="菜单名称" name="title">
              <a-input v-model:value="menuForm.title" placeholder="请输入菜单名称" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="路由名称" name="name">
              <a-input v-model:value="menuForm.name" placeholder="请输入路由名称" />
            </a-form-item>
          </a-col>
        </a-row>
        
        <a-form-item label="路由路径" name="path">
          <a-input v-model:value="menuForm.path" placeholder="请输入路由路径" />
        </a-form-item>
        
        <a-form-item label="组件路径" name="component">
          <a-input v-model:value="menuForm.component" placeholder="请输入组件路径，如: ~/pages/admin/users.vue" />
        </a-form-item>
        
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="图标" name="icon">
              <a-input v-model:value="menuForm.icon" placeholder="请输入图标名称" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="重定向路径" name="redirect">
              <a-input v-model:value="menuForm.redirect" placeholder="请输入重定向路径" />
            </a-form-item>
          </a-col>
        </a-row>
        
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="是否隐藏">
              <a-switch v-model:checked="menuForm.hideInMenu" />
              <span style="margin-left: 10px;">{{ menuForm.hideInMenu ? '隐藏' : '显示' }}</span>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="是否保活">
              <a-switch v-model:checked="menuForm.keepAlive" />
              <span style="margin-left: 10px;">{{ menuForm.keepAlive ? '保活' : '不保活' }}</span>
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { 
  getAdminMenusApi,
  createMenuApi,
  updateMenuApi,
  deleteMenusApi
} from '~/api/common/menu.js'
import { PlusOutlined } from '@ant-design/icons-vue'

// 表格相关
const dataSource = ref([])
const loading = ref(false)

// 表格列定义
const columns = [
  {
    title: '名称',
    dataIndex: 'title',
    width: 200
  },
  {
    title: '路径',
    dataIndex: 'path'
  },
  {
    title: '组件',
    dataIndex: 'component'
  },
  {
    title: '图标',
    dataIndex: 'icon'
  },
  {
    title: '隐藏',
    dataIndex: 'hideInMenu',
    width: 100
  },
  {
    title: '排序',
    dataIndex: 'weight',
    width: 80
  },
  {
    title: '操作',
    dataIndex: 'action',
    width: 150
  }
]

// 模态框相关
const modalVisible = ref(false)
const modalConfirmLoading = ref(false)
const modalTitle = ref('新建菜单')
const isEdit = ref(false)
const currentEditId = ref(null)

// 菜单表单
const menuFormRef = ref()
const menuForm = reactive({
  parentId: undefined,
  weight: 0,
  path: '',
  title: '',
  name: '',
  component: '',
  icon: '',
  redirect: '',
  keepAlive: false,
  hideInMenu: false
})

// 表单验证规则
const menuFormRules = {
  title: [{ required: true, message: '请输入菜单名称' }],
  path: [{ required: true, message: '请输入路由路径' }],
  component: [{ required: true, message: '请输入组件路径' }]
}

// 菜单树数据（用于父级菜单选择）
const menuTreeData = ref([])

// 获取菜单列表
const loadMenus = async () => {
  try {
    loading.value = true
    const res = await getAdminMenusApi()
    
    // 构建菜单树
    const menuMap = {}
    const rootMenus = []
    
    // 先建立映射关系
    res.data.list.forEach(menu => {
      menuMap[menu.id] = { ...menu, items: [] }
    })
    
    // 构建树结构
    res.data.list.forEach(menu => {
      if (menu.parentId === 0) {
        rootMenus.push(menuMap[menu.id])
      } else if (menuMap[menu.parentId]) {
        menuMap[menu.parentId].items.push(menuMap[menu.id])
      }
    })
    
    dataSource.value = rootMenus
    
    // 构建菜单树数据（用于选择父级菜单）
    menuTreeData.value = [{
      id: 0,
      title: '根菜单',
      children: [...rootMenus]
    }]
  } catch (err) {
    console.error('获取菜单列表失败:', err)
    message.error('获取菜单列表失败')
  } finally {
    loading.value = false
  }
}

// 处理新建
const handleCreate = () => {
  modalTitle.value = '新建菜单'
  isEdit.value = false
  currentEditId.value = null
  // 重置表单
  Object.assign(menuForm, {
    parentId: undefined,
    weight: 0,
    path: '',
    title: '',
    name: '',
    component: '',
    icon: '',
    redirect: '',
    keepAlive: false,
    hideInMenu: false
  })
  modalVisible.value = true
}

// 处理编辑
const handleEdit = async (record) => {
  modalTitle.value = '编辑菜单'
  isEdit.value = true
  currentEditId.value = record.id
  
  // 填充表单数据
  Object.assign(menuForm, {
    parentId: record.parentId,
    weight: record.weight,
    path: record.path,
    title: record.title,
    name: record.name,
    component: record.component,
    icon: record.icon,
    redirect: record.redirect,
    keepAlive: record.keepAlive,
    hideInMenu: record.hideInMenu
  })
  
  modalVisible.value = true
}

// 处理删除
const handleDelete = (record) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除菜单 "${record.title}" 吗？`,
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        await deleteMenusApi({ id: record.id })
        message.success('删除成功')
        loadMenus()
      } catch (err) {
        console.error('删除菜单失败:', err)
        message.error('删除菜单失败')
      }
    }
  })
}

// 模态框确认
const handleModalOk = () => {
  menuFormRef.value
    .validate()
    .then(async () => {
      modalConfirmLoading.value = true
      try {
        if (isEdit.value) {
          // 编辑
          await updateMenuApi({ 
            id: currentEditId.value,
            ...menuForm
          })
          message.success('更新成功')
        } else {
          // 新建
          await createMenuApi(menuForm)
          message.success('创建成功')
        }
        modalVisible.value = false
        loadMenus()
      } catch (err) {
        console.error('保存菜单失败:', err)
        message.error('保存菜单失败')
      } finally {
        modalConfirmLoading.value = false
      }
    })
    .catch((err) => {
      console.error('表单验证失败:', err)
    })
}

// 模态框取消
const handleModalCancel = () => {
  modalVisible.value = false
}

// 初始化加载
onMounted(() => {
  loadMenus()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
</script>
</file>