package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/domain/constant"
	"github.com/usual2970/retell/pkg/db"
	xredis "github.com/usual2970/retell/pkg/redis"
	"gorm.io/gorm"

	"github.com/redis/go-redis/v9"
)

var accountCacheExp = time.Hour * 24

type UserAccountRepository struct{}

func NewUserAccountRepository() *UserAccountRepository {
	return &UserAccountRepository{}
}

func (r *UserAccountRepository) DelByUserIdFromCache(ctx context.Context, userID int64) error {
	rd, err := xredis.GetPaasRedis()
	if err != nil {
		return err
	}

	key := r.getUserIDKey(userID)

	return rd.Del(ctx, key).Err()
}

// GetByUserID 根据用户ID获取账号信息
func (r *UserAccountRepository) GetByUserID(ctx context.Context, userID int64) (*domain.UserAccount, error) {
	// 先从缓存获取
	rd, err := xredis.GetPaasRedis()
	if err != nil {
		return nil, err
	}

	key := r.getUserIDKey(userID)
	rsBts, err := rd.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if string(rsBts) == notFoundPlaceHolder {
		return nil, constant.ErrRecordNotFound
	}

	if len(rsBts) > 0 {
		account := &domain.UserAccount{}
		if err := json.Unmarshal(rsBts, account); err != nil {
			return nil, err
		}
		return account, nil
	}

	// 缓存未命中，从数据库获取
	db, err := db.GetPaasDB()
	if err != nil {
		return nil, err
	}

	account := &domain.UserAccount{}
	if err := db.Where("user_id = ?", userID).
		Preload("UserProfile").
		Preload("UserPrivateInfo").
		First(account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 缓存空值
			rd.Set(ctx, key, notFoundPlaceHolder, accountCacheExp)
			return nil, constant.ErrRecordNotFound
		}
		return nil, err
	}

	// 缓存
	accountBts, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	rd.Set(ctx, key, string(accountBts), accountCacheExp)
	return account, nil
}

// GetOneByOpenid 根据openid获取账号信息
func (r *UserAccountRepository) GetOneByOpenid(ctx context.Context, openid string) (*domain.UserAccount, error) {
	// 先从缓存获取
	rd, err := xredis.GetPaasRedis()
	if err != nil {
		return nil, err
	}

	key := r.getOpenidKey(openid)
	rsBts, err := rd.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if string(rsBts) == notFoundPlaceHolder {
		return nil, constant.ErrRecordNotFound
	}

	if len(rsBts) > 0 {
		account := &domain.UserAccount{}
		if err := json.Unmarshal(rsBts, account); err != nil {
			return nil, err
		}
		return account, nil
	}

	// 缓存未命中，从数据库获取
	db, err := db.GetPaasDB()
	if err != nil {
		return nil, err
	}

	account := &domain.UserAccount{}
	if err := db.Where("openid = ? and deregistered_at is null", openid).
		Preload("UserProfile").
		Preload("UserPrivateInfo").
		First(account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 缓存空值
			rd.Set(ctx, key, notFoundPlaceHolder, accountCacheExp)
			return nil, constant.ErrRecordNotFound
		}
		return nil, err
	}

	// 缓存
	accountBts, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	rd.Set(ctx, key, string(accountBts), accountCacheExp)
	return account, nil
}

// Save 保存账号信息，同时保存关联的用户资料和私有信息
func (r *UserAccountRepository) Save(ctx context.Context, account *domain.UserAccount,
	profile *domain.UserProfile, privateInfo *domain.UserPrivateInfo,
	beforeSave func(account *domain.UserAccount) bool) error {
	db, err := db.GetPaasDB()
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// 保存账号信息
		if err := tx.Save(account).Error; err != nil {
			return err
		}

		// 调用beforeSave回调
		if beforeSave != nil {
			beforeSave(account)
		}

		if err := tx.Save(account).Error; err != nil {
			return err
		}

		// 保存用户资料
		if profile != nil {
			if err := tx.Save(profile).Error; err != nil {
				return err
			}
		}

		// 保存私有信息
		if privateInfo != nil {
			if err := tx.Save(privateInfo).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteFromCache 从缓存中删除账号信息
func (r *UserAccountRepository) DeleteFromCache(ctx context.Context, openid string) error {
	rd, err := xredis.GetPaasRedis()
	if err != nil {
		return err
	}

	key := r.getOpenidKey(openid)

	return rd.Del(ctx, key).Err()
}

// getOpenidKey 获取openid缓存键
func (r *UserAccountRepository) getOpenidKey(openid string) string {
	return "paas:user:openid:" + openid
}

// getUserIDKey 获取userID缓存键
func (r *UserAccountRepository) getUserIDKey(userID int64) string {
	return "paas:user:id:" + fmt.Sprintf("%d", userID)
}
