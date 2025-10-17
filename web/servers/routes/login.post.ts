export default eventHandler(async (event) => {
  const body = await readBody(event)
  const { username, password } = body
  
  // 使用环境变量或默认值
  const apiUrl = `http://127.0.0.1:8001/api/v1/login`
  
  try {
    // 调用真实的后端登录接口
    const response = await fetch(apiUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username,
        password
      })
    })
    
    const result = await response.json()
    
    if (result.code === 0) {
      // 登录成功，返回真实的访问令牌
      return {
        code: 200,
        data: {
          token: result.data.accessToken,
        },
        msg: '登录成功',
      }
    } else {
      // 登录失败
      setResponseStatus(event, 401)
      return {
        code: 401,
        msg: result.message || '用户名或密码错误',
      }
    }
  } catch (error) {
    // 处理网络或其他错误
    setResponseStatus(event, 500)
    return {
      code: 500,
      msg: '服务器内部错误',
    }
  }
})