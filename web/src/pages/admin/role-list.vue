<template>
  <div class="role-list-container">
    <a-card>
      <template #title>
        <div class="card-header">
          <span>角色管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon>
              <PlusOutlined />
            </template>
            新建角色
          </a-button>
        </div>
      </template>
      
      <a-form
        layout="inline"
        :model="searchForm"
        class="search-form"
      >
        <a-form-item label="角色名称">
          <a-input v-model:value="searchForm.name" placeholder="请输入角色名称" />
        </a-form-item>
        <a-form-item label="角色ID">
          <a-input v-model:value="searchForm.sid" placeholder="请输入角色ID" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">查询</a-button>
          <a-button style="margin-left: 8px" @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
      
      <a-table
        :columns="columns"
        :data-source="dataSource"
        :pagination="pagination"
        :loading="loading"
        @change="handleTableChange"
        bordered
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'status'">
            <a-tag :color="record.status === 1 ? 'green' : 'red'">
              {{ record.status === 1 ? '启用' : '禁用' }}
            </a-tag>
          </template>
          <template v-else-if="column.dataIndex === 'action'">
            <a-space>
              <a-button type="link" @click="handleEdit(record)">编辑</a-button>
              <a-button type="link" @click="handlePermission(record)">权限</a-button>
              <a-button type="link" @click="handleDelete(record)" danger>删除</a-button>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>
    
    <!-- 角色编辑/创建对话框 -->
    <a-modal
      v-model:open="modalVisible"
      :title="modalTitle"
      :confirm-loading="modalConfirmLoading"
      @ok="handleModalOk"
      @cancel="handleModalCancel"
    >
      <a-form
        ref="roleFormRef"
        :model="roleForm"
        :rules="roleFormRules"
        layout="vertical"
      >
        <a-form-item label="角色名称" name="name">
          <a-input v-model:value="roleForm.name" placeholder="请输入角色名称" />
        </a-form-item>
        <a-form-item label="角色ID" name="sid">
          <a-input v-model:value="roleForm.sid" placeholder="请输入角色ID" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { 
  getRoles, 
  createRole, 
  updateRole, 
  deleteRole,
  getRole
} from '~/api/common/role.js'
import { PlusOutlined } from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'

// 路由
const router = useRouter()

// 表格相关
const dataSource = ref([])
const loading = ref(false)
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

// 搜索表单
const searchForm = reactive({
  name: '',
  sid: ''
})

// 表格列定义
const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    width: 80
  },
  {
    title: '角色名称',
    dataIndex: 'name'
  },
  {
    title: '角色ID',
    dataIndex: 'sid'
  },
  {
    title: '创建时间',
    dataIndex: 'created_at'
  },
  {
    title: '操作',
    dataIndex: 'action',
    width: 200
  }
]

// 模态框相关
const modalVisible = ref(false)
const modalConfirmLoading = ref(false)
const modalTitle = ref('新建角色')
const isEdit = ref(false)
const currentEditId = ref(null)

// 角色表单
const roleFormRef = ref()
const roleForm = reactive({
  name: '',
  sid: ''
})

// 表单验证规则
const roleFormRules = {
  name: [{ required: true, message: '请输入角色名称' }],
  sid: [{ required: true, message: '请输入角色ID' }]
}

// 获取角色列表
const loadRoles = async () => {
  try {
    loading.value = true
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      name: searchForm.name,
      sid: searchForm.sid
    }
    
    const res = await getRoles(params)
    dataSource.value = res.data.items
    pagination.total = res.data.total
  } catch (err) {
    console.error('获取角色列表失败:', err)
    message.error('获取角色列表失败')
  } finally {
    loading.value = false
  }
}

// 处理表格分页、排序、筛选变化
const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadRoles()
}

// 处理搜索
const handleSearch = () => {
  pagination.current = 1
  loadRoles()
}

// 处理重置
const handleReset = () => {
  searchForm.name = ''
  pagination.current = 1
  loadRoles()
}

// 处理新建
const handleCreate = () => {
  modalTitle.value = '新建角色'
  isEdit.value = false
  currentEditId.value = null
  // 重置表单
  Object.assign(roleForm, {
    name: '',
    sid: ''
  })
  modalVisible.value = true
}

// 处理编辑
const handleEdit = async (record) => {
  modalTitle.value = '编辑角色'
  isEdit.value = true
  currentEditId.value = record.id
  
  try {
    const res = await getRole(record.id)
    Object.assign(roleForm, res.data)
  } catch (err) {
    console.error('获取角色详情失败:', err)
    message.error('获取角色详情失败')
    return
  }
  
  modalVisible.value = true
}

// 处理权限设置
const handlePermission = (record) => {
  router.push(`/admin/role-permission/${record.id}`)
}

// 处理删除
const handleDelete = (record) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除角色 "${record.name}" 吗？`,
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        await deleteRole(record.id)
        message.success('删除成功')
        loadRoles()
      } catch (err) {
        console.error('删除角色失败:', err)
        message.error('删除角色失败')
      }
    }
  })
}

// 模态框确认
const handleModalOk = () => {
  roleFormRef.value
    .validate()
    .then(async () => {
      modalConfirmLoading.value = true
      try {
        if (isEdit.value) {
          // 编辑
          await updateRole(currentEditId.value, roleForm)
          message.success('更新成功')
        } else {
          // 新建
          await createRole(roleForm)
          message.success('创建成功')
        }
        modalVisible.value = false
        loadRoles()
      } catch (err) {
        console.error('保存角色失败:', err)
        message.error('保存角色失败')
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
  loadRoles()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 16px;
}
</style>
</script>
</file>