<script setup>
import { ref, onMounted } from 'vue'
import { useRequest } from 'vue-request'
import axios from 'axios'
import { Card, Form, Input, Button, message, Row, Col, Divider } from 'ant-design-vue'
import { useUserStore } from '~/stores/user'

const userStore = useUserStore()

// 个人资料表单
const profileForm = ref({
  id: '',
  username: '',
  nickname: '',
  email: '',
  phone: '',
  created_at: '',
  updated_at: ''
})

// 密码表单
const passwordForm = ref({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

// 表单加载状态
const loading = ref(false)
const passwordLoading = ref(false)

// 表单引用
const profileFormRef = ref()
const passwordFormRef = ref()

// 获取当前用户个人资料
const fetchProfile = async () => {
  try {
    // 添加调试信息
    console.log('Fetching profile...')
    const res = await axios.get('/api/admin/profile')
    
    if (res.data.code === 0) {
      profileForm.value = {
        ...profileForm.value,
        ...res.data.data
      }
    } else {
      message.error(res.data.msg || '获取个人资料失败')
    }
  } catch (error) {
    console.error('Failed to fetch profile:', error)
    message.error('获取个人资料失败')
  }
}

// 更新个人资料
const updateProfile = async () => {
  try {
    await profileFormRef.value.validate()
    
    loading.value = true
    
    const data = {
      nickname: profileForm.value.nickname,
      email: profileForm.value.email,
      phone: profileForm.value.phone
    }
    
    const res = await axios.put('/api/admin/profile', data)
    
    if (res.data.code === 0) {
      message.success('个人资料更新成功')
      // 更新用户存储中的信息
      userStore.setUser({
        ...userStore.user,
        nickname: profileForm.value.nickname,
        email: profileForm.value.email,
        phone: profileForm.value.phone
      })
    } else {
      message.error(res.data.msg || '更新个人资料失败')
    }
  } catch (error) {
    console.error('Failed to update profile:', error)
    message.error('更新个人资料失败')
  } finally {
    loading.value = false
  }
}

// 更新密码表单规则
const passwordRules = ref({
  old_password: [
    { required: true, message: '请输入旧密码' }
  ],
  new_password: [
    { required: true, message: '请输入新密码' },
    { min: 6, message: '密码长度不能少于6位' }
  ],
  confirm_password: [
    { required: true, message: '请确认新密码' },
    {
      validator: (_, value) => {
        if (value && value !== passwordForm.value.new_password) {
          return Promise.reject('两次输入的密码不一致')
        }
        return Promise.resolve()
      }
    }
  ]
})

// 更新密码
const updatePassword = async () => {
  try {
    await passwordFormRef.value.validate()
    
    passwordLoading.value = true
    
    const data = {
      old_password: passwordForm.value.old_password,
      new_password: passwordForm.value.new_password
    }
    
    const res = await axios.put('/api/admin/profile/password', data)
    
    if (res.data.code === 0) {
      message.success('密码更新成功，请重新登录')
      // 清空表单
      passwordForm.value = {
        old_password: '',
        new_password: '',
        confirm_password: ''
      }
      // 1秒后跳转到登录页
      setTimeout(() => {
        userStore.logout()
        window.location.href = '/login'
      }, 1000)
    } else {
      message.error(res.data.msg || '更新密码失败')
    }
  } catch (error) {
    console.error('Failed to update password:', error)
    message.error('更新密码失败')
  } finally {
    passwordLoading.value = false
  }
}

onMounted(() => {
  fetchProfile()
})
</script>

<template>
  <page-container title="个人资料">
    <a-row :gutter="16">
      <a-col :span="12">
        <a-card title="基本信息">
          <a-form
            ref="profileFormRef"
            :model="profileForm"
            :label-col="{ span: 6 }"
            :wrapper-col="{ span: 18 }"
          >
            <a-form-item label="用户ID">
              <a-input v-model:value="profileForm.id" disabled />
            </a-form-item>
            
            <a-form-item label="用户名">
              <a-input v-model:value="profileForm.username" disabled />
            </a-form-item>
            
            <a-form-item label="昵称" name="nickname">
              <a-input v-model:value="profileForm.nickname" />
            </a-form-item>
            
            <a-form-item label="邮箱" name="email">
              <a-input v-model:value="profileForm.email" type="email" />
            </a-form-item>
            
            <a-form-item label="手机号" name="phone">
              <a-input v-model:value="profileForm.phone" />
            </a-form-item>
            
            <a-form-item label="创建时间">
              <a-input 
                :value="profileForm.created_at ? new Date(profileForm.created_at).toLocaleString() : ''" 
                disabled 
              />
            </a-form-item>
            
            <a-form-item label="更新时间">
              <a-input 
                :value="profileForm.updated_at ? new Date(profileForm.updated_at).toLocaleString() : ''" 
                disabled 
              />
            </a-form-item>
            
            <a-form-item :wrapper-col="{ span: 18, offset: 6 }">
              <a-button type="primary" :loading="loading" @click="updateProfile">
                更新资料
              </a-button>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>
      
      <a-col :span="12">
        <a-card title="修改密码">
          <a-form
            ref="passwordFormRef"
            :model="passwordForm"
            :rules="passwordRules"
            :label-col="{ span: 6 }"
            :wrapper-col="{ span: 18 }"
          >
            <a-form-item label="旧密码" name="old_password">
              <a-input 
                v-model:value="passwordForm.old_password" 
                type="password"
                autocomplete="current-password"
              />
            </a-form-item>
            
            <a-form-item label="新密码" name="new_password">
              <a-input 
                v-model:value="passwordForm.new_password" 
                type="password"
                autocomplete="new-password"
              />
            </a-form-item>
            
            <a-form-item label="确认密码" name="confirm_password">
              <a-input 
                v-model:value="passwordForm.confirm_password" 
                type="password"
                autocomplete="new-password"
              />
            </a-form-item>
            
            <a-form-item :wrapper-col="{ span: 18, offset: 6 }">
              <a-button type="primary" :loading="passwordLoading" @click="updatePassword">
                更新密码
              </a-button>
            </a-form-item>
          </a-form>
        </a-card>
        
        <a-card title="安全提示" style="margin-top: 16px">
          <p>为了您的账户安全，请注意以下事项：</p>
          <ul>
            <li>定期更新密码</li>
            <li>不要使用简单密码</li>
            <li>不要在公共场合泄露个人信息</li>
            <li>发现异常登录请及时修改密码</li>
          </ul>
        </a-card>
      </a-col>
    </a-row>
  </page-container>
</template>

<style scoped>
ul {
  padding-left: 20px;
}

ul li {
  margin-bottom: 8px;
}
</style>