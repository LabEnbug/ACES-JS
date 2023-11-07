package database

import (
	"backend/auth"
	"backend/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
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

func StoreToken(token string, userId uint) error {
	// in fact, we could save a inverse token by using userId as key,
	// 20231107 edited, max 2 logged sessions, the older one will be revoked automatically
	expiration := config.TokenExpireTime // expire time setting
	userTokenKey := "user.token:" + strconv.Itoa(int(userId))

	var currentTokens *redis.IntCmd

	var err error

	_, err = rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		currentTokens = pipe.LLen(ctx, userTokenKey)
		err = pipe.RPush(ctx, userTokenKey, token).Err()
		if err != nil {
			return err
		}
		err = pipe.Expire(ctx, userTokenKey, expiration).Err()
		if err != nil {
			return err
		}
		err = pipe.Set(ctx, "token:"+token, "ok", expiration).Err()
		if err != nil {
			return err
		}
		return nil
	})

	tokensCount, err := currentTokens.Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if tokensCount >= config.TokenMaxDevice {
		_, err = rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			// pop left (first)
			oldToken, err := pipe.LPop(ctx, userTokenKey).Result()
			if err != nil && err != redis.Nil {
				return err
			}

			// delete selected token
			if oldToken != "" {
				err = pipe.Del(ctx, "token:"+oldToken).Err()
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	return err
}

func RevokeToken(token string) error {
	userId, _, err := auth.GetInfoFromToken(token)
	if err != nil {
		return err
	}
	userTokenKey := "user.token:" + strconv.Itoa(int(userId))

	_, err = rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		// delete token from user token list
		err := pipe.LRem(ctx, userTokenKey, 1, token).Err()
		if err != nil {
			return err
		}

		// delete token
		err = pipe.Del(ctx, "token:"+token).Err()
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func CheckTokenIsExist(token string) (bool, error) {
	// checking token could prevent token revoked by manager
	// or user (like logout all devices or change password),
	// but these features will not be implemented now
	ttl, err := rdb.TTL(ctx, "token:"+token).Result()
	if err != nil || ttl == -2 {
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
