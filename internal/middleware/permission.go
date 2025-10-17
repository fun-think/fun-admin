package middleware

import (
	"fun-admin/pkg"
	"fun-admin/pkg/jwt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gin-gonic/gin"
)

// PermissionMiddleware 权限中间件，用于控制 admin 资源的访问权限
func PermissionMiddleware(e *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从上下文获取用户信息（假设通过 JWT 或其他方式设置）
		v, exists := ctx.Get("claims")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权访问",
			})
			ctx.Abort()
			return
		}

		uid := v.(*jwt.MyCustomClaims).UserId

		// 超级管理员跳过权限检查
		if convertor.ToString(uid) == pkg.AdminUserID {
			ctx.Next()
			return
		}

		// 获取请求的资源和操作
		sub := convertor.ToString(uid)

		// 构造资源路径（例如：api:/api/v1/users）
		obj := pkg.ApiResourcePrefix + ctx.Request.URL.Path

		// 获取 HTTP 方法
		act := ctx.Request.Method

		// 检查权限
		allowed, err := e.Enforce(sub, obj, act)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "权限检查失败",
			})
			ctx.Abort()
			return
		}

		if !allowed {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
