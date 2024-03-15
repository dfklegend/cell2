package dbop

// 避免多个服务器key冲突
var redisKeyPrefix = ""

func SetRedisKeyPrefix(prefix string) {
	redisKeyPrefix = prefix + "-"
}

func MakeRedisKey(key string) string {
	return redisKeyPrefix + key
}
