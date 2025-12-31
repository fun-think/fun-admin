<script setup>
import { message } from 'ant-design-vue'
import {
  createDictionaryType,
  deleteDictionaryType,
  getDictionaryTypes,
  refreshDictionaryCache,
  updateDictionaryType,
} from '@/api/dictionary.js'

// 表格数据
const tableData = ref([])
const loading = ref(false)
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  onChange(current, pageSize) {
    pagination.pageSize = pageSize
    pagination.current = current
    fetchTableData()
  },
})

// 表单相关
const formModalVisible = ref(false)
const formModalTitle = ref('')
const formRef = ref()
const formModel = reactive({
  name: '',
  code: '',
  description: '',
  status: 1,
  sort: 0,
})

// 详情相关
const detailDrawerVisible = ref(false)
const currentRecord = ref(null)

// 表格列定义
const columns = shallowRef([
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
    width: 80,
  },
  {
    title: '名称',
    dataIndex: 'name',
    key: 'name',
    width: 150,
  },
  {
    title: '编码',
    dataIndex: 'code',
    key: 'code',
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
  },
  {
    title: '排序',
    dataIndex: 'sort',
    key: 'sort',
    width: 80,
  },
  {
    title: '数据量',
    dataIndex: 'data_count',
    key: 'data_count',
    width: 100,
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    key: 'created_at',
    width: 150,
  },
  {
    title: '操作',
    key: 'actions',
    fixed: 'right',
    width: 200,
  },
])

// 获取表格数据
async function fetchTableData() {
  if (loading.value)
    return
  loading.value = true
  try {
    const res = await getDictionaryTypes({
      page: pagination.current,
      page_size: pagination.pageSize,
    })

    if (res.code === 0) {
      tableData.value = res.data.items || []
      pagination.total = res.data.total || 0
    }
  }
  catch (error) {
    console.error('Failed to fetch dictionary types:', error)
    message.error('获取字典类型失败')
  }
  finally {
    loading.value = false
  }
}

// 新增记录
function addRecord() {
  Object.assign(formModel, {
    name: '',
    code: '',
    description: '',
    status: 1,
    sort: 0,
  })
  formModalTitle.value = '新增字典类型'
  formModalVisible.value = true
}

// 编辑记录
function editRecord(record) {
  Object.assign(formModel, record)
  formModalTitle.value = '编辑字典类型'
  formModalVisible.value = true
}

// 查看字典数据
function viewDictionaryData(record) {
  currentRecord.value = record
  detailDrawerVisible.value = true
}

// 删除记录
async function deleteRecord(id) {
  try {
    const res = await deleteDictionaryType(id)
    if (res.code === 0) {
      message.success('删除成功')
      await fetchTableData()
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

// 保存记录
async function saveRecord() {
  const close = message.loading('提交中......')
  try {
    await formRef.value.validateFields()

    let res = {}
    if (formModel.id) {
      res = await updateDictionaryType(formModel.id, formModel)
    }
    else {
      res = await createDictionaryType(formModel)
    }

    if (res.code === 0) {
      message.success(formModel.id ? '更新成功' : '创建成功')
      formModalVisible.value = false
      await fetchTableData()
    }
    else {
      message.error(res.msg || (formModel.id ? '更新失败' : '创建失败'))
    }
  }
  catch (error) {
    console.error('Failed to save record:', error)
    message.error('保存失败')
  }
  finally {
    close()
  }
}

// 刷新缓存
async function handleRefreshCache() {
  try {
    const res = await refreshDictionaryCache()
    if (res.code === 0) {
      message.success('缓存刷新成功')
    }
    else {
      message.error(res.msg || '缓存刷新失败')
    }
  }
  catch (error) {
    console.error('Failed to refresh cache:', error)
    message.error('缓存刷新失败')
  }
}

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入字典名称' },
  ],
  code: [
    { required: true, message: '请输入字典编码' },
    { pattern: /^[A-Z_]+$/, message: '字典编码只能包含大写字母和下划线' },
  ],
  status: [
    { required: true, message: '请选择状态' },
  ],
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
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-badge :status="record.status === 1 ? 'success' : 'error'" :text="record.status === 1 ? '启用' : '禁用'" />
          </template>
          <template v-else-if="column.key === 'created_at'">
            {{ record.created_at ? new Date(record.created_at).toLocaleString() : '-' }}
          </template>
          <template v-else-if="column.key === 'actions'">
            <a-space>
              <a-button type="link" size="small" @click="viewDictionaryData(record)">
                数据管理
              </a-button>
              <a-button type="link" size="small" @click="editRecord(record)">
                编辑
              </a-button>
              <a-popconfirm
                title="确定删除该字典类型？"
                @confirm="deleteRecord(record.id)"
              >
                <a-button type="link" size="small" danger>
                  删除
                </a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>

        <template #emptyText>
          <a-empty description="暂无数据" />
        </template>
      </a-table>
    </a-card>

    <!-- 表单弹窗 -->
    <a-modal
      v-model:open="formModalVisible"
      :title="formModalTitle"
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
      </a-form>
    </a-modal>

    <!-- 字典数据抽屉 -->
    <a-drawer
      v-model:open="detailDrawerVisible"
      :title="currentRecord ? `${currentRecord.name} - 数据管理` : '字典数据管理'"
      width="800"
      :footer-style="{ textAlign: 'right' }"
    >
      <dictionary-data
        v-if="currentRecord"
        :dictionary-type="currentRecord"
        @close="detailDrawerVisible = false"
      />

      <template #footer>
        <a-button @click="detailDrawerVisible = false">
          关闭
        </a-button>
      </template>
    </a-drawer>
  </page-container>
</template>

<style scoped>
</style>
