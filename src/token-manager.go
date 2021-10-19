package src

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type TokenManager struct {
	tokenClient *TokenServer
	db          *redis.Client
}

func NewTokenManager(
	db *redis.Client,
	tokenClient *TokenServer,
) *TokenManager {
	return &TokenManager{
		db:          db,
		tokenClient: tokenClient,
	}
}

// 获取token
func (t *TokenManager) GetToken() (string, error) {
	ctx := context.TODO()
	token, err := t.db.Get(ctx, "token").Result()
	// token存在
	if len(token) > 0 {
		if err != nil {
			return "", err
		}
		return token, nil
	}
	// token为空，需要更新token
	id := uuid.New().String()
	t.db.Watch(ctx, t.UpdaterWatcher(ctx, id), "updater")
	// 若redis中的updater与自身提交的updater一致，说明具有更新token的权利，否在进入持续获取新token的阶段
	if v, _ := t.db.Get(ctx, "updater").Result(); v == id {
		fmt.Printf("routine %v is updating token\n", id)
		token := t.tokenClient.NewToken()
		t.db.TxPipelined(ctx, func(p redis.Pipeliner) error {
			// 假设令牌5秒后会自动过期，因此我们为redis中存在的令牌设置同等时间的有效期
			p.Set(ctx, "token", token, time.Second*5)
			p.Del(ctx, "updater")
			return nil
		})
	}
	// 等待并获取新的token
	for {
		temp, _ := t.db.Get(ctx, "token").Result()
		if temp != token && len(temp) != 0 {
			return temp, nil
		}
		time.Sleep(time.Millisecond * 1)
	}
}

// 乐观锁更新对象监视器
//
// ctx: 上下文
//
// id: 更新者id
func (TokenManager) UpdaterWatcher(
	ctx context.Context,
	id string,
) func(tx *redis.Tx) error {
	return func(tx *redis.Tx) error {
		key := "updater"
		updater, err := tx.Get(ctx, key).Result()
		if err != nil && err != redis.Nil {
			return err
		}
		// 若此时updater已存在，中断设置updater
		if len(updater) > 0 {
			return errors.New("error: token is updating by other routine")
		}
		// 设置updater
		updater = id
		_, err = tx.Pipelined(ctx, func(p redis.Pipeliner) error {
			p.Set(ctx, "updater", updater, 0)
			return nil
		})
		return err
	}
}
