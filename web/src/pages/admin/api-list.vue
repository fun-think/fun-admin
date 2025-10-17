<template>
  <div class="api-list-container">
    <a-card>
      <template #title>
        <div class="card-header">
          <span>API管理</span>
          <a-button type="primary" @click="handleCreate">
            <template #icon>
              <PlusOutlined />
            </template>
            新建API
          </a-button>
        </div>
      </template>
      
      <a-form
        layout="inline"
        :model="searchForm"
        class="search-form"
      >
        <a-form-item label="分组">
          <a-input v-model:value="searchForm.group" placeholder="请输入API分组" />
        </a-form-item>
        <a-form-item label="路径">
          <a-input v-model:value="searchForm.path" placeholder="请输入API路径" />
        </a-form-item>
        <a-form-item label="方法">
          <a-select v-model:value="searchForm.method" style="width: 120px" placeholder="请选择方法">
            <a-select-option value="">全部</a-select-option>
            <a-select-option value="GET">GET</a-select-option>
            <a-select-option value="POST">POST</a-select-option>
            <a-select-option value="PUT">PUT</a-select-option>
            <a-select-option value="DELETE">DELETE</a-select-option>
          </a-select>
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
          <template v-if="column.dataIndex === 'method'">
            <a-tag :color="getMethodTagColor(record.method)">
              {{ record.method }}
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
    
    <!-- API编辑/创建对话框 -->
    <a-modal
      v-model:open="modalVisible"
      :title="modalTitle"
      :confirm-loading="modalConfirmLoading"
      @ok="handleModalOk"
      @cancel="handleModalCancel"
    >
      <a-form
        ref="apiFormRef"
        :model="apiForm"
        :rules="apiFormRules"
        layout="vertical"
      >
        <a-form-item label="分组" name="group">
          <a-input v-model:value="apiForm.group" placeholder="请输入API分组" />
        </a-form-item>
        <a-form-item label="名称" name="name">
          <a-input v-model:value="apiForm.name" placeholder="请输入API名称" />
        </a-form-item>
        <a-form-item label="路径" name="path">
          <a-input v-model:value="apiForm.path" placeholder="请输入API路径" />
        </a-form-item>
        <a-form-item label="方法" name="method">
          <a-select v-model:value="apiForm.method" placeholder="请选择方法">
            <a-select-option value="GET">GET</a-select-option>
            <a-select-option value="POST">POST</a-select-option>
            <a-select-option value="PUT">PUT</a-select-option>
            <a-select-option value="DELETE">DELETE</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { 
  getApis, 
  createApi, 
  updateApi, 
  deleteApi,
  getApi
} from '~/api/common/api.js'
import { PlusOutlined } from '@ant-design/icons-vue'

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
  group: '',
  path: '',
  method: ''
})

// 表格列定义
const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    width: 80
  },
  {
    title: '分组',
    dataIndex: 'group'
  },
  {
    title: '名称',
    dataIndex: 'name'
  },
  {
    title: '路径',
    dataIndex: 'path'
  },
  {
    title: '方法',
    dataIndex: 'method',
    width: 120
  },
  {
    title: '创建时间',
    dataIndex: 'created_at'
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
const modalTitle = ref('新建API')
const isEdit = ref(false)
const currentEditId = ref(null)

// API表单
const apiFormRef = ref()
const apiForm = reactive({
  group: '',
  name: '',
  path: '',
  method: undefined
})

// 表单验证规则
const apiFormRules = {
  group: [{ required: true, message: '请输入API分组' }],
  name: [{ required: true, message: '请输入API名称' }],
  path: [{ required: true, message: '请输入API路径' }],
  method: [{ required: true, message: '请选择方法' }]
}

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

// 获取API列表
const loadApis = async () => {
  try {
    loading.value = true
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      group: searchForm.group,
      path: searchForm.path,
      method: searchForm.method
    }
    
    const res = await getApis(params)
    dataSource.value = res.data.items
    pagination.total = res.data.total
  } catch (err) {
    console.error('获取API列表失败:', err)
    message.error('获取API列表失败')
  } finally {
    loading.value = false
  }
}

// 处理表格分页、排序、筛选变化
const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadApis()
}

// 处理搜索
const handleSearch = () => {
  pagination.current = 1
  loadApis()
}

// 处理重置
const handleReset = () => {
  searchForm.path = ''
  searchForm.method = ''
  pagination.current = 1
  loadApis()
}

// 处理新建
const handleCreate = () => {
  modalTitle.value = '新建API'
  isEdit.value = false
  currentEditId.value = null
  // 重置表单
  Object.assign(apiForm, {
    group: '',
    name: '',
    path: '',
    method: undefined
  })
  modalVisible.value = true
}

// 处理编辑
const handleEdit = async (record) => {
  modalTitle.value = '编辑API'
  isEdit.value = true
  currentEditId.value = record.id
  
  try {
    const res = await getApi(record.id)
    Object.assign(apiForm, res.data)
  } catch (err) {
    console.error('获取API详情失败:', err)
    message.error('获取API详情失败')
    return
  }
  
  modalVisible.value = true
}

// 处理删除
const handleDelete = (record) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除API "${record.path}" 吗？`,
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        await deleteApi(record.id)
        message.success('删除成功')
        loadApis()
      } catch (err) {
        console.error('删除API失败:', err)
        message.error('删除API失败')
      }
    }
  })
}

// 模态框确认
const handleModalOk = () => {
  apiFormRef.value
    .validate()
    .then(async () => {
      modalConfirmLoading.value = true
      try {
        if (isEdit.value) {
          // 编辑
          await updateApi(currentEditId.value, apiForm)
          message.success('更新成功')
        } else {
          // 新建
          await createApi(apiForm)
          message.success('创建成功')
        }
        modalVisible.value = false
        loadApis()
      } catch (err) {
        console.error('保存API失败:', err)
        message.error('保存API失败')
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
  loadApis()
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