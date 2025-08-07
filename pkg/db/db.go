package db

import (
	"fmt"
	"sync"

	"github.com/usual2970/retell/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBManager 数据库连接管理器
type DBManager struct {
	dbs   map[string]*gorm.DB
	mutex sync.RWMutex
}

var (
	instance *DBManager
	once     sync.Once
)

// getInstance 获取 DBManager 单例
func getInstance() *DBManager {
	once.Do(func() {
		instance = &DBManager{
			dbs: make(map[string]*gorm.DB),
		}
	})
	return instance
}

// Init 初始化所有数据库连接
func Init() error {
	manager := getInstance()
	conf := config.GetConfig()

	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	for name, dbConf := range conf.Databases {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
			dbConf.User,
			dbConf.Password,
			dbConf.Host,
			dbConf.Port,
			dbConf.Database,
		)

		gormConf := &gorm.Config{}
		if config.IsDevelopment() {
			gormConf = &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			}
		}

		db, err := gorm.Open(mysql.Open(dsn), gormConf)
		if err != nil {
			return fmt.Errorf("failed to connect database %s: %v", name, err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("failed to get database instance %s: %v", name, err)
		}

		maxIdleConns := 10
		if dbConf.MaxIdleConns > 0 {
			maxIdleConns = dbConf.MaxIdleConns
		}
		maxOpenConns := 100
		if dbConf.MaxOpenConns > 0 {
			maxOpenConns = dbConf.MaxOpenConns
		}

		// 设置连接池参数
		sqlDB.SetMaxIdleConns(maxIdleConns)
		sqlDB.SetMaxOpenConns(maxOpenConns)

		manager.dbs[name] = db
	}

	return nil
}

// GetDB 获取指定名称的数据库连接
func GetDB(name string) (*gorm.DB, error) {
	manager := getInstance()
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	db, ok := manager.dbs[name]
	if !ok {
		return nil, fmt.Errorf("database %s not found", name)
	}
	return db, nil
}

func GetPaasDB() (*gorm.DB, error) {
	return GetDB("paas")
}

// MustGetDB 获取指定名称的数据库连接，如果不存在则panic
func MustGetDB(name string) *gorm.DB {
	db, err := GetDB(name)
	if err != nil {
		panic(err)
	}
	return db
}

// Close 关闭所有数据库连接
func Close() error {
	manager := getInstance()
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	for name, db := range manager.dbs {
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("failed to get database instance %s: %v", name, err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close database %s: %v", name, err)
		}
		delete(manager.dbs, name)
	}
	return nil
}
