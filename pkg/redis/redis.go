package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/usual2970/retell/pkg/config"
)

// RedisManager Redis连接管理器
type RedisManager struct {
	clients map[string]*redis.Client
	mutex   sync.RWMutex
}

var (
	instance *RedisManager
	once     sync.Once
)

// getInstance 获取 RedisManager 单例
func getInstance() *RedisManager {
	once.Do(func() {
		instance = &RedisManager{
			clients: make(map[string]*redis.Client),
		}
	})
	return instance
}

// Init 初始化所有Redis连接
func Init() error {
	manager := getInstance()
	conf := config.GetConfig()

	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	for name, redisConf := range conf.Redises {
		poolSize := 10
		if redisConf.PoolSize > 0 {
			poolSize = redisConf.PoolSize
		}
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", redisConf.Host, redisConf.Port),
			Password: redisConf.Password,
			DB:       redisConf.DB,
			// 连接池设置
			PoolSize: poolSize,
		})

		// 测试连接
		ctx := context.Background()
		if err := client.Ping(ctx).Err(); err != nil {
			return fmt.Errorf("failed to connect redis %s: %v", name, err)
		}

		manager.clients[name] = client
	}

	return nil
}

// GetRedis 获取指定名称的Redis连接
func GetRedis(name string) (*redis.Client, error) {
	manager := getInstance()
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	client, ok := manager.clients[name]
	if !ok {
		return nil, fmt.Errorf("redis %s not found", name)
	}
	return client, nil
}

func GetPaasRedis() (*redis.Client, error) {
	return GetRedis("paas")
}

// MustGetRedis 获取指定名称的Redis连接，如果不存在则panic
func MustGetRedis(name string) *redis.Client {
	client, err := GetRedis(name)
	if err != nil {
		panic(err)
	}
	return client
}

// Close 关闭所有Redis连接
func Close() error {
	manager := getInstance()
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	for name, client := range manager.clients {
		if err := client.Close(); err != nil {
			return fmt.Errorf("failed to close redis %s: %v", name, err)
		}
		delete(manager.clients, name)
	}
	return nil
}
