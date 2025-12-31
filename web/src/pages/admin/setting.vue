<script setup>
import { onMounted, ref } from 'vue'
import { message } from 'ant-design-vue'

// 表单数据
const formModel = ref({
  site_name: '',
  site_description: '',
  enable_registration: false,
  default_role: 'user',
  timezone: 'Asia/Shanghai',
})

// 表单加载状态
const loading = ref(false)

// 表单规则
const formRules = ref({
  site_name: [
    { required: true, message: '请输入网站名称' },
  ],
  default_role: [
    { required: true, message: '请选择默认角色' },
  ],
  timezone: [
    { required: true, message: '请选择时区' },
  ],
})

// 表单引用
const formRef = ref()

// 获取当前语言
function getCurrentLanguage() {
  return localStorage.getItem('admin-language') || 'zh-CN'
}

// 获取设置数据
async function fetchSettings() {
  loading.value = true
  try {
    // 模拟获取设置数据
    // 实际项目中应该调用后端 API 获取真实数据
    setTimeout(() => {
      formModel.value = {
        site_name: 'Fun Admin',
        site_description: '一个功能强大的后台管理系统',
        enable_registration: true,
        default_role: 'user',
        timezone: 'Asia/Shanghai',
      }
      loading.value = false
    }, 500)
  }
  catch (error) {
    console.error('Failed to fetch settings:', error)
    message.error('获取设置失败')
    loading.value = false
  }
}

// 保存设置
async function saveSettings() {
  try {
    await formRef.value.validate()

    loading.value = true

    // 模拟保存设置
    // 实际项目中应该调用后端 API 保存数据
    setTimeout(() => {
      loading.value = false
      message.success('设置保存成功')
    }, 1000)
  }
  catch (error) {
    console.error('Failed to save settings:', error)
    message.error('保存设置失败')
  }
}

// 重置表单
function resetForm() {
  formRef.value.resetFields()
}

onMounted(() => {
  fetchSettings()
})
</script>

<template>
  <page-container title="系统设置">
    <a-card>
      <a-form
        ref="formRef"
        :model="formModel"
        :rules="formRules"
        :label-col="{ span: 4 }"
        :wrapper-col="{ span: 14 }"
      >
        <a-form-item label="网站名称" name="site_name">
          <a-input
            v-model:value="formModel.site_name"
            placeholder="请输入网站名称"
          />
        </a-form-item>

        <a-form-item label="网站描述" name="site_description">
          <a-textarea
            v-model:value="formModel.site_description"
            placeholder="请输入网站描述"
            :rows="4"
          />
        </a-form-item>

        <a-form-item label="允许注册" name="enable_registration">
          <a-switch
            v-model:checked="formModel.enable_registration"
            checked-children="是"
            un-checked-children="否"
          />
        </a-form-item>

        <a-form-item label="默认角色" name="default_role">
          <a-select
            v-model:value="formModel.default_role"
            placeholder="请选择默认角色"
          >
            <a-select-option value="user">
              普通用户
            </a-select-option>
            <a-select-option value="guest">
              访客
            </a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item label="时区" name="timezone">
          <a-select
            v-model:value="formModel.timezone"
            placeholder="请选择时区"
          >
            <a-select-option value="UTC">
              UTC
            </a-select-option>
            <a-select-option value="Asia/Shanghai">
              亚洲/上海
            </a-select-option>
            <a-select-option value="America/New_York">
              美洲/纽约
            </a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
          <a-button
            type="primary"
            :loading="loading"
            @click="saveSettings"
          >
            保存设置
          </a-button>
          <a-button
            style="margin-left: 10px"
            @click="resetForm"
          >
            重置
          </a-button>

          <!-- 语言切换 -->
          <a-dropdown style="margin-left: 10px">
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
        </a-form-item>
      </a-form>
    </a-card>
  </page-container>
</template>
