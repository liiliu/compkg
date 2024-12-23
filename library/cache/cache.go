package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
	"weihu_server/library/common"
	"weihu_server/library/config"
	"weihu_server/library/logger"
)

var ctx = context.Background()
var redisClient *redis.Client
var redisPrefix string

func Initial() {
	config.MustHas("server.name")
	config.MustHas("redis.idleTimeout")
	config.MustHas("redis.poolSize")
	config.MustHas("redis.host")

	logger.Info(common.LogTagRedis, "Initial redis ...")

	redisPrefix = config.GetString("redis.keyPrefix")

	redisClient = redis.NewClient(&redis.Options{
		Addr:        config.GetString("redis.host"),
		Password:    config.GetString("redis.password"),
		DB:          config.GetInt("redis.db"),
		IdleTimeout: config.GetDuration("redis.idleTimeout") * time.Second,
		PoolSize:    config.GetInt("redis.poolSize"),
		MaxConnAge:  config.GetDuration("redis.maxConnAge") * time.Second,
	})

	// go-redis库v8版本相关命令都需要传递context.Context参数,Background 返回一个非空的Context,它永远不会被取消，没有值，也没有期限。
	ctx1, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := redisClient.Ping(ctx1).Result()
	if err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return
	}
	logger.Info(common.LogTagRedis, pong)
}

// GetClient 获取客户端
func GetClient() *redis.Client {
	return redisClient
}

// Remove 删除
func Remove(name string) error {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	_, err := redisClient.Del(ctx, key).Result()
	if err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return errors.New(err.Error())
	}
	return nil
}

// Save 保存
func Save(name string, values interface{}, timeout time.Duration) error {
	serialized, err := json.Marshal(values)
	if err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return errors.New(err.Error())
	}
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	_, err = redisClient.Set(ctx, key, serialized, timeout).Result()
	if err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return errors.New(err.Error())
	}
	return nil
}

// Expire 设置过期时间
func Expire(name string, timeout time.Duration) error {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	_, err := redisClient.Expire(ctx, key, timeout).Result()
	if err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return errors.New(err.Error())
	}
	return nil
}

// Exists 是否存在
func Exists(name string) bool {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	result, err := redisClient.Exists(ctx, key).Result()
	if err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return false
	}
	fmt.Printf("result:%d\n", result)
	if result == 0 {
		return false
	} else {
		return true
	}
}

// Get 获取
func Get(name string, values interface{}) error {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	value, err := redisClient.Get(ctx, key).Bytes()
	if err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return errors.New(err.Error())
	}

	err = json.Unmarshal(value, &values)
	if err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return errors.New(err.Error())
	}

	return nil
}

// SaveString 保存字符串
func SaveString(name string, str string, timeout time.Duration) error {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	_, err := redisClient.Set(ctx, key, str, timeout).Result()
	if err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return errors.New(err.Error())
	}

	return nil
}

// GetString 获取字符串
func GetString(name string) string {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		//logger.Error(common.LogTagRedis, err.Error())
		return ""
	}

	return value
}

// Llen 返回列表长度
func Llen(name string) (int64, error) {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	value, err := redisClient.LLen(ctx, key).Result()
	if err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return 0, errors.New(err.Error())
	}
	return value, nil
}

// Lpop 移除并返回列表
func Lpop(name string) ([]byte, error) {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	value, err := redisClient.LPop(ctx, key).Bytes()
	if err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return nil, errors.New(err.Error())
	}
	return value, nil
}

// Rpush 将一个或多个值插入到列表的尾部
func Rpush(name string, values interface{}) error {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	_, err := redisClient.RPush(ctx, key, values).Result()
	if err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return errors.New(err.Error())
	}
	return nil
}

// TTL 返回剩余时间(秒)
func TTL(name string) (time.Duration, error) {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	value, err := redisClient.TTL(ctx, key).Result()
	if err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return 0, errors.New(err.Error())
	}
	return value, nil
}

// ZAdd 为有序集合添加元素
func ZAdd(name string, score int64, value interface{}) error {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	_, err := redisClient.ZAdd(ctx, key, &redis.Z{
		Score:  float64(score),
		Member: value,
	}).Result()
	if err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return errors.New(err.Error())
	}
	return nil
}

// ZRange 获取有序集合所有元素
func ZRange(name string) ([]string, error) {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	value, err := redisClient.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result()
	if err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return nil, errors.New(err.Error())
	}
	return value, nil
}

// ZCount 获取有序集合指定分数范围内的元素数量
func ZCount(name string, min, max int64) (int64, error) {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	value, err := redisClient.ZCount(ctx, key, strconv.FormatInt(min, 10), strconv.FormatInt(max, 10)).Result()
	if err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return 0, errors.New(err.Error())
	}
	return value, nil
}

// ZRemRangeByScore 删除有序集合指定分数范围的元素数量
func ZRemRangeByScore(name string, min, max int64) error {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	_, err := redisClient.ZRemRangeByScore(ctx, key, strconv.FormatInt(min, 10), strconv.FormatInt(max, 10)).Result()
	if err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return errors.New(err.Error())
	}
	return nil
}

// ZCard 获取有序集合总长度
func ZCard(name string) (int64, error) {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	value, err := redisClient.ZCard(ctx, key).Result()
	if err != nil {
		logger.Error(common.LogTagRedis, err.Error())
		return 0, errors.New(err.Error())
	}
	return value, nil
}
