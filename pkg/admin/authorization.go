package admin

import "context"

// Authorizable 提供资源级/动作级鉴权的可选接口
// 返回 nil 表示允许，返回错误表示拒绝并携带原因
// 具体实现可基于 Casbin/自定义策略
// 注意：所有方法均为可选，未实现则默认放行

type Authorizable interface {
	CanList(ctx context.Context) error
	CanView(ctx context.Context, id interface{}) error
	CanCreate(ctx context.Context, data map[string]interface{}) error
	CanUpdate(ctx context.Context, id interface{}, data map[string]interface{}) error
	CanDelete(ctx context.Context, id interface{}) error
}
