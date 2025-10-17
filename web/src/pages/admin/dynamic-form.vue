<script setup>
import { ref, onMounted, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  getResourceConfig,
  getResourceRecord,
  createResourceRecord,
  updateResourceRecord
} from '~/api/admin/resources.js'
import { uploadFile, getStorageConfig } from '~/api/common/file.js'

// 引入 Ant Design Vue 组件
import {
  Input,
  InputNumber,
  Select,
  Textarea,
  Checkbox,
  DatePicker,
  Upload,
  Button,
  message,
  Progress
} from 'ant-design-vue'

const route = useRoute()
const router = useRouter()

// 获取资源标识符和操作类型
const resourceSlug = route.params.slug
const operation = route.params.operation // 'create' 或 'edit'
const recordId = route.params.id // 编辑时的记录ID

// 表单相关状态
const formRef = ref()
const formModel = ref({})
const formFields = ref([])
const resourceConfig = ref(null)
const loading = ref(false)
const formRules = ref({})
const relationshipData = ref({}) // 存储关联数据
const fileList = ref({}) // 存储文件上传列表
const uploadProgress = ref({}) // 存储上传进度
const storageConfig = ref(null) // 存储配置

// 获取当前语言
const getCurrentLanguage = () => {
  // 这里可以从 localStorage 或其他地方获取当前语言设置
  return localStorage.getItem('admin-language') || 'zh-CN'
}

// 获取存储配置
const fetchStorageConfig = async () => {
  try {
    const res = await getStorageConfig()
    if (res.code === 0) {
      storageConfig.value = res.data
    }
  } catch (error) {
    console.error('Failed to fetch storage config:', error)
  }
}

// 获取资源配置
const fetchResourceConfig = async () => {
  try {
    const res = await getResourceConfig(resourceSlug, { language: getCurrentLanguage() })
    if (res.code === 0) {
      resourceConfig.value = res.data
      formFields.value = res.data.fields
      
      // 初始化表单模型和验证规则
      const model = {}
      const rules = {}
      
      res.data.fields.forEach(field => {
        // 根据字段类型设置默认值
        switch (field.type) {
          case 'boolean':
            model[field.name] = false
            break
          case 'number':
            model[field.name] = undefined
            break
          case 'file':
            model[field.name] = null
            fileList.value[field.name] = []
            break
          default:
            model[field.name] = undefined
        }
        
        // 构建验证规则
        const fieldRules = []
        if (field.required) {
          fieldRules.push({ required: true, message: `${field.label}为必填项` })
        }
        rules[field.name] = fieldRules
        
        // 如果是关联字段，加载关联数据
        if (field.type === 'relationship') {
          loadRelationshipData(field)
        }
      })
      
      formModel.value = model
      formRules.value = rules
      
      // 如果是编辑操作，获取记录详情
      if (operation === 'edit' && recordId) {
        fetchRecordDetail(recordId)
      }
    }
  } catch (error) {
    console.error('Failed to fetch resource config:', error)
    message.error('获取资源配置失败')
  }
}

// 加载关联数据
const loadRelationshipData = async (field) => {
  try {
    const res = await getResourceData(field.related_resource, {
      page: 1,
      page_size: 100, // 获取所有关联数据
      language: getCurrentLanguage()
    })
    
    if (res.code === 0) {
      relationshipData.value[field.name] = res.data.items || []
    }
  } catch (error) {
    console.error(`Failed to load relationship data for ${field.name}:`, error)
    message.error(`加载${field.label}数据失败`)
  }
}

// 获取记录详情（编辑时）
const fetchRecordDetail = async (id) => {
  loading.value = true
  try {
    const res = await getResourceRecord(resourceSlug, id, { language: getCurrentLanguage() })
    if (res.code === 0) {
      formModel.value = res.data
      
      // 处理文件字段
      formFields.value.forEach(field => {
        if (field.type === 'file' && res.data[field.name]) {
          // 如果有已上传的文件，初始化文件列表
          fileList.value[field.name] = [{
            uid: '-1',
            name: res.data[field.name].split('/').pop(),
            status: 'done',
            url: res.data[field.name]
          }]
        }
      })
    }
  } catch (error) {
    console.error('Failed to fetch record detail:', error)
    message.error('获取记录详情失败')
  } finally {
    loading.value = false
  }
}

// 表单提交
const handleSubmit = async () => {
  try {
    // 验证表单
    await formRef.value.validate()
    
    loading.value = true
    
    // 处理文件上传
    const formData = new FormData()
    const jsonData = {}
    
    // 分离文件字段和普通字段
    for (const key in formModel.value) {
      if (formModel.value.hasOwnProperty(key)) {
        const field = formFields.value.find(f => f.name === key)
        if (field && field.type === 'file') {
          // 文件字段
          if (fileList.value[key] && fileList.value[key].length > 0) {
            const file = fileList.value[key][0].originFileObj
            if (file) {
              formData.append('file', file)
              // 这里应该上传文件并获取URL，简化处理直接使用文件名
              jsonData[key] = `/uploads/${file.name}`
            }
          }
        } else {
          // 普通字段
          jsonData[key] = formModel.value[key]
        }
      }
    }
    
    let res
    if (operation === 'create') {
      // 创建操作
      res = await createResourceRecord(resourceSlug, jsonData, { language: getCurrentLanguage() })
    } else {
      // 编辑操作
      res = await updateResourceRecord(resourceSlug, recordId, jsonData, { language: getCurrentLanguage() })
    }
    
    if (res.code === 0) {
      message.success(operation === 'create' ? '创建成功' : '更新成功')
      // 返回列表页
      router.push(`/admin/${resourceSlug}`)
    } else {
      // 显示后端返回的错误信息
      if (res.errors) {
        // 将后端错误映射到表单字段
        const fields = Object.keys(res.errors)
        fields.forEach(field => {
          formRef.value.setFields([{
            name: field,
            errors: res.errors[field]
          }])
        })
      } else {
        message.error(res.message || (operation === 'create' ? '创建失败' : '更新失败'))
      }
    }
  } catch (error) {
    console.error('Failed to submit form:', error)
    message.error(operation === 'create' ? '创建失败' : '更新失败')
  } finally {
    loading.value = false
  }
}

// 返回列表页
const handleBack = () => {
  router.push(`/admin/${resourceSlug}`)
}

// 文件上传相关函数
const beforeUpload = (file, field) => {
  // 文件类型验证
  if (field.accept) {
    const acceptTypes = field.accept.split(',').map(type => type.trim())
    const fileExtension = file.name.split('.').pop().toLowerCase()
    const isValidType = acceptTypes.some(type => {
      if (type.startsWith('.')) {
        return fileExtension === type.substring(1)
      } else if (type.includes('/*')) {
        return file.type.startsWith(type.split('/*')[0])
      }
      return file.type === type
    })
    
    if (!isValidType) {
      message.error(`不支持的文件类型：${file.name}`)
      return false
    }
  }
  
  // 文件大小验证
  if (field.maxSize && file.size > field.maxSize * 1024 * 1024) {
    message.error(`文件大小不能超过 ${field.maxSize}MB`)
    return false
  }
  
  return true
}

const handleUpload = async (file, fieldName) => {
  try {
    uploadProgress.value[fieldName] = 0
    
    const field = formFields.value.find(f => f.name === fieldName)
    const options = {
      onUploadProgress: (progressEvent) => {
        const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
        uploadProgress.value[fieldName] = progress
      }
    }
    
    // 使用云存储
    if (storageConfig.value && storageConfig.value.default_type) {
      options.storageType = storageConfig.value.default_type
    }
    
    const res = await uploadFile(file, options)
    
    if (res.code === 0) {
      formModel.value[fieldName] = res.data.url
      // 添加到文件列表
      fileList.value[fieldName] = [{
        uid: file.uid || Date.now(),
        name: file.name,
        status: 'done',
        url: res.data.url,
        size: file.size,
        type: file.type
      }]
      message.success('文件上传成功')
    } else {
      message.error(res.msg || '文件上传失败')
    }
  } catch (error) {
    console.error('Upload failed:', error)
    message.error('文件上传失败')
  } finally {
    uploadProgress.value[fieldName] = 0
  }
}

const handleChange = (info, fieldName) => {
  const { file } = info
  
  if (file.status === 'uploading') {
    fileList.value[fieldName] = [file]
  } else if (file.status === 'done') {
    fileList.value[fieldName] = [file]
  } else if (file.status === 'error') {
    message.error(`${file.name} 文件上传失败`)
  }
}

const handleRemove = (fieldName) => {
  fileList.value[fieldName] = []
  formModel.value[fieldName] = null
  uploadProgress.value[fieldName] = 0
}

// 根据字段类型渲染表单控件
const renderField = (field) => {
  switch (field.type) {
    case 'text':
      return h(
        Input,
        {
          value: formModel.value[field.name],
          'onUpdate:value': (value) => {
            formModel.value[field.name] = value
          },
          placeholder: `请输入${field.label}`
        }
      )
    
    case 'email':
      return h(
        Input,
        {
          type: 'email',
          value: formModel.value[field.name],
          'onUpdate:value': (value) => {
            formModel.value[field.name] = value
          },
          placeholder: `请输入${field.label}`
        }
      )
    
    case 'number':
      return h(
        InputNumber,
        {
          value: formModel.value[field.name],
          'onUpdate:value': (value) => {
            formModel.value[field.name] = value
          },
          style: { width: '100%' },
          placeholder: `请输入${field.label}`
        }
      )
    
    case 'select':
      return h(
        Select,
        {
          value: formModel.value[field.name],
          'onUpdate:value': (value) => {
            formModel.value[field.name] = value
          },
          placeholder: `请选择${field.label}`,
          options: field.options?.map(option => ({
            label: option.label,
            value: option.value
          })) || []
        }
      )
    
    case 'textarea':
      return h(
        Textarea,
        {
          value: formModel.value[field.name],
          'onUpdate:value': (value) => {
            formModel.value[field.name] = value
          },
          placeholder: `请输入${field.label}`,
          rows: field.rows || 4
        }
      )
    
    case 'boolean':
      return h(
        Checkbox,
        {
          checked: formModel.value[field.name],
          'onUpdate:checked': (checked) => {
            formModel.value[field.name] = checked
          }
        },
        { default: () => field.label }
      )
    
    case 'date':
      return h(
        DatePicker,
        {
          value: formModel.value[field.name],
          'onUpdate:value': (value) => {
            formModel.value[field.name] = value
          },
          placeholder: `请选择${field.label}`,
          style: { width: '100%' }
        }
      )
    
    case 'datetime':
      return h(
        DatePicker,
        {
          value: formModel.value[field.name],
          'onUpdate:value': (value) => {
            formModel.value[field.name] = value
          },
          placeholder: `请选择${field.label}`,
          showTime: true,
          style: { width: '100%' }
        }
      )
    
    case 'relationship':
      // 获取关联数据
      const relatedData = relationshipData.value[field.name] || []
      return h(
        Select,
        {
          value: formModel.value[field.name],
          'onUpdate:value': (value) => {
            formModel.value[field.name] = value
          },
          placeholder: `请选择${field.label}`,
          options: relatedData.map(item => ({
            label: item[field.display_field || 'name'] || item.id,
            value: item.id
          }))
        }
      )
    
    case 'file':
      return h(
        'div',
        {},
        [
          h(
            Upload,
            {
              name: 'file',
              fileList: fileList.value[field.name] || [],
              beforeUpload: (file) => beforeUpload(file, field),
              customRequest: ({ file }) => handleUpload(file, field.name),
              onChange: (info) => handleChange(info, field.name),
              onRemove: () => handleRemove(field.name),
              accept: field.accept || '*',
              maxCount: 1,
              showUploadList: {
                showPreviewIcon: true,
                showRemoveIcon: true,
                showDownloadIcon: true
              }
            },
            {
              default: () => h(Button, { type: 'primary' }, () => '上传文件')
            }
          ),
          uploadProgress.value[field.name] > 0 && uploadProgress.value[field.name] < 100 ? h(
            Progress,
            {
              percent: uploadProgress.value[field.name],
              size: 'small',
              style: { marginTop: '8px' }
            }
          ) : null,
          h(
            'div',
            { style: { marginTop: '8px', fontSize: '12px', color: '#666' } },
            [
              field.accept ? `支持格式：${field.accept} ` : '',
              field.maxSize ? `大小限制：${field.maxSize}MB` : ''
            ].filter(Boolean).join(' | ')
          ),
          formModel.value[field.name] ? h(
            'div',
            { style: { marginTop: '8px' } },
            h('a', { href: formModel.value[field.name], target: '_blank' }, '查看文件')
          ) : null
        ]
      )
    
    default:
      return h(
        Input,
        {
          value: formModel.value[field.name],
          'onUpdate:value': (value) => {
            formModel.value[field.name] = value
          },
          placeholder: `请输入${field.label}`
        }
      )
  }
}

onMounted(() => {
  fetchResourceConfig()
  fetchStorageConfig()
})
</script>

<template>
  <page-container :title="operation === 'create' ? '新增' : '编辑'">
    <a-card>
      <a-form
        ref="formRef"
        :model="formModel"
        :rules="formRules"
        :label-col="{ span: 4 }"
        :wrapper-col="{ span: 14 }"
      >
        <a-form-item
          v-for="field in formFields"
          :key="field.name"
          :label="field.type === 'boolean' ? '' : field.label"
          :name="field.name"
        >
          <component
            :is="renderField(field)"
            v-model:value="formModel[field.name]"
          />
        </a-form-item>
        
        <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
          <a-button
            type="primary"
            :loading="loading"
            @click="handleSubmit"
          >
            {{ operation === 'create' ? '创建' : '更新' }}
          </a-button>
          <a-button
            style="margin-left: 10px"
            @click="handleBack"
          >
            返回
          </a-button>
          
          <!-- 语言切换 -->
          <a-dropdown style="margin-left: 10px">
            <template #overlay>
              <a-menu @click="({ key }) => { localStorage.setItem('admin-language', key); location.reload(); }">
                <a-menu-item key="zh-CN">中文</a-menu-item>
                <a-menu-item key="en">English</a-menu-item>
              </a-menu>
            </template>
            <a-button>
              语言 <DownOutlined />
            </a-button>
          </a-dropdown>
        </a-form-item>
      </a-form>
    </a-card>
  </page-container>
</template>