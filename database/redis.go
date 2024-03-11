package database

import (
	"context"
	"fmt"
	"time"

	"github.com/0meet1/zero-framework/global"

	"github.com/go-redis/redis/v8"
)

const (
	DATABASE_REDIS = "zero.database.redis"
)

func InitRedisDatabase(hooks ...redis.Hook) {

	if len(global.StringValue("zero.redis.sentinel")) > 0 {
		if len(global.StringValue("zero.redis.password")) > 0 {
			failoverClient := redis.NewFailoverClient(&redis.FailoverOptions{
				MasterName:    global.StringValue("zero.redis.sentinel"),
				SentinelAddrs: global.SliceStringValue("zero.redis.hosts"),
				Password:      global.StringValue("zero.redis.password"),
				DB:            global.IntValue("zero.redis.database"),
				IdleTimeout:   time.Duration(global.IntValue("zero.redis.idleTimeout")),
				PoolSize:      global.IntValue("zero.redis.maxActive"),
				MinIdleConns:  global.IntValue("zero.redis.maxIdle"),
			})

			keeper := &xRedisKeeper{}
			keeper.init(failoverClient, hooks...)
			global.Key(DATABASE_REDIS, keeper)
		} else {

			failoverClient := redis.NewFailoverClient(&redis.FailoverOptions{
				MasterName:    global.StringValue("zero.redis.sentinel"),
				SentinelAddrs: global.SliceStringValue("zero.redis.hosts"),
				DB:            global.IntValue("zero.redis.database"),
				IdleTimeout:   time.Duration(global.IntValue("zero.redis.idleTimeout")),
				PoolSize:      global.IntValue("zero.redis.maxActive"),
				MinIdleConns:  global.IntValue("zero.redis.maxIdle"),
			})

			keeper := &xRedisKeeper{}
			keeper.init(failoverClient, hooks...)
			global.Key(DATABASE_REDIS, keeper)
		}
	} else {
		if len(global.StringValue("zero.redis.password")) > 0 {
			client := redis.NewClient(&redis.Options{
				Addr:         fmt.Sprintf("%s:%d", global.StringValue("zero.redis.hostname"), global.IntValue("zero.redis.hostport")),
				Password:     global.StringValue("zero.redis.password"),
				DB:           global.IntValue("zero.redis.database"),
				IdleTimeout:  time.Duration(global.IntValue("zero.redis.idleTimeout")),
				PoolSize:     global.IntValue("zero.redis.maxActive"),
				MinIdleConns: global.IntValue("zero.redis.maxIdle"),
			})

			keeper := &xRedisKeeper{}
			keeper.init(client, hooks...)
			global.Key(DATABASE_REDIS, keeper)
		} else {
			client := redis.NewClient(&redis.Options{
				Addr:         fmt.Sprintf("%s:%d", global.StringValue("zero.redis.hostname"), global.IntValue("zero.redis.hostport")),
				DB:           global.IntValue("zero.redis.database"),
				IdleTimeout:  time.Duration(global.IntValue("zero.redis.idleTimeout")),
				PoolSize:     global.IntValue("zero.redis.maxActive"),
				MinIdleConns: global.IntValue("zero.redis.maxIdle"),
			})

			keeper := &xRedisKeeper{}
			keeper.init(client, hooks...)
			global.Key(DATABASE_REDIS, keeper)
		}
	}
}

type RedisKeeper interface {
	RedisClient() *redis.Client
	RedisDel(key ...string) error
	RedisSet(key string, value string) error
	RedisSetEx(key string, value string, interval int) error
	RedisGet(key string) (string, error)
}

type xRedisKeeper struct {
	redisContext context.Context
	redisClient  *redis.Client
}

func (xrk *xRedisKeeper) init(client *redis.Client, hooks ...redis.Hook) {
	xrk.redisContext = context.Background()
	client.Ping(xrk.redisContext)
	xrk.redisClient = client
	if len(hooks) > 0 {
		for _, hook := range hooks {
			xrk.redisClient.AddHook(hook)
		}
	}
}

func (xrk *xRedisKeeper) RedisClient() *redis.Client {
	return xrk.redisClient
}

func (xrk *xRedisKeeper) RedisDel(key ...string) error {
	return xrk.redisClient.Del(xrk.redisContext, key...).Err()
}

func (xrk *xRedisKeeper) RedisSet(key string, value string) error {
	return xrk.redisClient.Set(xrk.redisContext, key, value, 0).Err()
}

func (xrk *xRedisKeeper) RedisSetEx(key string, value string, interval int) error {
	return xrk.redisClient.SetEX(xrk.redisContext, key, value, time.Duration(interval)*time.Second).Err()
}

func (xrk *xRedisKeeper) RedisGet(key string) (string, error) {
	cmd := xrk.redisClient.Do(xrk.redisContext, "exists", key)
	exists, err := cmd.Int()
	if err != nil {
		return "", err
	}

	if exists <= 0 {
		return "", nil
	}
	cmd = xrk.redisClient.Do(xrk.redisContext, "get", key)
	val, err := cmd.Text()
	if err != nil {
		return "", err
	}
	return val, nil
}
