<script setup>
import { ref, computed, onMounted, h } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { message, Tag } from 'ant-design-vue'
import CrudTableModal from './components/crud-table-modal.vue'
import {
  createResourceRecord,
  deleteResourceRecord,
  getResourceConfig,
  getResourceData,
  runResourceAction,
  updateResourceRecord,
} from '@/api/resources.js'

const resourceSlug = 'crud_items'
const columns = ref([])
const resourceMetaFields = ref([])
const actionColumn = {
  title: '操作',
  dataIndex: 'action',
  key: 'action',
  width: 150,
  fixed: 'right',
}

const state = ref({
  loading: false,
  dataSource: [],
  queryParams: {
    name: '',
    value: '',
    remark: '',
  },
  pagination: {
    current: 1,
    pageSize: 10,
    total: 0,
  },
})

const crudTableModal = ref()
const editingRecord = ref(null)
const resourceMeta = ref(null)
const selectedRowKeys = ref([])
const rowSelection = computed(() => ({
  selectedRowKeys: selectedRowKeys.value,
  onChange: keys => {
    selectedRowKeys.value = keys
  },
}))
const actionButtons = computed(() => {
  const actions = resourceMeta.value?.actions || []
  return actions.filter(action => ['reset_values', 'bulk_delete'].includes(action.name))
})
const actionLoading = ref(false)

async function query() {
  state.value.loading = true
  try {
    const { queryParams, pagination } = state.value
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
    }

    Object.entries(queryParams).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        params[`search_${key}`] = value
      }
    })
    const res = await getResourceData(resourceSlug, params)
    if (res?.code !== 0) {
      message.error(res?.message || '查询失败')
      return
    }
    const payload = res.data || {}
    state.value.dataSource = payload.items || []
    state.value.pagination.total = payload.total || 0
    state.value.pagination.current = payload.page || pagination.current
    state.value.pagination.pageSize = payload.page_size || pagination.pageSize
  }
  catch (error) {
    console.error('Failed to query data:', error)
    message.error('查询失败')
  }
  finally {
    state.value.loading = false
  }
}

function initQuery() {
  state.value.pagination.current = 1
  query()
}

function resetQuery() {
  state.value.queryParams = {
    name: '',
    value: '',
    remark: '',
  }
  state.value.pagination.current = 1
  query()
}

function handleTableChange(pagination) {
  state.value.pagination.current = pagination.current || state.value.pagination.current
  state.value.pagination.pageSize = pagination.pageSize || state.value.pagination.pageSize
  query()
}

async function handleDelete(record) {
  if (!record?.id) {
    return message.error('id 不能为空')
  }
  try {
    const res = await deleteResourceRecord(resourceSlug, record.id)
    if (res?.code === 0) {
      message.success('删除成功')
      await query()
    }
    else {
      message.error(res?.message || '删除失败')
    }
  }
  catch (error) {
    console.error('Failed to delete record:', error)
    message.error('删除失败')
  }
}

function handleAdd() {
  editingRecord.value = null
  crudTableModal.value?.open()
}

function handleEdit(record) {
  editingRecord.value = record
  crudTableModal.value?.open(record)
}

async function handleModalOk(formData) {
  try {
    if (editingRecord.value?.id) {
      const res = await updateResourceRecord(resourceSlug, editingRecord.value.id, formData)
      if (res?.code !== 0) {
        return message.error(res?.message || '更新失败')
      }
      message.success('更新成功')
    }
    else {
      const res = await createResourceRecord(resourceSlug, formData)
      if (res?.code !== 0) {
        return message.error(res?.message || '新增失败')
      }
      message.success('新增成功')
    }
    await query()
  }
  catch (error) {
    console.error('Failed to save record:', error)
    message.error('保存失败')
  }
  finally {
    editingRecord.value = null
  }
}

function handleModalCancel() {
  editingRecord.value = null
}

async function handleAction(actionName, label) {
  if (selectedRowKeys.value.length === 0) {
    return message.warning('请先选择至少一条记录')
  }
  actionLoading.value = true
  try {
    const res = await runResourceAction(resourceSlug, actionName, { ids: selectedRowKeys.value })
    if (res?.code === 0) {
      message.success(`${label} 成功`)
      await query()
      selectedRowKeys.value = []
    }
    else {
      message.error(res?.message || `${label} 失败`)
    }
  }
  catch (error) {
    console.error(`Failed to run action ${actionName}:`, error)
    message.error(`${label} 失败`)
  }
  finally {
    actionLoading.value = false
  }
}

function buildColumnsFromMeta(meta) {
  if (!Array.isArray(meta)) {
    columns.value = [actionColumn]
    return
  }
  columns.value = meta.map(item => {
    const column = {
      title: item.label,
      dataIndex: item.name,
      key: item.name,
      sorter: item.sortable,
      align: item.align?.toLowerCase?.() || 'left',
      width: item.width,
    }
    if (item.type === 'datetime' || item.type === 'date') {
      column.customRender = ({ text }) => (text ? new Date(text).toLocaleString() : '-')
    }
    if (item.type === 'badge' && item.badgeMap) {
      column.customRender = ({ text }) => h(Tag, { color: 'blue' }, () => item.badgeMap[text] || text)
    }
    return column
  })
  columns.value.push(actionColumn)
}

async function loadResourceMeta() {
  try {
    const res = await getResourceConfig(resourceSlug, { language: localStorage.getItem('admin-language') || 'zh-CN' })
    if (res?.code !== 0) {
      return
    }
    resourceMeta.value = res.data
    resourceMetaFields.value = res.data?.fields || []
    buildColumnsFromMeta(res.data?.columns)
  }
  catch (error) {
    console.error('Failed to load resource metadata:', error)
  }
}

onMounted(async () => {
  await loadResourceMeta()
  initQuery()
})
</script>

<template>
  <page-container>
    <a-card mb-4>
      <a-form class="system-crud-wrapper" :label-col="{ span: 7 }" :model="state.queryParams">
        <a-row :gutter="[15, 0]">
          <a-col :span="6">
            <a-form-item name="name" label="名">
              <a-input v-model:value="state.queryParams.name" placeholder="请输入名" />
            </a-form-item>
          </a-col>
          <a-col :span="6">
            <a-form-item name="value" label="值">
              <a-input v-model:value="state.queryParams.value" placeholder="请输入值" />
            </a-form-item>
          </a-col>
          <a-col :span="6">
            <a-form-item name="remark" label="备注">
              <a-input v-model:value="state.queryParams.remark" placeholder="请输入备注" />
            </a-form-item>
          </a-col>
          <a-col :span="6">
            <a-space flex justify-end w-full>
              <a-button :loading="state.loading" type="primary" @click="initQuery">
                查询
              </a-button>
              <a-button :loading="state.loading" @click="resetQuery">
                重置
              </a-button>
            </a-space>
          </a-col>
        </a-row>
      </a-form>
    </a-card>

    <a-card title="增删改查表格">
      <template #extra>
        <a-space size="middle">
          <a-button type="primary" @click="handleAdd">
            <template #icon>
              <PlusOutlined />
            </template>
            新增
          </a-button>
          <a-button
            v-for="action in actionButtons"
            :key="action.name"
            :loading="actionLoading"
            :danger="action.name === 'bulk_delete'"
            @click="handleAction(action.name, action.label)"
          >
            {{ action.label }}
          </a-button>
        </a-space>
      </template>
      <a-table
        row-key="id"
        :loading="state.loading"
        :columns="columns"
        :data-source="state.dataSource"
        :pagination="state.pagination"
        :row-selection="rowSelection"
        @change="handleTableChange"
      >
        <template #bodyCell="scope">
          <template v-if="scope?.column?.dataIndex === 'action'">
            <div flex gap-2>
              <a-button type="link" @click="handleEdit(scope?.record)">
                编辑
              </a-button>
              <a-popconfirm
                title="确定删除该条数据？" ok-text="确定" cancel-text="取消"
                @confirm="() => handleDelete(scope?.record)"
              >
                <a-button type="link">
                  删除
                </a-button>
              </a-popconfirm>
            </div>
          </template>
        </template>
      </a-table>
    </a-card>

    <CrudTableModal
      ref="crudTableModal"
      :fields="resourceMetaFields"
      @ok="handleModalOk"
      @cancel="handleModalCancel"
    />
  </page-container>
</template>

<style lang="less" scoped>
.system-crud-wrapper{
    .ant-form-item{
      margin: 0;
    }
  }
</style>
