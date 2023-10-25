package database

import (
	"backend/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var rdb *redis.Client
var ctx = context.Background()

func InitRedisPool() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.DefaultRedis.Host, config.DefaultRedis.Port),
		Password: config.DefaultRedis.Pass,
		DB:       config.DefaultRedis.Channel,
		PoolSize: 10,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}

func CloseRedisPool() {
	err := rdb.Close()
	if err != nil {
		return
	}
}

func StoreToken(token string) error {
	// in fact, we could save a inverse token by using userId as key,
	// but it's not necessary to do it now
	expiration := time.Duration(72) * time.Hour // 3 days valid
	err := rdb.Set(ctx, "token:"+token, "ok", expiration).Err()
	return err
}

func RevokeToken(token string) error {
	err := rdb.Del(ctx, "token:"+token).Err()
	return err
}

func CheckTokenIsExist(token string) (bool, error) {
	// checking token could prevent token revoked by manager
	// or user (like logout all devices or change password),
	// but these features will not be implemented now
	ttl, err := rdb.TTL(ctx, "token:"+token).Result()
	if err != nil {
		return false, err
	}

	if ttl == -2 {
		return false, fmt.Errorf("token expired")
	}

	_, err = rdb.Get(ctx, "token:"+token).Result()

	if err == redis.Nil {
		return false, fmt.Errorf("token not exist")
	} else if err != nil {
		return false, err
	}

	return true, nil
}
