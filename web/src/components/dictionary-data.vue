<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import {
  Badge,
  Button,
  Card,
  Form,
  Input,
  InputNumber,
  Modal,
  Popconfirm,
  Select,
  Space,
  Table,
  Tag,
  Upload,
  message,
} from 'ant-design-vue'
import {
  batchDeleteDictionaryData,
  createDictionaryData,
  deleteDictionaryData,
  exportDictionaryData,
  getDictionaryData,
  importDictionaryData,
  updateDictionaryData,
} from '@/api/dictionary.js'

const props = defineProps({
  dictionaryType: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['close'])

// 表格数据
const tableData = ref([])
const loading = ref(false)
const selectedRowKeys = ref([])
const pagination = ref({
  current: 1,
  pageSize: 10,
  total: 0,
})

// 表单相关
const formModalVisible = ref(false)
const formModalTitle = ref('')
const formRef = ref()
const formModel = ref({
  label: '',
  value: '',
  description: '',
  status: 1,
  sort: 0,
  css_class: '',
  remark: '',
})

// 搜索相关
const searchForm = ref({
  label: '',
  status: '',
})

// 表格列定义
const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
    width: 80,
  },
  {
    title: '标签',
    dataIndex: 'label',
    key: 'label',
    width: 150,
  },
  {
    title: '值',
    dataIndex: 'value',
    key: 'value',
    width: 150,
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
    ellipsis: true,
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
        text: text === 1 ? '启用' : '禁用',
      },
    ),
  },
  {
    title: '排序',
    dataIndex: 'sort',
    key: 'sort',
    width: 80,
  },
  {
    title: 'CSS类',
    dataIndex: 'css_class',
    key: 'css_class',
    width: 120,
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    key: 'created_at',
    width: 150,
    customRender: ({ text }) => text ? new Date(text).toLocaleString() : '-',
  },
  {
    title: '操作',
    key: 'actions',
    fixed: 'right',
    width: 150,
    customRender: ({ record }) => {
      return h(Space, {}, () => [
        h(
          Button,
          {
            type: 'link',
            size: 'small',
            onClick: () => editRecord(record),
          },
          () => '编辑',
        ),
        h(
          Popconfirm,
          {
            title: '确定删除该字典数据？',
            onConfirm: () => deleteRecord(record.id),
          },
          {
            default: () => h(
              Button,
              {
                type: 'link',
                size: 'small',
                danger: true,
              },
              () => '删除',
            ),
          },
        ),
      ])
    },
  },
]

// 监听字典类型变化
watch(() => props.dictionaryType, () => {
  if (props.dictionaryType) {
    fetchTableData()
  }
}, { immediate: true })

// 获取表格数据
async function fetchTableData() {
  if (!props.dictionaryType)
    return

  loading.value = true
  try {
    const res = await getDictionaryData(props.dictionaryType.code, {
      page: pagination.value.current,
      page_size: pagination.value.pageSize,
      ...searchForm.value,
    })

    if (res.code === 0) {
      tableData.value = res.data.items || []
      pagination.value.total = res.data.total || 0
    }
  }
  catch (error) {
    console.error('Failed to fetch dictionary data:', error)
    message.error('获取字典数据失败')
  }
  finally {
    loading.value = false
  }
}

// 新增记录
function addRecord() {
  formModel.value = {
    label: '',
    value: '',
    description: '',
    status: 1,
    sort: 0,
    css_class: '',
    remark: '',
  }
  formModalTitle.value = '新增字典数据'
  formModalVisible.value = true
}

// 编辑记录
function editRecord(record) {
  formModel.value = { ...record }
  formModalTitle.value = '编辑字典数据'
  formModalVisible.value = true
}

// 删除记录
async function deleteRecord(id) {
  try {
    const res = await deleteDictionaryData(props.dictionaryType.code, id)
    if (res.code === 0) {
      message.success('删除成功')
      fetchTableData()
    }
    else {
      message.error(res.msg || '删除失败')
    }
  }
  catch (error) {
    console.error('Failed to delete record:', error)
    message.error('删除失败')
  }
}

// 批量删除
async function batchDelete() {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择要删除的记录')
    return
  }

  try {
    const res = await batchDeleteDictionaryData(props.dictionaryType.code, selectedRowKeys.value)
    if (res.code === 0) {
      message.success('批量删除成功')
      selectedRowKeys.value = []
      fetchTableData()
    }
    else {
      message.error(res.msg || '批量删除失败')
    }
  }
  catch (error) {
    console.error('Failed to batch delete:', error)
    message.error('批量删除失败')
  }
}

// 保存记录
async function saveRecord() {
  try {
    await formRef.value.validateFields()

    const isEdit = !!formModel.value.id
    const res = isEdit
      ? await updateDictionaryData(props.dictionaryType.code, formModel.value.id, formModel.value)
      : await createDictionaryData(props.dictionaryType.code, formModel.value)

    if (res.code === 0) {
      message.success(isEdit ? '更新成功' : '创建成功')
      formModalVisible.value = false
      fetchTableData()
    }
    else {
      message.error(res.msg || (isEdit ? '更新失败' : '创建失败'))
    }
  }
  catch (error) {
    console.error('Failed to save record:', error)
    message.error('保存失败')
  }
}

// 搜索
function handleSearch() {
  pagination.value.current = 1
  fetchTableData()
}

// 重置搜索
function resetSearch() {
  searchForm.value = {
    label: '',
    status: '',
  }
  pagination.value.current = 1
  fetchTableData()
}

// 导出数据
async function handleExport() {
  try {
    const res = await exportDictionaryData(props.dictionaryType.code, 'excel')
    if (res.code === 0) {
      // 创建下载链接
      const link = document.createElement('a')
      link.href = res.data.url
      link.download = res.data.filename
      link.click()
      message.success('导出成功')
    }
  }
  catch (error) {
    console.error('Failed to export:', error)
    message.error('导出失败')
  }
}

// 导入数据
async function handleImport(file) {
  try {
    const res = await importDictionaryData(props.dictionaryType.code, file)
    if (res.code === 0) {
      message.success('导入成功')
      fetchTableData()
    }
    else {
      message.error(res.msg || '导入失败')
    }
  }
  catch (error) {
    console.error('Failed to import:', error)
    message.error('导入失败')
  }
}

// 表格变化处理
function handleTableChange(pag) {
  pagination.value.current = pag.current
  pagination.value.pageSize = pag.pageSize
  fetchTableData()
}

// 选择变化处理
function handleSelectionChange(keys) {
  selectedRowKeys.value = keys
}

// 表单验证规则
const formRules = {
  label: [
    { required: true, message: '请输入标签' },
  ],
  value: [
    { required: true, message: '请输入值' },
  ],
  status: [
    { required: true, message: '请选择状态' },
  ],
}
</script>

<template>
  <div>
    <!-- 搜索区域 -->
    <a-card size="small" style="margin-bottom: 16px">
      <a-form layout="inline">
        <a-form-item label="标签">
          <a-input
            v-model:value="searchForm.label"
            placeholder="请输入标签"
            allow-clear
          />
        </a-form-item>

        <a-form-item label="状态">
          <a-select
            v-model:value="searchForm.status"
            placeholder="请选择状态"
            allow-clear
            style="width: 120px"
          >
            <a-select-option :value="1">
              启用
            </a-select-option>
            <a-select-option :value="0">
              禁用
            </a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSearch">
              搜索
            </a-button>
            <a-button @click="resetSearch">
              重置
            </a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- 操作按钮 -->
    <div style="margin-bottom: 16px; display: flex; justify-content: space-between;">
      <a-space>
        <a-button type="primary" @click="addRecord">
          新增数据
        </a-button>
        <a-button
          danger
          :disabled="selectedRowKeys.length === 0"
          @click="batchDelete"
        >
          批量删除
        </a-button>
        <a-button @click="handleExport">
          导出
        </a-button>
        <a-upload
          :show-upload-list="false"
          :before-upload="() => false"
          :custom-request="({ file }) => handleImport(file)"
          accept=".xlsx,.xls,.csv"
        >
          <a-button>导入</a-button>
        </a-upload>
      </a-space>
    </div>

    <!-- 表格 -->
    <a-table
      row-key="id"
      :columns="columns"
      :data-source="tableData"
      :pagination="pagination"
      :loading="loading"
      :row-selection="{
        selectedRowKeys,
        onChange: handleSelectionChange,
      }"
      @change="handleTableChange"
    >
      <template #emptyText>
        <a-empty description="暂无数据" />
      </template>
    </a-table>

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
        <a-form-item label="标签" name="label">
          <a-input v-model:value="formModel.label" placeholder="请输入标签" />
        </a-form-item>

        <a-form-item label="值" name="value">
          <a-input v-model:value="formModel.value" placeholder="请输入值" />
        </a-form-item>

        <a-form-item label="描述" name="description">
          <a-textarea
            v-model:value="formModel.description"
            placeholder="请输入描述"
            :rows="3"
          />
        </a-form-item>

        <a-form-item label="状态" name="status">
          <a-radio-group v-model:value="formModel.status">
            <a-radio :value="1">
              启用
            </a-radio>
            <a-radio :value="0">
              禁用
            </a-radio>
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

        <a-form-item label="CSS类" name="css_class">
          <a-input v-model:value="formModel.css_class" placeholder="请输入CSS类名" />
        </a-form-item>

        <a-form-item label="备注" name="remark">
          <a-textarea
            v-model:value="formModel.remark"
            placeholder="请输入备注"
            :rows="2"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<style scoped>
</style>
