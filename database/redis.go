package database

import (
	"context"
	"strconv"
	"time"
	"zero-framework/global"

	"github.com/go-redis/redis/v8"
)

const (
	DATABASE_REDIS = "zero.database.redis"
)

func InitRedisDatabase() {

	if len(global.StringValue("zero.redis.sentinel")) > 0 {
		if len(global.StringValue("zero.redis.password")) > 0 {
			failoverClient := redis.NewFailoverClient(&redis.FailoverOptions{
				MasterName:    global.StringValue("zero.redis.sentinel"),
				SentinelAddrs: []string{global.StringValue("zero.redis.hostname") + ":" + strconv.Itoa(global.IntValue("zero.redis.hostport"))},
				Password:      global.StringValue("zero.redis.password"),
				DB:            global.IntValue("zero.redis.database"),
				IdleTimeout:   time.Duration(global.IntValue("zero.redis.idleTimeout")),
				PoolSize:      global.IntValue("zero.redis.maxActive"),
				MinIdleConns:  global.IntValue("zero.redis.maxIdle"),
			})

			keeper := &xRedisKeeper{}
			keeper.init(failoverClient)
			global.Key(DATABASE_REDIS, keeper)
		} else {
			failoverClient := redis.NewFailoverClient(&redis.FailoverOptions{
				MasterName:    global.StringValue("zero.redis.sentinel"),
				SentinelAddrs: []string{global.StringValue("zero.redis.hostname") + ":" + strconv.Itoa(global.IntValue("zero.redis.hostport"))},
				DB:            global.IntValue("zero.redis.database"),
				IdleTimeout:   time.Duration(global.IntValue("zero.redis.idleTimeout")),
				PoolSize:      global.IntValue("zero.redis.maxActive"),
				MinIdleConns:  global.IntValue("zero.redis.maxIdle"),
			})

			keeper := &xRedisKeeper{}
			keeper.init(failoverClient)
			global.Key(DATABASE_REDIS, keeper)
		}
	} else {
		if len(global.StringValue("zero.redis.password")) > 0 {
			client := redis.NewClient(&redis.Options{
				Addr:         global.StringValue("zero.redis.hostname") + ":" + strconv.Itoa(global.IntValue("zero.redis.hostport")),
				Password:     global.StringValue("zero.redis.password"),
				DB:           global.IntValue("zero.redis.database"),
				IdleTimeout:  time.Duration(global.IntValue("zero.redis.idleTimeout")),
				PoolSize:     global.IntValue("zero.redis.maxActive"),
				MinIdleConns: global.IntValue("zero.redis.maxIdle"),
			})

			keeper := &xRedisKeeper{}
			keeper.init(client)
			global.Key(DATABASE_REDIS, keeper)
		} else {
			client := redis.NewClient(&redis.Options{
				Addr:         global.StringValue("zero.redis.hostname") + ":" + strconv.Itoa(global.IntValue("zero.redis.hostport")),
				DB:           global.IntValue("zero.redis.database"),
				IdleTimeout:  time.Duration(global.IntValue("zero.redis.idleTimeout")),
				PoolSize:     global.IntValue("zero.redis.maxActive"),
				MinIdleConns: global.IntValue("zero.redis.maxIdle"),
			})

			keeper := &xRedisKeeper{}
			keeper.init(client)
			global.Key(DATABASE_REDIS, keeper)
		}
	}
}

type xRedisKeeper struct {
	redisContext context.Context
	redisClient  *redis.Client
}

func (xrk *xRedisKeeper) init(client *redis.Client) {
	xrk.redisContext = context.Background()
	client.Ping(xrk.redisContext)
	xrk.redisClient = client
}

func (xrk *xRedisKeeper) RedisSetEx(key string, value string, interval int) error {
	cmd := xrk.redisClient.Do(xrk.redisContext, "setex", key, strconv.Itoa(interval+60), value)
	return cmd.Err()
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
