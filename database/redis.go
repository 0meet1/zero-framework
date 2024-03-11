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

type RedisKeyspaceExpiredObserver interface {
	OnMessage(*redis.Message) error
}

func InitRedisDatabase(observers ...RedisKeyspaceExpiredObserver) {

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
			keeper.init(failoverClient, observers...)
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
			keeper.init(failoverClient, observers...)
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
			keeper.init(client, observers...)
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
			keeper.init(client, observers...)
			global.Key(DATABASE_REDIS, keeper)
		}
	}
}

type RedisKeeper interface {
	Client() *redis.Client
	Del(...string) error
	Set(string, string) error
	SetEx(string, string, int) error
	Get(string) (string, error)
}

type xRedisKeeper struct {
	redisContext context.Context
	redisClient  *redis.Client
	observers    []RedisKeyspaceExpiredObserver
}

func (xrk *xRedisKeeper) init(client *redis.Client, observers ...RedisKeyspaceExpiredObserver) {
	xrk.redisContext = context.Background()
	client.Ping(xrk.redisContext)
	xrk.redisClient = client
	xrk.observers = observers
	if len(observers) > 0 {
		_, err := xrk.redisClient.ConfigSet(xrk.redisContext, "notify-keyspace-events", "Ex").Result()
		if err != nil {
			panic(err)
		}
		go func() {
			pubsub := client.Subscribe(xrk.redisContext, "__keyevent@0__:expired")
			defer pubsub.Close()
			for eventMessage := range pubsub.Channel() {
				for _, observer := range xrk.observers {
					err := observer.OnMessage(eventMessage)
					if err != nil {
						global.Logger().Error(fmt.Sprintf("redis observer error: %s", err.Error()))
					}
				}
			}
		}()
	}
}

func (xrk *xRedisKeeper) Client() *redis.Client {
	return xrk.redisClient
}

func (xrk *xRedisKeeper) Del(key ...string) error {
	return xrk.redisClient.Del(xrk.redisContext, key...).Err()
}

func (xrk *xRedisKeeper) Set(key string, value string) error {
	return xrk.redisClient.Set(xrk.redisContext, key, value, 0).Err()
}

func (xrk *xRedisKeeper) SetEx(key string, value string, interval int) error {
	return xrk.redisClient.SetEX(xrk.redisContext, key, value, time.Duration(interval)*time.Second).Err()
}

func (xrk *xRedisKeeper) SetNX(key string, value string, interval int) error {
	return xrk.redisClient.SetNX(xrk.redisContext, key, value, time.Duration(interval)*time.Second).Err()
}

func (xrk *xRedisKeeper) Expire(key string, interval int) error {
	return xrk.redisClient.Expire(xrk.redisContext, key, time.Duration(interval)*time.Second).Err()
}

func (xrk *xRedisKeeper) Exists(key string) (int64, error) {
	return xrk.redisClient.Exists(xrk.redisContext, key).Result()
}

func (xrk *xRedisKeeper) Get(key string) (string, error) {
	vcount, err := xrk.Exists(key)
	if err != nil {
		return "", err
	}
	if vcount <= 0 {
		return "", nil
	}
	return xrk.redisClient.Get(xrk.redisContext, key).Result()
}
