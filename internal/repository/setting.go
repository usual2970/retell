package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/domain/constant"
	"github.com/usual2970/retell/pkg/db"
	"github.com/usual2970/retell/pkg/logger"
	xredis "github.com/usual2970/retell/pkg/redis"
	"gorm.io/gorm"
)

var settingCacheExp = time.Hour * 24

const MessageRetentionDaysKey = "message_retention_days"

type SettingRepository struct{}

func NewSettingRepository() *SettingRepository {
	return &SettingRepository{}
}

func (r *SettingRepository) getSettingKey(key string, userId ...int64) string {
	if len(userId) > 0 && userId[0] > 0 {
		return "paas:setting:key:" + key + ":user:" + strconv.FormatInt(userId[0], 10)
	}
	return "paas:setting:key:" + key
}

func (r *SettingRepository) SetSetting(ctx context.Context, setting *domain.Setting) error {

	// 更新数据库
	db, err := db.GetPaasDB()
	if err != nil {
		return err
	}

	if err := db.Save(setting).Error; err != nil {
		return err
	}

	rd, err := xredis.GetPaasRedis()
	if err != nil {
		return err
	}

	cacheKey := r.getSettingKey(setting.Key, setting.UserID)

	// 删除缓存
	if err := rd.Del(ctx, cacheKey).Err(); err != nil {
		logger.WithField("module", "repository").WithField("key", cacheKey).Error("Failed to delete cache")
	}

	return nil

}

func (r *SettingRepository) GetSetting(ctx context.Context, key string, userId ...int64) (*domain.Setting, error) {
	// 先从缓存中获取
	rd, err := xredis.GetPaasRedis()
	if err != nil {
		return nil, err
	}

	cacheKey := r.getSettingKey(key, userId...)

	rsBts, err := rd.Get(ctx, cacheKey).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if string(rsBts) == notFoundPlaceHolder {
		return nil, constant.ErrRecordNotFound
	}

	if len(rsBts) > 0 {
		setting := &domain.Setting{}
		if err := json.Unmarshal(rsBts, setting); err == nil {
			return setting, nil
		}
		// 若解析缓存数据失败，继续从数据库获取
	}

	// 缓存中不存在或解析失败，从数据库获取
	db, err := db.GetPaasDB()
	if err != nil {
		return nil, err
	}

	setting := &domain.Setting{}
	query := db.Where("`key` = ?", key)
	if len(userId) > 0 {
		query = query.Where("user_id = ?", userId[0])
	}
	if err := query.First(setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 缓存空值
			rd.Set(ctx, cacheKey, notFoundPlaceHolder, settingCacheExp)
			return nil, constant.ErrRecordNotFound
		}
		return nil, err
	}

	// 缓存结果
	settingBts, err := json.Marshal(setting)
	if err != nil {
		return nil, err
	}

	rd.Set(ctx, cacheKey, string(settingBts), settingCacheExp)

	return setting, nil
}
