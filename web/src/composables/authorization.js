export const STORAGE_AUTHORIZE_KEY = 'Authorization'

// 修改useAuthorization函数，确保正确处理Bearer token
export const useAuthorization = createGlobalState(() => {
  const token = useStorage(STORAGE_AUTHORIZE_KEY, null, undefined, { mergeDefaults: true })

  // 创建一个包装器来处理Bearer前缀
  return computed({
    get: () => {
      // 如果token不为空且不包含Bearer前缀，则添加前缀
      if (token.value && !token.value.startsWith('Bearer ')) {
        return `Bearer ${token.value}`
      }
      return token.value
    },
    set: (newToken) => {
      // 存储时不带Bearer前缀
      if (newToken && newToken.startsWith('Bearer ')) {
        token.value = newToken.substring(7) // 移除"Bearer "前缀
      }
      else {
        token.value = newToken
      }
    },
  })
})
