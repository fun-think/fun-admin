<script setup>
import { ref, onMounted, computed } from 'vue'
import {
  Card,
  Button,
  Select,
  message,
  Space,
  Tabs,
  Alert,
  Table,
  Tag,
  Badge,
  Descriptions,
  Modal,
  Form,
  Input,
  Switch,
  Divider,
  Row,
  Col,
  CopyOutlined,
  DownloadOutlined,
  ReloadOutlined
} from 'ant-design-vue'
import { useGet, usePost } from '~/utils/request'

// 文档生成相关
const generateLoading = ref(false)
const selectedFormat = ref('swagger')
const includeResources = ref(true)
const includeHandlers = ref(true)
const outputFileName = ref('api-docs')

// 文档预览相关
const docsData = ref(null)
const previewLoading = ref(false)
const activePreviewTab = ref('overview')

// 资源和处理器数据
const resources = ref([])
const handlers = ref([])
const dataLoading = ref(false)

// 格式选项
const formatOptions = [
  { label: 'Swagger (OpenAPI 2.0)', value: 'swagger' },
  { label: 'OpenAPI 3.0', value: 'openapi' },
  { label: 'Markdown', value: 'markdown' },
  { label: 'HTML', value: 'html' }
]

// 生成文档
const generateDocs = async () => {
  generateLoading.value = true
  try {
    const res = await usePost('/api/admin/docs/generate', {
      format: selectedFormat.value,
      include_resources: includeResources.value,
      include_handlers: includeHandlers.value,
      filename: outputFileName.value
    })
    
    if (res.code === 0) {
      docsData.value = res.data
      message.success('文档生成成功')
      activePreviewTab.value = 'overview'
    } else {
      message.error(res.msg || '文档生成失败')
    }
  } catch (error) {
    console.error('Failed to generate docs:', error)
    message.error('文档生成失败')
  } finally {
    generateLoading.value = false
  }
}

// 下载文档
const downloadDocs = async () => {
  if (!docsData.value) {
    message.warning('请先生成文档')
    return
  }
  
  try {
    const res = await useGet('/api/admin/docs/download', {
      format: selectedFormat.value,
      filename: outputFileName.value
    })
    
    if (res.code === 0) {
      // 创建下载链接
      const link = document.createElement('a')
      link.href = res.data.url
      link.download = res.data.filename
      link.click()
      message.success('文档下载成功')
    } else {
      message.error(res.msg || '文档下载失败')
    }
  } catch (error) {
    console.error('Failed to download docs:', error)
    message.error('文档下载失败')
  }
}

// 获取资源数据
const fetchResources = async () => {
  dataLoading.value = true
  try {
    const res = await useGet('/api/admin/resources')
    if (res.code === 0) {
      resources.value = res.data || []
    }
  } catch (error) {
    console.error('Failed to fetch resources:', error)
  } finally {
    dataLoading.value = false
  }
}

// 复制到剪贴板
const copyToClipboard = (text) => {
  navigator.clipboard.writeText(text).then(() => {
    message.success('已复制到剪贴板')
  }).catch(() => {
    message.error('复制失败')
  })
}

// API端点表格列
const endpointColumns = [
  {
    title: '方法',
    dataIndex: 'method',
    key: 'method',
    width: 100,
    customRender: ({ text }) => {
      const methodColors = {
        GET: 'green',
        POST: 'blue',
        PUT: 'orange',
        DELETE: 'red',
        PATCH: 'purple'
      }
      return h(Tag, { color: methodColors[text] || 'default' }, () => text)
    }
  },
  {
    title: '路径',
    dataIndex: 'path',
    key: 'path',
    ellipsis: true
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
    ellipsis: true
  },
  {
    title: '认证',
    dataIndex: 'auth',
    key: 'auth',
    width: 100,
    customRender: ({ text }) => h(
      Badge,
      {
        status: text ? 'success' : 'default',
        text: text ? '需要' : '不需要'
      }
    )
  },
  {
    title: '操作',
    key: 'actions',
    fixed: 'right',
    width: 100,
    customRender: ({ record }) => h(
      Button,
      {
        type: 'link',
        size: 'small',
        icon: h(CopyOutlined),
        onClick: () => copyToClipboard(`${record.method} ${record.path}`)
      },
      () => '复制'
    )
  }
]

// 计算统计数据
const stats = computed(() => {
  if (!docsData.value) return {}
  
  return {
    totalEndpoints: docsData.value.endpoints?.length || 0,
    totalResources: docsData.value.resources?.length || 0,
    authenticatedEndpoints: docsData.value.endpoints?.filter(e => e.auth).length || 0,
    publicEndpoints: docsData.value.endpoints?.filter(e => !e.auth).length || 0
  }
})

// 分组统计
const methodStats = computed(() => {
  if (!docsData.value?.endpoints) return {}
  
  return docsData.value.endpoints.reduce((acc, endpoint) => {
    acc[endpoint.method] = (acc[endpoint.method] || 0) + 1
    return acc
  }, {})
})

onMounted(() => {
  fetchResources()
})
</script>

<template>
  <page-container title="API文档生成">
    <a-row :gutter="16">
      <a-col :span="8">
        <a-card title="生成配置" size="small">
          <a-form layout="vertical">
            <a-form-item label="文档格式">
              <a-select v-model:value="selectedFormat">
                <a-select-option
                  v-for="option in formatOptions"
                  :key="option.value"
                  :value="option.value"
                >
                  {{ option.label }}
                </a-select-option>
              </a-select>
            </a-form-item>
            
            <a-form-item label="输出文件名">
              <a-input v-model:value="outputFileName" placeholder="请输入文件名" />
            </a-form-item>
            
            <a-form-item>
              <a-checkbox v-model:checked="includeResources">
                包含资源文档
              </a-checkbox>
            </a-form-item>
            
            <a-form-item>
              <a-checkbox v-model:checked="includeHandlers">
                包含处理器文档
              </a-checkbox>
            </a-form-item>
            
            <a-form-item>
              <a-space direction="vertical" style="width: 100%">
                <a-button 
                  type="primary" 
                  :loading="generateLoading" 
                  @click="generateDocs"
                  block
                >
                  生成文档
                </a-button>
                
                <a-button 
                  :disabled="!docsData"
                  @click="downloadDocs"
                  block
                >
                  下载文档
                </a-button>
                
                <a-button 
                  :disabled="!docsData"
                  @click="generateDocs"
                  block
                >
                  <template #icon><ReloadOutlined /></template>
                  重新生成
                </a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-card>
        
        <a-card title="统计信息" size="small" style="margin-top: 16px">
          <a-descriptions v-if="docsData" :column="1" size="small">
            <a-descriptions-item label="总端点数">
              <Badge :count="stats.totalEndpoints" />
            </a-descriptions-item>
            <a-descriptions-item label="资源数量">
              <Badge :count="stats.totalResources" />
            </a-descriptions-item>
            <a-descriptions-item label="需要认证">
              <Badge :count="stats.authenticatedEndpoints" color="green" />
            </a-descriptions-item>
            <a-descriptions-item label="公开端点">
              <Badge :count="stats.publicEndpoints" color="blue" />
            </a-descriptions-item>
          </a-descriptions>
          
          <a-empty v-else description="暂无数据" />
        </a-card>
      </a-col>
      
      <a-col :span="16">
        <a-card>
          <a-tabs v-if="docsData" v-model:activeKey="activePreviewTab">
            <a-tab-pane key="overview" tab="概览">
              <a-alert
                message="API文档概览"
                :description="`共 ${stats.totalEndpoints} 个API端点，${stats.totalResources} 个资源`"
                type="info"
                show-icon
                style="margin-bottom: 16px"
              />
              
              <a-descriptions title="基本信息" :column="2" bordered>
                <a-descriptions-item label="文档格式">
                  <Tag color="blue">{{ selectedFormat.toUpperCase() }}</Tag>
                </a-descriptions-item>
                <a-descriptions-item label="生成时间">
                  {{ new Date().toLocaleString() }}
                </a-descriptions-item>
                <a-descriptions-item label="版本">
                  v1.0.0
                </a-descriptions-item>
                <a-descriptions-item label="基础路径">
                  /api
                </a-descriptions-item>
              </a-descriptions>
              
              <Divider />
              
              <h4>HTTP方法统计</h4>
              <a-row :gutter="16">
                <a-col v-for="(count, method) in methodStats" :key="method" :span="6">
                  <a-card size="small">
                    <div style="text-align: center">
                      <Tag :color="{
                        GET: 'green',
                        POST: 'blue',
                        PUT: 'orange',
                        DELETE: 'red',
                        PATCH: 'purple'
                      }[method]">{{ method }}</Tag>
                      <div style="font-size: 24px; font-weight: bold; margin: 8px 0">
                        {{ count }}
                      </div>
                    </div>
                  </a-card>
                </a-col>
              </a-row>
            </a-tab-pane>
            
            <a-tab-pane key="endpoints" tab="API端点">
              <a-table
                :columns="endpointColumns"
                :data-source="docsData.endpoints"
                :pagination="{ pageSize: 10 }"
                size="small"
              >
                <template #title>
                  <span>API端点列表 ({{ docsData.endpoints?.length || 0 }})</span>
                </template>
              </a-table>
            </a-tab-pane>
            
            <a-tab-pane key="resources" tab="资源文档">
              <div v-if="docsData.resources && docsData.resources.length > 0">
                <a-card v-for="resource in docsData.resources" :key="resource.name" size="small" style="margin-bottom: 16px">
                  <template #title>
                    <Space>
                      <span>{{ resource.name }}</span>
                      <Tag color="blue">{{ resource.slug }}</Tag>
                    </Space>
                  </template>
                  
                  <a-descriptions :column="2" size="small">
                    <a-descriptions-item label="描述">
                      {{ resource.description || '无描述' }}
                    </a-descriptions-item>
                    <a-descriptions-item label="端点数量">
                      {{ resource.endpoints?.length || 0 }}
                    </a-descriptions-item>
                  </a-descriptions>
                  
                  <Divider />
                  
                  <h5>可用操作</h5>
                  <a-space>
                    <Tag v-for="action in resource.actions" :key="action" color="green">
                      {{ action }}
                    </Tag>
                  </a-space>
                </a-card>
              </div>
              <a-empty v-else description="暂无资源数据" />
            </a-tab-pane>
            
            <a-tab-pane key="schema" tab="数据模型">
              <a-alert
                message="数据模型"
                description="系统中的数据结构和字段定义"
                type="info"
                show-icon
                style="margin-bottom: 16px"
              />
              
              <div v-if="docsData.schemas && docsData.schemas.length > 0">
                <a-card v-for="schema in docsData.schemas" :key="schema.name" size="small" style="margin-bottom: 16px">
                  <template #title>
                    <Space>
                      <span>{{ schema.name }}</span>
                      <Tag color="purple">{{ schema.type }}</Tag>
                    </Space>
                  </template>
                  
                  <a-table
                    :columns="[
                      { title: '字段名', dataIndex: 'name', key: 'name' },
                      { title: '类型', dataIndex: 'type', key: 'type' },
                      { title: '必填', dataIndex: 'required', key: 'required' },
                      { title: '描述', dataIndex: 'description', key: 'description' }
                    ]"
                    :data-source="schema.fields"
                    :pagination="false"
                    size="small"
                  />
                </a-card>
              </div>
              <a-empty v-else description="暂无数据模型" />
            </a-tab-pane>
          </a-tabs>
          
          <a-empty v-else description="请先生成API文档" />
        </a-card>
      </a-col>
    </a-row>
  </page-container>
</template>

<style scoped>
.ant-card-small {
  margin-bottom: 16px;
}

.ant-statistic {
  text-align: center;
}
</style>