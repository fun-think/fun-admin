package admin

import "context"

// CreateHook 资源创建生命周期钩子（可选）
// 返回 error 终止流程

type CreateHook interface {
	BeforeCreate(ctx context.Context, data map[string]interface{}) error
	AfterCreate(ctx context.Context, data map[string]interface{}) error
}

// UpdateHook 资源更新生命周期钩子（可选）

type UpdateHook interface {
	BeforeUpdate(ctx context.Context, id interface{}, data map[string]interface{}) error
	AfterUpdate(ctx context.Context, id interface{}, data map[string]interface{}) error
}

// DeleteHook 资源删除生命周期钩子（可选）

type DeleteHook interface {
	BeforeDelete(ctx context.Context, id interface{}) error
	AfterDelete(ctx context.Context, id interface{}) error
}
