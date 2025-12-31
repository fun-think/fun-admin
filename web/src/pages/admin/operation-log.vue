<script setup>
import { onMounted, reactive, ref, shallowRef } from 'vue'
import { message } from 'ant-design-vue'
import { deleteOperationLog, deleteOperationLogs, getOperationLogs } from '@/api/operation-log.js'

// 操作日志数据
const operationLogs = ref([])
const loading = ref(false)
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  onChange(current, pageSize) {
    pagination.pageSize = pageSize
    pagination.current = current
    fetchOperationLogs()
  },
})

// 搜索表单
const searchForm = reactive({
  keyword: '',
  resource: '',
  action: '',
})

// 获取操作日志列表
async function fetchOperationLogs(page = 1) {
  if (loading.value)
    return
  loading.value = true
  try {
    const params = {
      page,
      page_size: pagination.pageSize,
      search_keyword: searchForm.keyword,
      resource: searchForm.resource,
      action: searchForm.action,
    }

    const res = await getOperationLogs(params)

    if (res.data.code === 0) {
      operationLogs.value = res.data.data.items || []
      pagination.total = res.data.data.total || 0
      pagination.current = page
    }
    else {
      message.error(res.data.msg || '获取操作日志失败')
    }
  }
  catch (error) {
    console.error('Failed to fetch operation logs:', error)
    message.error('获取操作日志失败')
  }
  finally {
    loading.value = false
  }
}

// 搜索
async function handleSearch() {
  pagination.current = 1
  await fetchOperationLogs()
}

// 重置搜索
async function handleReset() {
  Object.assign(searchForm, {
    keyword: '',
    resource: '',
    action: '',
  })
  await fetchOperationLogs()
}

// 删除操作日志
async function handleDelete(id) {
  const close = message.loading('删除中......')
  try {
    const res = await deleteOperationLog(id)

    if (res.data.code === 0) {
      message.success('删除成功')
      await fetchOperationLogs(pagination.current)
    }
    else {
      message.error(res.data.msg || '删除失败')
    }
  }
  catch (error) {
    console.error('Failed to delete operation logger:', error)
    message.error('删除失败')
  }
  finally {
    close()
  }
}

// 批量删除
const selectedRowKeys = ref([])
async function handleBatchDelete() {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请先选择要删除的日志')
    return
  }

  const close = message.loading('删除中......')
  try {
    const res = await deleteOperationLogs({ ids: selectedRowKeys.value })

    if (res.data.code === 0) {
      message.success('批量删除成功')
      selectedRowKeys.value = []
      await fetchOperationLogs(pagination.current)
    }
    else {
      message.error(res.data.msg || '批量删除失败')
    }
  }
  catch (error) {
    console.error('Failed to batch delete operation logs:', error)
    message.error('批量删除失败')
  }
  finally {
    close()
  }
}

// 表格列定义
const columns = shallowRef([
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '操作用户',
    dataIndex: 'user_name',
    key: 'user_name',
  },
  {
    title: 'IP地址',
    dataIndex: 'ip',
    key: 'ip',
  },
  {
    title: '请求方法',
    dataIndex: 'method',
    key: 'method',
  },
  {
    title: '请求路径',
    dataIndex: 'path',
    key: 'path',
  },
  {
    title: '状态码',
    dataIndex: 'status_code',
    key: 'status_code',
  },
  {
    title: '执行时长(ms)',
    dataIndex: 'duration',
    key: 'duration',
  },
  {
    title: '操作描述',
    dataIndex: 'description',
    key: 'description',
  },
  {
    title: '资源',
    dataIndex: 'resource',
    key: 'resource',
  },
  {
    title: '操作类型',
    dataIndex: 'action',
    key: 'action',
  },
  {
    title: '操作时间',
    dataIndex: 'created_at',
    key: 'created_at',
  },
  {
    title: '操作',
    key: 'actions',
  },
])

// 表格行选择配置
const rowSelection = {
  onChange: (selectedRowKeysValue) => {
    selectedRowKeys.value = selectedRowKeysValue
  },
}

onMounted(() => {
  fetchOperationLogs()
})
</script>

<template>
  <page-container title="操作日志">
    <a-card>
      <!-- 搜索区域 -->
      <a-form layout="inline" class="mb-4">
        <a-form-item label="关键词">
          <a-input
            v-model:value="searchForm.keyword"
            placeholder="用户名/路径/描述"
            allow-clear
          />
        </a-form-item>

        <a-form-item label="资源">
          <a-input
            v-model:value="searchForm.resource"
            placeholder="资源名称"
            allow-clear
          />
        </a-form-item>

        <a-form-item label="操作类型">
          <a-select
            v-model:value="searchForm.action"
            placeholder="请选择操作类型"
            style="width: 120px"
            allow-clear
          >
            <a-select-option value="create">
              创建
            </a-select-option>
            <a-select-option value="update">
              更新
            </a-select-option>
            <a-select-option value="delete">
              删除
            </a-select-option>
            <a-select-option value="read">
              读取
            </a-select-option>
            <a-select-option value="other">
              其他
            </a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item>
          <a-button type="primary" @click="handleSearch">
            搜索
          </a-button>
          <a-button style="margin-left: 8px" @click="handleReset">
            重置
          </a-button>

          <a-popconfirm
            v-if="selectedRowKeys.length > 0"
            title="确定删除选中的日志？"
            ok-text="确定"
            cancel-text="取消"
            @confirm="handleBatchDelete"
          >
            <a-button type="primary" danger style="margin-left: 8px">
              批量删除 ({{ selectedRowKeys.length }})
            </a-button>
          </a-popconfirm>
        </a-form-item>
      </a-form>

      <!-- 表格 -->
      <a-table
        row-key="id"
        :loading="loading"
        :columns="columns"
        :data-source="operationLogs"
        :pagination="pagination"
        :row-selection="rowSelection"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'created_at'">
            {{ record.created_at ? new Date(record.created_at).toLocaleString() : '' }}
          </template>
          <template v-else-if="column.key === 'actions'">
            <a-space>
              <a-button type="link" size="small" @click="() => { message.info('查看详情功能待实现') }">
                查看
              </a-button>
              <a-popconfirm
                title="确定删除该条日志？"
                @confirm="handleDelete(record.id)"
              >
                <a-button type="link" size="small" danger>
                  删除
                </a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>
  </page-container>
</template>

<style scoped>
.mb-4 {
  margin-bottom: 1rem;
}
</style>
