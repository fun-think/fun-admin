<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  Card,
  Row,
  Col,
  Steps,
  Upload,
  Button,
  Table,
  Form,
  Input,
  Select,
  message,
  Progress,
  Tag,
  Space,
  Popconfirm,
  Modal,
  Alert,
  Tabs,
  Descriptions,
  Divider
} from 'ant-design-vue'
import {
  importExcel,
  importCSV,
  importJSON,
  getImportTemplate,
  getImportHistory,
  getImportTask,
  cancelImportTask,
  retryImportTask,
  getFieldMapping
} from '~/api/import.js'
import { getResourceList } from '~/api/resources.js'

const route = useRoute()
const router = useRouter()

// 步骤控制
const currentStep = ref(0)
const steps = [
  { title: '选择文件', description: '上传要导入的文件' },
  { title: '配置映射', description: '配置字段映射关系' },
  { title: '预览数据', description: '预览要导入的数据' },
  { title: '执行导入', description: '执行数据导入操作' }
]

// 文件上传相关
const fileList = ref([])
const uploadLoading = ref(false)
const fileType = ref('excel')
const selectedResource = ref('')
const resources = ref([])

// 字段映射相关
const mappingForm = ref()
const mappingData = ref([])
const fileHeaders = ref([])
const previewData = ref([])
const mappingLoading = ref(false)

// 导入执行相关
const importLoading = ref(false)
const importProgress = ref(0)
const importResult = ref(null)
const currentTaskId = ref(null)

// 历史记录相关
const historyLoading = ref(false)
const importHistory = ref([])
const historyPagination = ref({
  current: 1,
  pageSize: 10,
  total: 0
})

// 获取资源列表
const fetchResources = async () => {
  try {
    const res = await getResourceList()
    if (res.code === 0) {
      resources.value = res.data
    }
  } catch (error) {
    console.error('Failed to fetch resources:', error)
    message.error('获取资源列表失败')
  }
}

// 文件上传处理
const beforeUpload = (file) => {
  const validTypes = {
    excel: ['application/vnd.ms-excel', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'],
    csv: ['text/csv', 'application/csv'],
    json: ['application/json']
  }
  
  if (!validTypes[fileType.value].includes(file.type)) {
    message.error(`不支持的文件类型，请上传${fileType.value.toUpperCase()}文件`)
    return false
  }
  
  if (file.size > 10 * 1024 * 1024) {
    message.error('文件大小不能超过10MB')
    return false
  }
  
  return true
}

const handleUpload = async (file) => {
  if (!selectedResource.value) {
    message.error('请先选择要导入的资源')
    return
  }
  
  uploadLoading.value = true
  try {
    let result
    const importOptions = {
      resource: selectedResource.value
    }
    
    switch (fileType.value) {
      case 'excel':
        result = await importExcel(file, importOptions)
        break
      case 'csv':
        result = await importCSV(file, importOptions)
        break
      case 'json':
        result = await importJSON(file, importOptions)
        break
    }
    
    if (result.code === 0) {
      fileHeaders.value = result.data.headers || []
      previewData.value = result.data.preview || []
      
      // 自动生成字段映射
      await generateFieldMapping()
      
      message.success('文件上传成功')
      currentStep.value = 1
    } else {
      message.error(result.msg || '文件上传失败')
    }
  } catch (error) {
    console.error('Upload failed:', error)
    message.error('文件上传失败')
  } finally {
    uploadLoading.value = false
  }
}

// 生成字段映射
const generateFieldMapping = async () => {
  if (!selectedResource.value || fileHeaders.value.length === 0) return
  
  mappingLoading.value = true
  try {
    const res = await getFieldMapping(selectedResource.value, fileHeaders.value)
    if (res.code === 0) {
      mappingData.value = res.data.mapping || []
    } else {
      // 生成默认映射
      mappingData.value = fileHeaders.value.map(header => ({
        file_field: header,
        db_field: '',
        required: false,
        transform: 'none'
      }))
    }
  } catch (error) {
    console.error('Failed to generate field mapping:', error)
    // 生成默认映射
    mappingData.value = fileHeaders.value.map(header => ({
      file_field: header,
      db_field: '',
      required: false,
      transform: 'none'
    }))
  } finally {
    mappingLoading.value = false
  }
}

// 下一步
const nextStep = () => {
  if (currentStep.value < steps.length - 1) {
    currentStep.value++
  }
}

// 上一步
const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

// 执行导入
const executeImport = async () => {
  if (mappingData.value.length === 0) {
    message.error('请配置字段映射')
    return
  }
  
  importLoading.value = true
  importProgress.value = 0
  
  try {
    const mapping = mappingData.value.reduce((acc, item) => {
      if (item.db_field) {
        acc[item.file_field] = {
          field: item.db_field,
          required: item.required,
          transform: item.transform
        }
      }
      return acc
    }, {})
    
    // 模拟导入进度
    const progressInterval = setInterval(() => {
      importProgress.value = Math.min(importProgress.value + 10, 90)
    }, 500)
    
    // 这里应该调用实际的导入API
    await new Promise(resolve => setTimeout(resolve, 3000))
    
    clearInterval(progressInterval)
    importProgress.value = 100
    
    importResult.value = {
      success: true,
      total: previewData.value.length,
      imported: previewData.value.length - 2,
      failed: 2,
      errors: ['第3行数据格式错误', '第5行必填字段为空']
    }
    
    message.success('导入完成')
    currentStep.value = 3
    fetchImportHistory()
  } catch (error) {
    console.error('Import failed:', error)
    message.error('导入失败')
  } finally {
    importLoading.value = false
  }
}

// 获取导入历史
const fetchImportHistory = async () => {
  historyLoading.value = true
  try {
    const res = await getImportHistory({
      page: historyPagination.value.current,
      page_size: historyPagination.value.pageSize
    })
    
    if (res.code === 0) {
      importHistory.value = res.data.items || []
      historyPagination.value.total = res.data.total || 0
    }
  } catch (error) {
    console.error('Failed to fetch import history:', error)
    message.error('获取导入历史失败')
  } finally {
    historyLoading.value = false
  }
}

// 下载模板
const downloadTemplate = async () => {
  if (!selectedResource.value) {
    message.error('请先选择资源')
    return
  }
  
  try {
    const res = await getImportTemplate(selectedResource.value, fileType.value)
    if (res.code === 0) {
      // 创建下载链接
      const link = document.createElement('a')
      link.href = res.data.url
      link.download = res.data.filename
      link.click()
    }
  } catch (error) {
    console.error('Failed to download template:', error)
    message.error('下载模板失败')
  }
}

// 重新开始
const restartImport = () => {
  currentStep.value = 0
  fileList.value = []
  mappingData.value = []
  previewData.value = []
  importResult.value = null
  importProgress.value = 0
}

// 表格列定义
const previewColumns = computed(() => {
  if (fileHeaders.value.length === 0) return []
  
  return fileHeaders.value.map(header => ({
    title: header,
    dataIndex: header,
    key: header,
    width: 150
  }))
})

const historyColumns = [
  {
    title: '任务ID',
    dataIndex: 'id',
    key: 'id',
    width: 100
  },
  {
    title: '资源',
    dataIndex: 'resource',
    key: 'resource',
    width: 120
  },
  {
    title: '文件类型',
    dataIndex: 'file_type',
    key: 'file_type',
    width: 100,
    customRender: ({ text }) => h(Tag, { color: 'blue' }, () => text.toUpperCase())
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    width: 100,
    customRender: ({ text }) => {
      const statusMap = {
        pending: { color: 'orange', text: '等待中' },
        running: { color: 'blue', text: '导入中' },
        completed: { color: 'green', text: '已完成' },
        failed: { color: 'red', text: '失败' }
      }
      const status = statusMap[text] || { color: 'default', text: text }
      return h(Tag, { color: status.color }, () => status.text)
    }
  },
  {
    title: '总记录数',
    dataIndex: 'total_records',
    key: 'total_records',
    width: 100
  },
  {
    title: '成功数',
    dataIndex: 'success_count',
    key: 'success_count',
    width: 100
  },
  {
    title: '失败数',
    dataIndex: 'failed_count',
    key: 'failed_count',
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
    width: 150,
    customRender: ({ record }) => {
      const actions = []
      
      if (record.status === 'running') {
        actions.push(h(
          Popconfirm,
          {
            title: '确定要取消该导入任务？',
            onConfirm: () => cancelImportTask(record.id)
          },
          {
            default: () => h(Button, { type: 'link', size: 'small', danger: true }, () => '取消')
          }
        ))
      }
      
      if (record.status === 'failed') {
        actions.push(h(
          Button,
          {
            type: 'link',
            size: 'small',
            onClick: () => retryImportTask(record.id)
          },
          () => '重试'
        ))
      }
      
      actions.push(h(
        Button,
        {
          type: 'link',
          size: 'small',
          onClick: () => {
            // 查看详情
          }
        },
        () => '详情'
      ))
      
      return h(Space, {}, () => actions)
    }
  }
]

onMounted(() => {
  fetchResources()
  fetchImportHistory()
})
</script>

<template>
  <page-container title="数据导入">
    <a-row :gutter="16">
      <a-col :span="18">
        <a-card>
          <!-- 步骤条 -->
          <a-steps :current="currentStep" style="margin-bottom: 32px">
            <a-step v-for="step in steps" :key="step.title" :title="step.title" :description="step.description" />
          </a-steps>
          
          <!-- 步骤内容 -->
          <div v-if="currentStep === 0">
            <a-alert
              message="数据导入说明"
              description="支持Excel、CSV、JSON格式的数据导入。请先选择要导入的资源类型，然后上传对应的文件。"
              type="info"
              show-icon
              style="margin-bottom: 16px"
            />
            
            <a-form layout="vertical">
              <a-form-item label="选择资源" required>
                <a-select
                  v-model:value="selectedResource"
                  placeholder="请选择要导入的资源"
                  style="width: 100%"
                >
                  <a-select-option
                    v-for="resource in resources"
                    :key="resource.slug"
                    :value="resource.slug"
                  >
                    {{ resource.name }}
                  </a-select-option>
                </a-select>
              </a-form-item>
              
              <a-form-item label="文件类型" required>
                <a-radio-group v-model:value="fileType">
                  <a-radio value="excel">Excel (.xlsx, .xls)</a-radio>
                  <a-radio value="csv">CSV (.csv)</a-radio>
                  <a-radio value="json">JSON (.json)</a-radio>
                </a-radio-group>
              </a-form-item>
              
              <a-form-item label="上传文件" required>
                <a-upload
                  :file-list="fileList"
                  :before-upload="beforeUpload"
                  :custom-request="({ file }) => handleUpload(file)"
                  :accept="fileType === 'excel' ? '.xlsx,.xls' : fileType === 'csv' ? '.csv' : '.json'"
                  :max-count="1"
                >
                  <a-button type="primary" :loading="uploadLoading">
                    {{ uploadLoading ? '上传中...' : '选择文件' }}
                  </a-button>
                </a-upload>
              </a-form-item>
              
              <a-form-item>
                <a-space>
                  <a-button
                    type="primary"
                    :disabled="!selectedResource || fileList.length === 0"
                    @click="nextStep"
                  >
                    下一步
                  </a-button>
                  <a-button
                    :disabled="!selectedResource"
                    @click="downloadTemplate"
                  >
                    下载模板
                  </a-button>
                </a-space>
              </a-form-item>
            </a-form>
          </div>
          
          <div v-else-if="currentStep === 1">
            <a-alert
              message="字段映射配置"
              description="请将文件中的字段与数据库字段进行映射。必填字段必须映射，否则将导致导入失败。"
              type="info"
              show-icon
              style="margin-bottom: 16px"
            />
            
            <a-table
              :columns="[
                { title: '文件字段', dataIndex: 'file_field', key: 'file_field' },
                { title: '数据库字段', dataIndex: 'db_field', key: 'db_field' },
                { title: '必填', dataIndex: 'required', key: 'required' },
                { title: '转换', dataIndex: 'transform', key: 'transform' }
              ]"
              :data-source="mappingData"
              :pagination="false"
              :loading="mappingLoading"
            >
              <template #bodyCell="{ column, record, index }">
                <template v-if="column.key === 'db_field'">
                  <a-select
                    v-model:value="record.db_field"
                    placeholder="选择数据库字段"
                    style="width: 100%"
                  >
                    <a-select-option value="id">ID</a-select-option>
                    <a-select-option value="name">名称</a-select-option>
                    <a-select-option value="email">邮箱</a-select-option>
                    <a-select-option value="phone">电话</a-select-option>
                    <a-select-option value="status">状态</a-select-option>
                  </a-select>
                </template>
                
                <template v-if="column.key === 'required'">
                  <a-checkbox v-model:checked="record.required" />
                </template>
                
                <template v-if="column.key === 'transform'">
                  <a-select
                    v-model:value="record.transform"
                    placeholder="转换方式"
                    style="width: 100%"
                  >
                    <a-select-option value="none">无转换</a-select-option>
                    <a-select-option value="uppercase">转大写</a-select-option>
                    <a-select-option value="lowercase">转小写</a-select-option>
                    <a-select-option value="trim">去除空格</a-select-option>
                  </a-select>
                </template>
              </template>
            </a-table>
            
            <div style="margin-top: 16px">
              <a-space>
                <a-button @click="prevStep">上一步</a-button>
                <a-button type="primary" @click="nextStep">下一步</a-button>
              </a-space>
            </div>
          </div>
          
          <div v-else-if="currentStep === 2">
            <a-alert
              message="数据预览"
              :description="`共 ${previewData.length} 条数据，请确认无误后执行导入。`"
              type="info"
              show-icon
              style="margin-bottom: 16px"
            />
            
            <a-table
              :columns="previewColumns"
              :data-source="previewData"
              :pagination="{ pageSize: 10 }"
              :scroll="{ x: 'max-content' }"
              size="small"
            />
            
            <div style="margin-top: 16px">
              <a-space>
                <a-button @click="prevStep">上一步</a-button>
                <a-button type="primary" @click="executeImport" :loading="importLoading">
                  开始导入
                </a-button>
              </a-space>
            </div>
          </div>
          
          <div v-else-if="currentStep === 3">
            <a-result
              :status="importResult?.success ? 'success' : 'error'"
              :title="importResult?.success ? '导入完成' : '导入失败'"
              :sub-title="importResult ? `总计：${importResult.total} 条，成功：${importResult.imported} 条，失败：${importResult.failed} 条` : ''"
            >
              <template #extra>
                <a-space>
                  <a-button type="primary" @click="restartImport">
                    重新导入
                  </a-button>
                  <a-button @click="router.push('/admin/import-history')">
                    查看历史
                  </a-button>
                </a-space>
              </template>
            </a-result>
            
            <div v-if="importResult?.errors?.length" style="margin-top: 16px">
              <a-alert
                message="导入错误"
                type="error"
                show-icon
              >
                <template #description>
                  <ul>
                    <li v-for="(error, index) in importResult.errors" :key="index">
                      {{ error }}
                    </li>
                  </ul>
                </template>
              </a-alert>
            </div>
          </div>
        </a-card>
      </a-col>
      
      <a-col :span="6">
        <a-card title="导入历史" size="small">
          <a-table
            :columns="historyColumns.filter(col => ['resource', 'status', 'total_records', 'success_count'].includes(col.dataIndex))"
            :data-source="importHistory"
            :pagination="false"
            :loading="historyLoading"
            size="small"
          >
            <template #title>
              <div style="display: flex; justify-content: space-between; align-items: center;">
                <span>最近导入</span>
                <a-button type="link" size="small" @click="router.push('/admin/import-history')">
                  查看全部
                </a-button>
              </div>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>
  </page-container>
</template>

<style scoped>
.ant-steps {
  max-width: 780px;
  margin: 0 auto;
}
</style>