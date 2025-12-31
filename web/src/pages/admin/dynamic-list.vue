<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  deleteResourceRecord,
  deleteResourceRecords,
  getResourceConfig,
  getResourceData,
  runResourceAction,
} from '@/api/resources.js'
import { useAccess } from '@/composables/access.js'

const route = useRoute()
const router = useRouter()

// 获取资源标识符
const resourceSlug = route.params.slug

// 表格相关状态
const columns = ref([])
const dataSource = ref([])
const loading = ref(false)
const pagination = ref({
  current: 1,
  pageSize: 10,
  total: 0,
})
// 排序
const orderBy = ref('')
const orderDirection = ref('')
// 表格尺寸
const tableSize = ref('middle') // default | middle | small
// 列可见性
const allColumns = ref([])
const visibleKeys = ref([])

// 搜索和过滤相关状态
const searchForm = reactive({})
const filterForm = reactive({})
// 保存的筛选器与 Chips
const savedFilters = ref([])
const savedFilterName = ref('')
const selectedSavedFilter = ref('')
const activeChips = ref([])
// 快速搜索关键词
const quickKeyword = ref('')

// 批量操作相关状态
const selectedRowKeys = ref([])
const selectedRows = ref([])

// 资源配置
const resourceConfig = ref(null)
// 白名单驱动的搜索/过滤定义
const searchableNames = ref([])
const filterDefs = ref([])

// 获取当前语言
function getCurrentLanguage() {
  // 这里可以从 localStorage 或其他地方获取当前语言设置
  return localStorage.getItem('admin-language') || 'zh-CN'
}

// 获取资源配置
async function fetchResourceConfig() {
  // 仪表板不需要资源配置
  if (resourceSlug === 'dashboard') {
    return
  }

  try {
    const res = await getResourceConfig(resourceSlug, { language: getCurrentLanguage() })
    if (res.code === 0) {
      resourceConfig.value = res.data
      // 根据字段配置生成表格列（优先使用后端 columns 元信息）
      if (res.data.columns && res.data.columns.length) {
        generateColumnsFromColumns(res.data.columns)
      }
      else {
        generateColumns(res.data.fields)
      }
      // 基于白名单构建搜索与过滤定义
      buildSearchAndFilters()
      // 读取已保存筛选器
      loadSavedFilters()
      // 获取数据
      fetchData()
    }
    else {
      message.error(res.msg || '获取资源配置失败')
    }
  }
  catch (error) {
    console.error('Failed to fetch resource config:', error)
    message.error('获取资源配置失败')
  }
}

// 基于后端白名单构建搜索与过滤 UI 定义
function buildSearchAndFilters() {
  const fields = resourceConfig.value?.fields || []
  const fieldMap = Object.fromEntries(fields.map(f => [f.name, f]))
  // 搜索字段：优先使用后端 searchable_fields
  if (Array.isArray(resourceConfig.value?.searchable_fields) && resourceConfig.value.searchable_fields.length) {
    searchableNames.value = resourceConfig.value.searchable_fields.filter(n => fieldMap[n])
  }
  else {
    searchableNames.value = fields.filter(f => ['text', 'email', 'textarea'].includes(f.type)).map(f => f.name)
  }
  // 过滤字段：优先使用后端 filters 元信息
  if (Array.isArray(resourceConfig.value?.filters) && resourceConfig.value.filters.length) {
    filterDefs.value = resourceConfig.value.filters
  }
  else {
    filterDefs.value = fields
      .filter(f => ['select', 'boolean'].includes(f.type))
      .map(f => ({ name: f.name, label: f.label, type: f.type, options: f.options || [] }))
  }
  // 初次构建 chips
  buildActiveChips()
}

// 计算当前激活的 chips
function buildActiveChips() {
  const fields = resourceConfig.value?.fields || []
  const labelOf = name => fields.find(f => f.name === name)?.label || name
  const chips = []
  Object.entries(searchForm).forEach(([k, v]) => {
    if (v !== undefined && v !== null && v !== '')
      chips.push({ key: `search_${k}`, label: `${labelOf(k)}:${v}`, type: 'search', name: k })
  })
  Object.entries(filterForm).forEach(([k, v]) => {
    if (v !== undefined && v !== null && v !== '')
      chips.push({ key: `filter_${k}`, label: `${labelOf(k)}:${v}`, type: 'filter', name: k })
  })
  activeChips.value = chips
}

// 本地存储 key
const filterStorageKey = () => `admin:filters:${resourceSlug}`

// 读取/保存筛选器
function loadSavedFilters() {
  try {
    const raw = localStorage.getItem(filterStorageKey())
    savedFilters.value = raw ? JSON.parse(raw) : []
  }
  catch {
    savedFilters.value = []
  }
}
function saveCurrentFilter() {
  if (!savedFilterName.value.trim()) {
    message.warning('请输入筛选器名称')
    return
  }
  const payload = { name: savedFilterName.value.trim(), search: { ...searchForm }, filter: { ...filterForm } }
  const list = savedFilters.value.filter(f => f.name !== payload.name)
  list.unshift(payload)
  savedFilters.value = list
  localStorage.setItem(filterStorageKey(), JSON.stringify(list))
  message.success('已保存筛选器')
}
function applySavedFilter(name) {
  const found = savedFilters.value.find(f => f.name === name)
  if (!found)
    return
  Object.keys(searchForm).forEach(k => delete searchForm[k])
  Object.keys(filterForm).forEach(k => delete filterForm[k])
  Object.assign(searchForm, found.search || {})
  Object.assign(filterForm, found.filter || {})
  selectedSavedFilter.value = name
  buildActiveChips()
  fetchData(1)
}
function deleteSavedFilter(name) {
  savedFilters.value = savedFilters.value.filter(f => f.name !== name)
  localStorage.setItem(filterStorageKey(), JSON.stringify(savedFilters.value))
  if (selectedSavedFilter.value === name)
    selectedSavedFilter.value = ''
}

// 根据字段配置生成表格列
function generateColumns(fields) {
  const cols = fields.map((field) => {
    const col = {
      title: field.label,
      dataIndex: field.name,
      key: field.name,
    }

    // 根据字段类型处理显示
    switch (field.type) {
      case 'boolean':
        col.customRender = ({ text }) => text ? '是' : '否'
        break
      case 'date':
      case 'datetime':
        col.customRender = ({ text }) => text ? new Date(text).toLocaleString() : ''
        break
      case 'relationship':
        // 关联字段显示关联数据
        col.customRender = ({ record }) => {
          const relatedData = record[`${field.name}_data`]
          if (relatedData) {
            return relatedData[field.display_field || 'name'] || relatedData.id
          }
          return record[field.name] || ''
        }
        break
    }

    return col
  })

  // 添加操作列
  cols.push({
    title: '操作',
    dataIndex: 'action',
    key: 'action',
  })

  columns.value = cols
  allColumns.value = cols
  visibleKeys.value = cols.filter(c => c.dataIndex !== 'action').map(c => c.dataIndex)
}

// 根据后端列元信息构建 a-table 列
function generateColumnsFromColumns(colsMeta) {
  const cols = colsMeta.map((col) => {
    const column = {
      title: col.label,
      dataIndex: col.name,
      key: col.name,
      align: col.align || 'left',
      sorter: !!col.sortable,
    }
    if (col.width)
      column.width = col.width
    if (col.sticky)
      column.fixed = col.sticky
    // 通用渲染：布尔/日期/枚举/徽章/图片/头像/链接/自定义格式化
    column.customRender = ({ text, record }) => {
      if (col.badgeMap && text in col.badgeMap) {
        return h('span', { style: { display: 'inline-block', padding: '0 8px', borderRadius: '4px', backgroundColor: col.badgeMap[text], color: '#fff' } }, col.enumMap?.[text] || String(text))
      }
      if (col.enumMap && text in col.enumMap)
        return col.enumMap[text]
      if (col.type === 'boolean')
        return text ? '是' : '否'
      if (col.type === 'date' || col.type === 'datetime')
        return text ? new Date(text).toLocaleString() : ''
      if (col.type === 'image' || col.type === 'avatar') {
        const url = typeof text === 'string' ? text : ''
        if (!url)
          return ''
        return h('img', { src: url, style: { width: col.type === 'avatar' ? '32px' : '72px', height: col.type === 'avatar' ? '32px' : 'auto', borderRadius: col.type === 'avatar' ? '50%' : '4px', objectFit: 'cover' } })
      }
      if (col.type === 'link') {
        const url = col.urlField ? record[col.urlField] : (typeof text === 'string' ? text : '')
        return url ? h('a', { href: url, target: '_blank' }, text || '查看') : ''
      }
      if (col.formatter) {
        if (col.formatter === 'datetimeFromNow' && text)
          return new Date(text).toLocaleString()
      }
      return text ?? ''
    }
    return column
  })
  // 操作列
  cols.push({ title: '操作', dataIndex: 'action', key: 'action', fixed: 'right' })
  columns.value = cols
  allColumns.value = cols
  visibleKeys.value = cols.filter(c => c.dataIndex !== 'action').map(c => c.dataIndex)
}

// 获取数据
async function fetchData(page = 1) {
  loading.value = true
  try {
    // 构建查询参数
    const params = {
      page,
      page_size: pagination.value.pageSize,
      language: getCurrentLanguage(),
      order_by: orderBy.value || resourceConfig.value?.default_order?.field,
      order_direction: orderDirection.value || resourceConfig.value?.default_order?.direction,
    }

    // 添加过滤参数
    Object.keys(filterForm).forEach((key) => {
      if (filterForm[key] !== undefined && filterForm[key] !== null && filterForm[key] !== '') {
        params[key] = filterForm[key]
      }
    })

    // 添加搜索参数
    Object.keys(searchForm).forEach((key) => {
      if (searchForm[key] !== undefined && searchForm[key] !== null && searchForm[key] !== '') {
        params[`search_${key}`] = searchForm[key]
      }
    })

    const res = await getResourceData(resourceSlug, params)

    if (res.code === 0) {
      dataSource.value = res.data.items || []
      pagination.value.total = res.data.total || 0
      pagination.value.current = page
    }
    else {
      message.error(res.msg || '获取数据失败')
    }
  }
  catch (error) {
    console.error('Failed to fetch data:', error)
    message.error('获取数据失败')
  }
  finally {
    loading.value = false
  }
}

// 处理分页变化
function handlePageChange(page) {
  fetchData(page)
}

// 表格变化（分页、排序）
function handleTableChange(pager, _filters, sorter) {
  const current = pager?.current || 1
  const size = pager?.pageSize || pagination.value.pageSize
  pagination.value.current = current
  pagination.value.pageSize = size
  if (sorter && sorter.order) {
    orderBy.value = sorter.field || sorter.columnKey
    orderDirection.value = sorter.order === 'ascend' ? 'ASC' : 'DESC'
  }
  else {
    orderBy.value = ''
    orderDirection.value = ''
  }
  fetchData(current)
}

// 应用列可见性
function applyVisibleColumns() {
  const keys = new Set(visibleKeys.value)
  columns.value = allColumns.value.filter(c => c.dataIndex === 'action' || keys.has(c.dataIndex))
}

// 处理删除
async function handleDelete(record) {
  try {
    const res = await deleteResourceRecord(resourceSlug, record.id)

    if (res.code === 0) {
      message.success('删除成功')
      fetchData(pagination.value.current)
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

// 处理编辑
function handleEdit(record) {
  // 跳转到编辑页面
  router.push(`/admin/${resourceSlug}/edit/${record.id}`)
}

// 处理新增
function handleAdd() {
  // 跳转到新增页面
  router.push(`/admin/${resourceSlug}/create`)
}

// 重置搜索表单
function resetSearchForm() {
  Object.keys(searchForm).forEach((key) => {
    searchForm[key] = ''
  })
  quickKeyword.value = ''
  buildActiveChips()
  fetchData()
}

// 重置过滤表单
function resetFilterForm() {
  Object.keys(filterForm).forEach((key) => {
    filterForm[key] = ''
  })
  buildActiveChips()
  fetchData()
}

// 搜索
function handleSearch() {
  if (quickKeyword.value && searchableNames.value.length) {
    searchableNames.value.forEach((name) => {
      searchForm[name] = quickKeyword.value
    })
  }
  buildActiveChips()
  fetchData()
}

// 导出数据
function handleExport(format) {
  // 构建导出参数
  const params = new URLSearchParams()

  // 添加过滤参数
  Object.keys(filterForm).forEach((key) => {
    if (filterForm[key] !== undefined && filterForm[key] !== null && filterForm[key] !== '') {
      params.append(key, filterForm[key])
    }
  })

  // 添加搜索参数
  Object.keys(searchForm).forEach((key) => {
    if (searchForm[key] !== undefined && searchForm[key] !== null && searchForm[key] !== '') {
      params.append(`search_${key}`, searchForm[key])
    }
  })

  // 添加语言参数
  params.append('language', getCurrentLanguage())

  // 添加导出格式
  params.append('format', format)

  // 构建导出 URL
  const exportUrl = `/api/v1/${resourceSlug}/export?${params.toString()}`

  // 创建下载链接并触发下载
  const link = document.createElement('a')
  link.href = exportUrl
  link.download = `${resourceConfig.value.title}_导出_${new Date().toISOString().slice(0, 19).replace(/:/g, '-')}.${format}`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

// 批量删除
async function handleBatchDelete() {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请先选择要删除的记录')
    return
  }

  try {
    const res = await deleteResourceRecords(resourceSlug, selectedRowKeys.value)

    if (res.code === 0) {
      message.success('批量删除成功')
      // 清空选择
      selectedRowKeys.value = []
      selectedRows.value = []
      // 重新加载数据
      fetchData(pagination.value.current)
    }
    else {
      message.error(res.msg || '批量删除失败')
    }
  }
  catch (error) {
    console.error('Failed to batch delete records:', error)
    message.error('批量删除失败')
  }
}

// 表格行选择配置
const rowSelection = {
  onChange: (selectedRowKeysValue, selectedRowsValue) => {
    selectedRowKeys.value = selectedRowKeysValue
    selectedRows.value = selectedRowsValue
  },
}

// 触发动作（无参或批量）
async function triggerAction(action) {
  // 若是批量动作必须选择行
  const ids = action.bulk ? selectedRowKeys.value : []
  if (action.bulk && ids.length === 0) {
    message.warning('请先选择记录')
    return
  }
  await runAction(action, null, ids)
}

// 动作参数对话框
const actionModalVisible = ref(false)
const actionModalDef = ref(null)
const actionModalModel = reactive({})
const actionModalIds = ref([])

function openActionModal(action, ids) {
  actionModalDef.value = action
  actionModalIds.value = ids
  // 初始化参数模型
  Object.keys(actionModalModel).forEach(k => delete actionModalModel[k])
  ;(action.form_fields || []).forEach((f) => {
    actionModalModel[f.name] = undefined
  })
  actionModalVisible.value = true
}
async function submitActionModal() {
  const params = { ...actionModalModel }
  const res = await runResourceAction(resourceSlug, actionModalDef.value.name, { ids: actionModalIds.value, params })
  if (res.code === 0) {
    message.success('操作成功')
    actionModalVisible.value = false
    fetchData(pagination.value.current)
  }
  else {
    message.error(res.message || '操作失败')
  }
}

// 运行动作（行级或批量）
async function runAction(action, record, idsOverride) {
  const ids = idsOverride ?? (record ? [record.id] : [])
  if (action.form_fields && action.form_fields.length) {
    openActionModal(action, ids)
    return
  }
  const res = await runResourceAction(resourceSlug, action.name, { ids, params: {} })
  if (res.code === 0) {
    message.success('操作成功')
    fetchData(pagination.value.current)
  }
  else {
    message.error(res.message || '操作失败')
  }
}

const { hasPermission } = useAccess()
function canShowAction(act) {
  return hasPermission(act.permission)
}

// 渲染动作表单字段控件（与表单页简化对齐）
function getFormComponent(field) {
  switch (field.type) {
    case 'boolean':
      return 'a-switch'
    case 'number':
      return 'a-input-number'
    case 'select':
      return 'a-select'
    case 'date':
    case 'datetime':
      return 'a-date-picker'
    default:
      return 'a-input'
  }
}

onMounted(() => {
  fetchResourceConfig()
})
</script>

<template>
  <page-container :title="resourceConfig?.title">
    <a-card>
      <!-- 搜索和过滤区域 -->
      <a-form layout="inline" class="mb-4">
        <!-- 搜索字段 -->
        <!-- 快速搜索输入 -->
        <a-form-item>
          <a-input-search v-model:value="quickKeyword" placeholder="快速搜索" enter-button @search="handleSearch" />
        </a-form-item>
        <!-- 精确搜索字段（可选显示） -->
        <template v-for="name in searchableNames" :key="name">
          <a-form-item v-show="false" :label="(resourceConfig?.fields || []).find(f => f.name === name)?.label">
            <a-input v-model:value="searchForm[name]" :placeholder="`搜索${(resourceConfig?.fields || []).find(f => f.name === name)?.label}`" allow-clear />
          </a-form-item>
        </template>

        <!-- 过滤字段 -->
        <template v-for="f in filterDefs" :key="f.name">
          <a-form-item :label="f.label">
            <template v-if="f.type === 'select' && f.options && f.options.length">
              <a-select
                v-model:value="filterForm[f.name]"
                :placeholder="`筛选${f.label}`"
                style="width: 160px"
                allow-clear
              >
                <a-select-option
                  v-for="opt in f.options"
                  :key="opt.value"
                  :value="opt.value"
                >
                  {{ opt.label }}
                </a-select-option>
              </a-select>
            </template>
            <template v-else-if="f.type === 'boolean'">
              <a-select
                v-model:value="filterForm[f.name]"
                :placeholder="`筛选${f.label}`"
                style="width: 120px"
                allow-clear
              >
                <a-select-option :value="true">
                  是
                </a-select-option>
                <a-select-option :value="false">
                  否
                </a-select-option>
              </a-select>
            </template>
            <template v-else>
              <a-input v-model:value="filterForm[f.name]" :placeholder="`筛选${f.label}`" allow-clear style="width: 160px" />
            </template>
          </a-form-item>
        </template>

        <a-form-item>
          <a-button type="primary" @click="handleSearch">
            搜索
          </a-button>
          <a-button style="margin-left: 8px" @click="resetSearchForm">
            重置
          </a-button>

          <!-- 导出按钮 -->
          <a-dropdown v-if="resourceConfig?.exportable !== false" style="margin-left: 8px">
            <template #overlay>
              <a-menu>
                <a-menu-item key="csv" @click="handleExport('csv')">
                  导出 CSV
                </a-menu-item>
                <a-menu-item key="xlsx" @click="handleExport('xlsx')">
                  导出 Excel
                </a-menu-item>
              </a-menu>
            </template>
            <a-button>
              导出 <DownOutlined />
            </a-button>
          </a-dropdown>

          <!-- 语言切换 -->
          <a-dropdown style="margin-left: 8px">
            <template #overlay>
              <a-menu @click="({ key }) => { localStorage.setItem('admin-language', key); location.reload(); }">
                <a-menu-item key="zh-CN">
                  中文
                </a-menu-item>
                <a-menu-item key="en">
                  English
                </a-menu-item>
              </a-menu>
            </template>
            <a-button>
              语言 <DownOutlined />
            </a-button>
          </a-dropdown>

          <!-- 保存筛选器 -->
          <a-dropdown style="margin-left: 8px">
            <template #overlay>
              <a-menu style="width: 260px; padding: 8px 12px">
                <div style="font-weight: 500; margin-bottom: 8px">
                  保存当前筛选
                </div>
                <a-input v-model:value="savedFilterName" placeholder="筛选器名称" />
                <a-button type="primary" size="small" style="margin-top: 8px" @click="saveCurrentFilter">
                  保存
                </a-button>
                <div style="font-weight: 500; margin: 12px 0 8px">
                  已保存
                </div>
                <a-empty v-if="!savedFilters.length" description="暂无" />
                <div v-else>
                  <div v-for="it in savedFilters" :key="it.name" style="display:flex; justify-content:space-between; align-items:center; margin-bottom:6px">
                    <a @click="applySavedFilter(it.name)">{{ it.name }}</a>
                    <a-button type="link" danger size="small" @click.stop="deleteSavedFilter(it.name)">
                      删除
                    </a-button>
                  </div>
                </div>
              </a-menu>
            </template>
            <a-button>筛选器</a-button>
          </a-dropdown>

          <!-- 列设置 -->
          <a-dropdown style="margin-left: 8px">
            <template #overlay>
              <a-menu style="width: 220px; padding: 8px 12px">
                <div style="margin-bottom: 8px; font-weight: 500">
                  列可见性
                </div>
                <a-checkbox-group v-model:value="visibleKeys" @change="applyVisibleColumns">
                  <a-row :gutter="[8, 8]">
                    <a-col v-for="c in allColumns.filter(x => x.dataIndex !== 'action')" :key="c.key" :span="24">
                      <a-checkbox :value="c.dataIndex">
                        {{ c.title }}
                      </a-checkbox>
                    </a-col>
                  </a-row>
                </a-checkbox-group>
                <div style="margin: 12px 0 8px; font-weight: 500">
                  密度
                </div>
                <a-radio-group v-model:value="tableSize">
                  <a-radio-button value="small">
                    紧凑
                  </a-radio-button>
                  <a-radio-button value="middle">
                    中等
                  </a-radio-button>
                  <a-radio-button value="default">
                    宽松
                  </a-radio-button>
                </a-radio-group>
              </a-menu>
            </template>
            <a-button>列设置</a-button>
          </a-dropdown>
        </a-form-item>
      </a-form>

      <!-- 筛选器徽章（chips） -->
      <div v-if="activeChips.length" style="margin: -8px 0 8px 0">
        <a-space wrap>
          <a-tag v-for="chip in activeChips" :key="chip.key" closable @close="() => { if (chip.type === 'search') searchForm[chip.name] = ''; else filterForm[chip.name] = ''; buildActiveChips(); fetchData(1); }">
            {{ chip.label }}
          </a-tag>
        </a-space>
      </div>

      <template #extra>
        <a-space>
          <a-button v-if="resourceConfig?.creatable" type="primary" @click="handleAdd">
            新增
          </a-button>
          <!-- 头部动作（示例：批量动作） -->
          <a-dropdown v-if="(resourceConfig?.actions || []).some(a => canShowAction(a))">
            <template #overlay>
              <a-menu>
                <a-menu-item v-for="act in resourceConfig.actions.filter(a => canShowAction(a))" :key="act.name" @click="triggerAction(act)">
                  <span>{{ act.label }}</span>
                </a-menu-item>
              </a-menu>
            </template>
            <a-button>操作 <DownOutlined /></a-button>
          </a-dropdown>
        </a-space>
      </template>

      <a-table
        v-if="resourceSlug !== 'dashboard'"
        row-key="id"
        :loading="loading"
        :columns="columns"
        :data-source="dataSource"
        :size="tableSize"
        :pagination="{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: pagination.total,
          onChange: handlePageChange,
          showSizeChanger: true,
          pageSizeOptions: ['10', '20', '50', '100'],
        }"
        :row-selection="rowSelection"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'action'">
            <div class="flex gap-2">
              <template v-for="act in (resourceConfig?.actions || []).filter(a => canShowAction(a))" :key="act.name">
                <template v-if="canShowAction(act)">
                  <a-popconfirm v-if="act.confirm" :title="act.confirm" ok-text="确定" cancel-text="取消" @confirm="() => runAction(act, record)">
                    <a-button type="link">
                      {{ act.label }}
                    </a-button>
                  </a-popconfirm>
                  <a-button v-else type="link" @click="runAction(act, record)">
                    {{ act.label }}
                  </a-button>
                </template>
              </template>
            </div>
          </template>
        </template>
      </a-table>

      <!-- 动作参数对话框 -->
      <a-modal v-model:open="actionModalVisible" :title="actionModalDef?.label || '动作'" @ok="submitActionModal">
        <a-form label-col="{ span: 6 }" wrapper-col="{ span: 16 }">
          <a-form-item v-for="f in (actionModalDef?.form_fields || [])" :key="f.name" :label="f.label">
            <component :is="getFormComponent(f)" v-model:value="actionModalModel[f.name]" :placeholder="`请输入${f.label}`" />
          </a-form-item>
        </a-form>
      </a-modal>

      <!-- 批量操作栏 -->
      <div v-if="selectedRowKeys.length > 0" class="batch-actions">
        <div class="batch-info">
          已选择 {{ selectedRowKeys.length }} 项
          <a-button type="link" @click="selectedRowKeys = []">
            取消选择
          </a-button>
        </div>
        <a-space>
          <a-popconfirm v-if="resourceConfig?.deletable !== false" title="确定删除选中的数据？" ok-text="确定" cancel-text="取消" @confirm="handleBatchDelete">
            <a-button type="primary" danger>
              批量删除
            </a-button>
          </a-popconfirm>
          <a-dropdown v-if="(resourceConfig?.actions || []).some(a => a.bulk && canShowAction(a))">
            <template #overlay>
              <a-menu>
                <a-menu-item v-for="act in resourceConfig.actions.filter(a => a.bulk && canShowAction(a))" :key="act.name" @click="() => triggerAction(act)">
                  {{ act.label }}
                </a-menu-item>
              </a-menu>
            </template>
            <a-button>批量操作 <DownOutlined /></a-button>
          </a-dropdown>
        </a-space>
      </div>
    </a-card>
  </page-container>
</template>

<style scoped>
.flex {
  display: flex;
}

.gap-2 {
  gap: 0.5rem;
}

.mb-4 {
  margin-bottom: 1rem;
}

.batch-actions {
  margin-top: 16px;
  padding: 12px 0;
  border-top: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.batch-info {
  font-size: 14px;
  color: #666;
}
</style>
