<script setup>
import { ref, onMounted, computed } from 'vue'
import {
  Card,
  Table,
  Button,
  Form,
  Input,
  Modal,
  message,
  Popconfirm,
  Space,
  Tag,
  Drawer,
  Descriptions,
  Badge
} from 'ant-design-vue'
import {
  getDictTypes,
  createDictType,
  updateDictType,
  deleteDictType,
  refreshDictCache
} from '~/api/common/dict.js'

// 表格数据
const tableData = ref([])
const loading = ref(false)
const pagination = ref({
  current: 1,
  pageSize: 10,
  total: 0
})

// 表单相关
const formModalVisible = ref(false)
const formModalTitle = ref('')
const formRef = ref()
const formModel = ref({
  name: '',
  code: '',
  description: '',
  status: 1,
  sort: 0
})

// 详情相关
const detailDrawerVisible = ref(false)
const currentRecord = ref(null)

// 表格列定义
const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
    width: 80
  },
  {
    title: '名称',
    dataIndex: 'name',
    key: 'name',
    width: 150
  },
  {
    title: '编码',
    dataIndex: 'code',
    key: 'code',
    width: 150
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
    ellipsis: true
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    width: 100,
    customRender: ({ text }) => h(
      Badge,
      {
        status: text === 1 ? 'success' : 'error',
        text: text === 1 ? '启用' : '禁用'
      }
    )
  },
  {
    title: '排序',
    dataIndex: 'sort',
    key: 'sort',
    width: 80
  },
  {
    title: '数据量',
    dataIndex: 'data_count',
    key: 'data_count',
    width: 100
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    key: 'created_at',
    width: 150,
    customRender: ({ text }) => text ? new Date(text).toLocaleString() : '-'
  },
  {
    title: '操作',
    key: 'actions',
    fixed: 'right',
    width: 200,
    customRender: ({ record }) => {
      return h(Space, {}, () => [
        h(
          Button,
          {
            type: 'link',
            size: 'small',
            onClick: () => viewDictData(record)
          },
          () => '数据管理'
        ),
        h(
          Button,
          {
            type: 'link',
            size: 'small',
            onClick: () => editRecord(record)
          },
          () => '编辑'
        ),
        h(
          Popconfirm,
          {
            title: '确定删除该字典类型？',
            onConfirm: () => deleteRecord(record.id)
          },
          {
            default: () => h(
              Button,
              {
                type: 'link',
                size: 'small',
                danger: true
              },
              () => '删除'
            )
          }
        )
      ])
    }
  }
]

// 获取表格数据
const fetchTableData = async () => {
  loading.value = true
  try {
    const res = await getDictTypes({
      page: pagination.value.current,
      page_size: pagination.value.pageSize
    })
    
    if (res.code === 0) {
      tableData.value = res.data.items || []
      pagination.value.total = res.data.total || 0
    }
  } catch (error) {
    console.error('Failed to fetch dict types:', error)
    message.error('获取字典类型失败')
  } finally {
    loading.value = false
  }
}

// 新增记录
const addRecord = () => {
  formModel.value = {
    name: '',
    code: '',
    description: '',
    status: 1,
    sort: 0
  }
  formModalTitle.value = '新增字典类型'
  formModalVisible.value = true
}

// 编辑记录
const editRecord = (record) => {
  formModel.value = { ...record }
  formModalTitle.value = '编辑字典类型'
  formModalVisible.value = true
}

// 查看字典数据
const viewDictData = (record) => {
  currentRecord.value = record
  detailDrawerVisible.value = true
}

// 删除记录
const deleteRecord = async (id) => {
  try {
    const res = await deleteDictType(id)
    if (res.code === 0) {
      message.success('删除成功')
      fetchTableData()
    } else {
      message.error(res.msg || '删除失败')
    }
  } catch (error) {
    console.error('Failed to delete record:', error)
    message.error('删除失败')
  }
}

// 保存记录
const saveRecord = async () => {
  try {
    await formRef.value.validateFields()
    
    const isEdit = !!formModel.value.id
    const res = isEdit 
      ? await updateDictType(formModel.value.id, formModel.value)
      : await createDictType(formModel.value)
    
    if (res.code === 0) {
      message.success(isEdit ? '更新成功' : '创建成功')
      formModalVisible.value = false
      fetchTableData()
    } else {
      message.error(res.msg || (isEdit ? '更新失败' : '创建失败'))
    }
  } catch (error) {
    console.error('Failed to save record:', error)
    message.error('保存失败')
  }
}

// 刷新缓存
const handleRefreshCache = async () => {
  try {
    const res = await refreshDictCache()
    if (res.code === 0) {
      message.success('缓存刷新成功')
    } else {
      message.error(res.msg || '缓存刷新失败')
    }
  } catch (error) {
    console.error('Failed to refresh cache:', error)
    message.error('缓存刷新失败')
  }
}

// 表格变化处理
const handleTableChange = (pag) => {
  pagination.value.current = pag.current
  pagination.value.pageSize = pag.pageSize
  fetchTableData()
}

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入字典名称' }
  ],
  code: [
    { required: true, message: '请输入字典编码' },
    { pattern: /^[A-Z_]+$/, message: '字典编码只能包含大写字母和下划线' }
  ],
  status: [
    { required: true, message: '请选择状态' }
  ]
}

onMounted(() => {
  fetchTableData()
})
</script>

<template>
  <page-container title="字典类型管理">
    <a-card>
      <div style="margin-bottom: 16px; display: flex; justify-content: space-between;">
        <a-space>
          <a-button type="primary" @click="addRecord">
            新增字典类型
          </a-button>
          <a-button @click="handleRefreshCache">
            刷新缓存
          </a-button>
        </a-space>
      </div>
      
      <a-table
        row-key="id"
        :columns="columns"
        :data-source="tableData"
        :pagination="pagination"
        :loading="loading"
        @change="handleTableChange"
      >
        <template #emptyText>
          <a-empty description="暂无数据" />
        </template>
      </a-table>
    </a-card>
    
    <!-- 表单弹窗 -->
    <a-modal
      v-model:open="formModalVisible"
      :title="formModalTitle"
      :confirm-loading="loading"
      @ok="saveRecord"
      @cancel="formModalVisible = false"
    >
      <a-form
        ref="formRef"
        :model="formModel"
        :rules="formRules"
        layout="vertical"
      >
        <a-form-item label="字典名称" name="name">
          <a-input v-model:value="formModel.name" placeholder="请输入字典名称" />
        </a-form-item>
        
        <a-form-item label="字典编码" name="code">
          <a-input 
            v-model:value="formModel.code" 
            placeholder="请输入字典编码（大写字母和下划线）"
            :disabled="formModel.id"
          />
        </a-form-item>
        
        <a-form-item label="描述" name="description">
          <a-textarea 
            v-model:value="formModel.description" 
            placeholder="请输入字典描述"
            :rows="3"
          />
        </a-form-item>
        
        <a-form-item label="状态" name="status">
          <a-radio-group v-model:value="formModel.status">
            <a-radio :value="1">启用</a-radio>
            <a-radio :value="0">禁用</a-radio>
          </a-radio-group>
        </a-form-item>
        
        <a-form-item label="排序" name="sort">
          <a-input-number 
            v-model:value="formModel.sort" 
            placeholder="请输入排序值"
            :min="0"
            style="width: 100%"
          />
        </a-form-item>
      </a-form>
    </a-modal>
    
    <!-- 字典数据抽屉 -->
    <a-drawer
      v-model:open="detailDrawerVisible"
      :title="currentRecord ? `${currentRecord.name} - 数据管理` : '字典数据管理'"
      width="800"
      :footer-style="{ textAlign: 'right' }"
    >
      <dict-data-manager 
        v-if="currentRecord"
        :dict-type="currentRecord"
        @close="detailDrawerVisible = false"
      />
      
      <template #footer>
        <a-button @click="detailDrawerVisible = false">关闭</a-button>
      </template>
    </a-drawer>
  </page-container>
</template>

<style scoped>
</style>