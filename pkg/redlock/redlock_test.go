package redlock

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestLock(t *testing.T) {
	client := getTestRedis()
	locker, _ := NewLocker(client)

	rs, err := locker.Lock(context.Background(), "test3", time.Hour*10)

	t.Log(rs, err)

	rs1, err1 := locker.Lock(context.Background(), "test3", time.Hour*10)
	t.Log(rs1, err1)
}

func TestUnLock(t *testing.T) {
	client := getTestRedis()

	rs, err := client.Get(context.Background(), "test113").Result()
	t.Log(rs, err)

	locker, _ := NewLocker(client)

	rs, err = locker.Lock(context.Background(), "test9", time.Hour*10)

	t.Log(rs, err)

	err = locker.UnLock(context.Background(), "test9", rs)
	t.Log(err)

	rs1, err1 := locker.Lock(context.Background(), "test5", time.Hour*10)
	t.Log(rs1, err1)
}

func getTestRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", "127.0.0.1", "6379"),
		Password: "password",
		DB:       0,
		// 连接池设置
		PoolSize: 10,
	})

	return client

}
