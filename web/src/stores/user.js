import { getUserInfoApi } from '~@/api/common/user'
import { defineStore } from 'pinia'
import { useStorage } from '@vueuse/core'
import { useAuthorization } from '~/composables/authorization'

export const useUserStore = defineStore('user', () => {
  // 用户信息存储在localStorage中
  const user = useStorage('user', null)

  // 设置用户信息
  const setUser = (userData) => {
    user.value = userData
  }

  // 清除用户信息
  const logout = () => {
    user.value = null
  }

  // 检查是否已登录
  const isLoggedIn = computed(() => {
    return !!user.value
  })

  // 获取用户信息
  const getUserInfo = async () => {
    try {
      const res = await getUserInfoApi()
      if (res.code === 0) {
        setUser(res.data?.user || null)
        return res.data?.user || null
      }
      return null
    } catch (error) {
      console.error('获取用户信息失败:', error)
      return null
    }
  }

  return {
    user,
    setUser,
    logout,
    isLoggedIn,
    getUserInfo
  }
})