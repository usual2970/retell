package redlock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/xid"
	xredis "github.com/usual2970/retell/pkg/redis"
)

const sqLockPre = "paas:lock:%s"

var lock ILocker
var lockOnce sync.Once

func InitLock() {
	lockOnce.Do(func() {
		lock, _ = NewLocker()
	})
}

type ILocker interface {
	Lock(ctx context.Context, key string, expiration time.Duration) (string, error)
	UnLock(ctx context.Context, key, val string) error
}

type Locker struct {
	rc *redis.Client
}

func NewLocker(rc ...*redis.Client) (ILocker, error) {

	var redis *redis.Client
	if len(rc) > 0 {
		redis = rc[0]
	} else {
		var err error
		redis, err = xredis.GetPaasRedis()
		if err != nil {
			return nil, err
		}
	}
	lock = &Locker{
		rc: redis,
	}

	return lock, nil
}

// Lock 加锁
func (l *Locker) Lock(ctx context.Context, key string, expiration time.Duration) (string, error) {

	val := xid.New().String()
	key = getLockPrefix(key)

	rs, err := l.rc.SetNX(ctx, key, val, expiration).Result()
	if !rs || err != nil {
		return "", fmt.Errorf("获取redis lock失败：%v,%v", rs, err)
	}

	return val, nil
}

// UnLock 解锁
func (l *Locker) UnLock(ctx context.Context, key, val string) error {

	key = getLockPrefix(key)

	scriptStr := `
	if redis.call("get",KEYS[1]) == ARGV[1] then
    	return redis.call("del",KEYS[1])
	else
    	return 0
	end
	`
	_, err := redis.NewScript(scriptStr).Run(ctx, l.rc, []string{key}, val).Result()
	if err != nil {
		return fmt.Errorf("解锁失败：%s", err)
	}
	return nil
}

func Lock(ctx context.Context, key string, expiration time.Duration) (string, error) {
	return lock.Lock(ctx, key, expiration)
}

func UnLock(ctx context.Context, key, val string) error {
	return lock.UnLock(ctx, key, val)
}
func getLockPrefix(key string) string {
	return fmt.Sprintf(sqLockPre, key)
}
