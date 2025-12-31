<script setup>
import { message } from 'ant-design-vue'
import {
  getEmailConfig,
  getStorageConfig,
  getSystemConfig,
  sendTestEmail,
  testEmailConfig,
  testStorageConfig,
  updateEmailConfig,
  updateStorageConfig,
  updateSystemConfig,
} from '@/api/config.js'

// Tab页相关
const activeTab = ref('system')

// 系统配置
const systemConfig = ref({})
const systemLoading = ref(false)
const systemFormRef = ref()

// 邮件配置
const emailConfig = ref({})
const emailLoading = ref(false)
const emailFormRef = ref()
const testEmailLoading = ref(false)
const testEmailModalVisible = ref(false)
const testEmailForm = ref({
  to: '',
  subject: '测试邮件',
  content: '这是一封测试邮件，用于验证邮件配置是否正确。',
})

// 存储配置
const storageConfig = ref({})
const storageLoading = ref(false)
const storageFormRef = ref()
const testStorageLoading = ref(false)

// 获取系统配置
async function fetchSystemConfig() {
  if (systemLoading.value)
    return
  systemLoading.value = true
  try {
    const res = await getSystemConfig()
    if (res.code === 0) {
      systemConfig.value = res.data || {}
    }
  }
  catch (error) {
    console.error('Failed to fetch system config:', error)
    message.error('获取系统配置失败')
  }
  finally {
    systemLoading.value = false
  }
}

// 获取邮件配置
async function fetchEmailConfig() {
  if (emailLoading.value)
    return
  emailLoading.value = true
  try {
    const res = await getEmailConfig()
    if (res.code === 0) {
      emailConfig.value = res.data || {}
    }
  }
  catch (error) {
    console.error('Failed to fetch email config:', error)
    message.error('获取邮件配置失败')
  }
  finally {
    emailLoading.value = false
  }
}

// 获取存储配置
async function fetchStorageConfig() {
  if (storageLoading.value)
    return
  storageLoading.value = true
  try {
    const res = await getStorageConfig()
    if (res.code === 0) {
      storageConfig.value = res.data || {}
    }
  }
  catch (error) {
    console.error('Failed to fetch storage config:', error)
    message.error('获取存储配置失败')
  }
  finally {
    storageLoading.value = false
  }
}

// 保存系统配置
async function saveSystemConfig() {
  try {
    await systemFormRef.value.validateFields()

    systemLoading.value = true
    const res = await updateSystemConfig(systemConfig.value)

    if (res.code === 0) {
      message.success('系统配置保存成功')
    }
    else {
      message.error(res.msg || '保存失败')
    }
  }
  catch (error) {
    console.error('Failed to save system config:', error)
    message.error('保存失败')
  }
  finally {
    systemLoading.value = false
  }
}

// 保存邮件配置
async function saveEmailConfig() {
  try {
    await emailFormRef.value.validateFields()

    emailLoading.value = true
    const res = await updateEmailConfig(emailConfig.value)

    if (res.code === 0) {
      message.success('邮件配置保存成功')
    }
    else {
      message.error(res.msg || '保存失败')
    }
  }
  catch (error) {
    console.error('Failed to save email config:', error)
    message.error('保存失败')
  }
  finally {
    emailLoading.value = false
  }
}

// 测试邮件配置
async function testEmailConfigConnection() {
  try {
    await emailFormRef.value.validateFields()

    testEmailLoading.value = true
    const res = await testEmailConfig(emailConfig.value)

    if (res.code === 0) {
      message.success('邮件配置测试成功')
    }
    else {
      message.error(res.msg || '测试失败')
    }
  }
  catch (error) {
    console.error('Failed to test email config:', error)
    message.error('测试失败')
  }
  finally {
    testEmailLoading.value = false
  }
}

// 发送测试邮件
async function sendTestEmailFunc() {
  try {
    testEmailLoading.value = true
    const res = await sendTestEmail(
      testEmailForm.value.to,
      testEmailForm.value.subject,
      testEmailForm.value.content,
    )

    if (res.code === 0) {
      message.success('测试邮件发送成功')
      testEmailModalVisible.value = false
    }
    else {
      message.error(res.msg || '发送失败')
    }
  }
  catch (error) {
    console.error('Failed to send test email:', error)
    message.error('发送失败')
  }
  finally {
    testEmailLoading.value = false
  }
}

// 保存存储配置
async function saveStorageConfig() {
  try {
    await storageFormRef.value.validateFields()

    storageLoading.value = true
    const res = await updateStorageConfig(storageConfig.value)

    if (res.code === 0) {
      message.success('存储配置保存成功')
    }
    else {
      message.error(res.msg || '保存失败')
    }
  }
  catch (error) {
    console.error('Failed to save storage config:', error)
    message.error('保存失败')
  }
  finally {
    storageLoading.value = false
  }
}

// 测试存储配置
async function testStorageConfigConnection() {
  try {
    await storageFormRef.value.validateFields()

    testStorageLoading.value = true
    const res = await testStorageConfig(storageConfig.value)

    if (res.code === 0) {
      message.success('存储配置测试成功')
    }
    else {
      message.error(res.msg || '测试失败')
    }
  }
  catch (error) {
    console.error('Failed to test storage config:', error)
    message.error('测试失败')
  }
  finally {
    testStorageLoading.value = false
  }
}

// Tab切换处理
function handleTabChange(key) {
  activeTab.value = key
  switch (key) {
    case 'system':
      fetchSystemConfig()
      break
    case 'email':
      fetchEmailConfig()
      break
    case 'storage':
      fetchStorageConfig()
      break
  }
}

// 表单验证规则
const systemRules = {
  app_name: [
    { required: true, message: '请输入应用名称' },
  ],
  app_url: [
    { required: true, message: '请输入应用URL' },
  ],
}

const emailRules = {
  host: [
    { required: true, message: '请输入SMTP服务器' },
  ],
  port: [
    { required: true, message: '请输入端口号' },
  ],
  username: [
    { required: true, message: '请输入用户名' },
  ],
  password: [
    { required: true, message: '请输入密码' },
  ],
  from_email: [
    { required: true, message: '请输入发件人邮箱' },
  ],
}

const storageRules = {
  default_type: [
    { required: true, message: '请选择默认存储类型' },
  ],
}

const testEmailRules = {
  to: [
    { required: true, message: '请输入收件人邮箱' },
    { type: 'email', message: '请输入有效的邮箱地址' },
  ],
}

onMounted(() => {
  fetchSystemConfig()
})
</script>

<template>
  <page-container title="系统设置">
    <a-card>
      <a-tabs v-model:active-key="activeTab" @change="handleTabChange">
        <!-- 系统配置 -->
        <a-tab-pane key="system" tab="系统配置">
          <a-alert
            message="系统基础配置"
            description="配置应用的基本信息，包括应用名称、URL、时区等。"
            type="info"
            show-icon
            style="margin-bottom: 16px"
          />

          <a-form
            ref="systemFormRef"
            :model="systemConfig"
            :rules="systemRules"
            layout="vertical"
          >
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="应用名称" name="app_name">
                  <a-input v-model:value="systemConfig.app_name" placeholder="请输入应用名称" />
                </a-form-item>
              </a-col>

              <a-col :span="12">
                <a-form-item label="应用URL" name="app_url">
                  <a-input v-model:value="systemConfig.app_url" placeholder="请输入应用URL" />
                </a-form-item>
              </a-col>
            </a-row>

            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="时区" name="timezone">
                  <a-select v-model:value="systemConfig.timezone" placeholder="请选择时区">
                    <a-select-option value="Asia/Shanghai">
                      Asia/Shanghai
                    </a-select-option>
                    <a-select-option value="UTC">
                      UTC
                    </a-select-option>
                  </a-select>
                </a-form-item>
              </a-col>

              <a-col :span="12">
                <a-form-item label="语言" name="language">
                  <a-select v-model:value="systemConfig.language" placeholder="请选择语言">
                    <a-select-option value="zh-CN">
                      中文
                    </a-select-option>
                    <a-select-option value="en">
                      English
                    </a-select-option>
                  </a-select>
                </a-form-item>
              </a-col>
            </a-row>

            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="是否开启注册" name="enable_register">
                  <a-switch v-model:checked="systemConfig.enable_register" />
                </a-form-item>
              </a-col>

              <a-col :span="12">
                <a-form-item label="是否开启验证码" name="enable_captcha">
                  <a-switch v-model:checked="systemConfig.enable_captcha" />
                </a-form-item>
              </a-col>
            </a-row>

            <a-form-item>
              <a-button type="primary" :loading="systemLoading" @click="saveSystemConfig">
                保存配置
              </a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- 邮件配置 -->
        <a-tab-pane key="email" tab="邮件配置">
          <a-alert
            message="邮件服务配置"
            description="配置SMTP服务器信息，用于发送系统邮件通知。"
            type="info"
            show-icon
            style="margin-bottom: 16px"
          />

          <a-form
            ref="emailFormRef"
            :model="emailConfig"
            :rules="emailRules"
            layout="vertical"
          >
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="SMTP服务器" name="host">
                  <a-input v-model:value="emailConfig.host" placeholder="请输入SMTP服务器地址" />
                </a-form-item>
              </a-col>

              <a-col :span="12">
                <a-form-item label="端口" name="port">
                  <a-input-number v-model:value="emailConfig.port" placeholder="请输入端口号" style="width: 100%" />
                </a-form-item>
              </a-col>
            </a-row>

            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="用户名" name="username">
                  <a-input v-model:value="emailConfig.username" placeholder="请输入用户名" />
                </a-form-item>
              </a-col>

              <a-col :span="12">
                <a-form-item label="密码" name="password">
                  <a-input-password v-model:value="emailConfig.password" placeholder="请输入密码" />
                </a-form-item>
              </a-col>
            </a-row>

            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="发件人邮箱" name="from_email">
                  <a-input v-model:value="emailConfig.from_email" placeholder="请输入发件人邮箱" />
                </a-form-item>
              </a-col>

              <a-col :span="12">
                <a-form-item label="发件人名称" name="from_name">
                  <a-input v-model:value="emailConfig.from_name" placeholder="请输入发件人名称" />
                </a-form-item>
              </a-col>
            </a-row>

            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="是否启用SSL" name="enable_ssl">
                  <a-switch v-model:checked="emailConfig.enable_ssl" />
                </a-form-item>
              </a-col>

              <a-col :span="12">
                <a-form-item label="是否启用TLS" name="enable_tls">
                  <a-switch v-model:checked="emailConfig.enable_tls" />
                </a-form-item>
              </a-col>
            </a-row>

            <a-form-item>
              <a-space>
                <a-button type="primary" :loading="emailLoading" @click="saveEmailConfig">
                  保存配置
                </a-button>
                <a-button :loading="testEmailLoading" @click="testEmailConfigConnection">
                  测试连接
                </a-button>
                <a-button @click="testEmailModalVisible = true">
                  发送测试邮件
                </a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <!-- 存储配置 -->
        <a-tab-pane key="storage" tab="存储配置">
          <a-alert
            message="文件存储配置"
            description="配置文件存储方式，支持本地存储和云存储服务。"
            type="info"
            show-icon
            style="margin-bottom: 16px"
          />

          <a-form
            ref="storageFormRef"
            :model="storageConfig"
            :rules="storageRules"
            layout="vertical"
          >
            <a-form-item label="默认存储类型" name="default_type">
              <a-select v-model:value="storageConfig.default_type" placeholder="请选择默认存储类型">
                <a-select-option value="local">
                  本地存储
                </a-select-option>
                <a-select-option value="oss">
                  阿里云OSS
                </a-select-option>
                <a-select-option value="s3">
                  Amazon S3
                </a-select-option>
              </a-select>
            </a-form-item>

            <a-divider>本地存储配置</a-divider>

            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="上传目录" name="local_path">
                  <a-input v-model:value="storageConfig.local_path" placeholder="请输入上传目录" />
                </a-form-item>
              </a-col>

              <a-col :span="12">
                <a-form-item label="访问URL" name="local_url">
                  <a-input v-model:value="storageConfig.local_url" placeholder="请输入访问URL" />
                </a-form-item>
              </a-col>
            </a-row>

            <a-divider>阿里云OSS配置</a-divider>

            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="AccessKey ID" name="oss_access_key_id">
                  <a-input v-model:value="storageConfig.oss_access_key_id" placeholder="请输入AccessKey ID" />
                </a-form-item>
              </a-col>

              <a-col :span="12">
                <a-form-item label="AccessKey Secret" name="oss_access_key_secret">
                  <a-input-password v-model:value="storageConfig.oss_access_key_secret" placeholder="请输入AccessKey Secret" />
                </a-form-item>
              </a-col>
            </a-row>

            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="Endpoint" name="oss_endpoint">
                  <a-input v-model:value="storageConfig.oss_endpoint" placeholder="请输入Endpoint" />
                </a-form-item>
              </a-col>

              <a-col :span="12">
                <a-form-item label="Bucket" name="oss_bucket">
                  <a-input v-model:value="storageConfig.oss_bucket" placeholder="请输入Bucket名称" />
                </a-form-item>
              </a-col>
            </a-row>

            <a-divider>Amazon S3配置</a-divider>

            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="AccessKey ID" name="s3_access_key_id">
                  <a-input v-model:value="storageConfig.s3_access_key_id" placeholder="请输入AccessKey ID" />
                </a-form-item>
              </a-col>

              <a-col :span="12">
                <a-form-item label="AccessKey Secret" name="s3_access_key_secret">
                  <a-input-password v-model:value="storageConfig.s3_access_key_secret" placeholder="请输入AccessKey Secret" />
                </a-form-item>
              </a-col>
            </a-row>

            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="Region" name="s3_region">
                  <a-input v-model:value="storageConfig.s3_region" placeholder="请输入Region" />
                </a-form-item>
              </a-col>

              <a-col :span="12">
                <a-form-item label="Bucket" name="s3_bucket">
                  <a-input v-model:value="storageConfig.s3_bucket" placeholder="请输入Bucket名称" />
                </a-form-item>
              </a-col>
            </a-row>

            <a-form-item>
              <a-space>
                <a-button type="primary" :loading="storageLoading" @click="saveStorageConfig">
                  保存配置
                </a-button>
                <a-button :loading="testStorageLoading" @click="testStorageConfigConnection">
                  测试连接
                </a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- 测试邮件弹窗 -->
    <a-modal
      v-model:open="testEmailModalVisible"
      title="发送测试邮件"
      :confirm-loading="testEmailLoading"
      @ok="sendTestEmailFunc"
      @cancel="testEmailModalVisible = false"
    >
      <a-form
        :model="testEmailForm"
        :rules="testEmailRules"
        layout="vertical"
      >
        <a-form-item label="收件人邮箱" name="to">
          <a-input v-model:value="testEmailForm.to" placeholder="请输入收件人邮箱" />
        </a-form-item>

        <a-form-item label="邮件主题" name="subject">
          <a-input v-model:value="testEmailForm.subject" placeholder="请输入邮件主题" />
        </a-form-item>

        <a-form-item label="邮件内容" name="content">
          <a-textarea
            v-model:value="testEmailForm.content"
            placeholder="请输入邮件内容"
            :rows="6"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </page-container>
</template>

<style scoped>
</style>
