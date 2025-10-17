package cache

import "time"

// 定义默认过期时间
const (
	DefaultExpiration = 5 * time.Minute  // 默认5分钟过期
	ShortExpiration   = 1 * time.Minute  // 短期1分钟过期
	LongExpiration    = 30 * time.Minute // 长期30分钟过期
	NeverExpire       = 0 * time.Second  // 永不过期
)
