// 设置在线用户redis缓存
package models

import (
	"context"
	"ginchat/utils"
	"time"
)

func SetUserOnlineInfo(key string, val []byte, timeTTL time.Duration) {
	//用于管理 Redis 操作的超时、取消和其他相关信息
	ctx := context.Background()
	//设置redis的键值对并设置过期时间
	utils.Red.Set(ctx, key, val, timeTTL)
}
