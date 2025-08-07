package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/domain/constant"
	"github.com/usual2970/retell/pkg/db"
	xredis "github.com/usual2970/retell/pkg/redis"
	"gorm.io/gorm"
)

var (
	resourcePositionCacheExp = time.Hour * 24
	resourceCacheExp         = time.Hour * 24
)

type ResourcePositionRepository struct{}

func NewResourcePositionRepository() *ResourcePositionRepository {
	return &ResourcePositionRepository{}
}

func (r *ResourcePositionRepository) Update(
	ctx context.Context,
	position *domain.ResourcePosition,
	toAddResources, toUpdateResources, toDeleteResources []domain.Resource,
) error {
	db, err := db.GetPaasDB()
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// 更新资源位置
		if err := tx.Save(position).Error; err != nil {
			return err
		}

		// 添加新资源
		for _, resource := range toAddResources {
			resource.PositionID = position.ID
			if err := tx.Save(&resource).Error; err != nil {
				return err
			}
			r.deleteResourceCache(ctx, resource.URI)
		}

		// 更新已有资源
		for _, resource := range toUpdateResources {
			if err := tx.Save(&resource).Error; err != nil {
				return err
			}
			r.deleteResourceCache(ctx, resource.URI)
		}

		// 删除指定资源
		for _, resource := range toDeleteResources {
			if err := tx.Delete(&resource).Error; err != nil {
				return err
			}
			r.deleteResourceCache(ctx, resource.URI)
		}

		r.deletePositionCache(ctx, position.URI)
		r.deleteGamePositionCache(ctx, position.RelationURI)

		return nil
	})
}

// Save 保存资源位置及其资源，支持事务
func (r *ResourcePositionRepository) Save(
	ctx context.Context,
	position *domain.ResourcePosition,
	resources []domain.Resource,
	fn func(position *domain.ResourcePosition),
) error {
	db, err := db.GetPaasDB()
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// 保存资源位置
		if err := tx.Save(position).Error; err != nil {
			return err
		}

		// 调用回调函数，可以在保存后对 position 进行操作
		if fn != nil {
			fn(position)
		}

		// 保存资源列表
		if len(resources) > 0 {
			for _, resource := range resources {
				if err := tx.Save(&resource).Error; err != nil {
					return err
				}
			}
		}

		// 删除缓存
		r.deletePositionCache(ctx, position.URI)

		return nil
	})
}

// GetByURI 根据URI获取资源位置
func (r *ResourcePositionRepository) GetByURI(ctx context.Context, uri string) (*domain.ResourcePosition, error) {
	// 先从缓存获取
	rd, err := xredis.GetPaasRedis()
	if err != nil {
		return nil, err
	}

	key := r.getPositionCacheKey(uri)
	rsBts, err := rd.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if string(rsBts) == notFoundPlaceHolder {
		return nil, constant.ErrRecordNotFound
	}

	if len(rsBts) > 0 {
		position := &domain.ResourcePosition{}
		if err := json.Unmarshal(rsBts, position); err != nil {
			return nil, err
		}
		return position, nil
	}

	// 缓存未命中，从数据库获取
	db, err := db.GetPaasDB()
	if err != nil {
		return nil, err
	}

	position := &domain.ResourcePosition{}
	if err := db.Where("uri = ?", uri).
		Preload("Resources", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC").Where("is_show = ?", 1)

		}).
		First(position).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 缓存空值
			rd.Set(ctx, key, notFoundPlaceHolder, resourcePositionCacheExp)
			return nil, constant.ErrRecordNotFound
		}
		return nil, err
	}

	// 缓存结果
	positionBts, err := json.Marshal(position)
	if err != nil {
		return nil, err
	}

	rd.Set(ctx, key, string(positionBts), resourcePositionCacheExp)
	return position, nil
}

// GetResourceByURI 根据URI获取资源
func (r *ResourcePositionRepository) GetResourceByURI(ctx context.Context, uri string) (*domain.Resource, error) {
	// 先从缓存获取
	rd, err := xredis.GetPaasRedis()
	if err != nil {
		return nil, err
	}

	key := r.getResourceCacheKey(uri)
	rsBts, err := rd.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if string(rsBts) == notFoundPlaceHolder {
		return nil, constant.ErrRecordNotFound
	}

	if len(rsBts) > 0 {
		resource := &domain.Resource{}
		if err := json.Unmarshal(rsBts, resource); err != nil {
			return nil, err
		}
		return resource, nil
	}

	// 缓存未命中，从数据库获取
	db, err := db.GetPaasDB()
	if err != nil {
		return nil, err
	}

	resource := &domain.Resource{}
	if err := db.Where("uri = ?", uri).First(resource).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 缓存空值
			rd.Set(ctx, key, notFoundPlaceHolder, resourceCacheExp)
			return nil, constant.ErrRecordNotFound
		}
		return nil, err
	}

	// 缓存结果
	resourceBts, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}

	rd.Set(ctx, key, string(resourceBts), resourceCacheExp)
	return resource, nil
}

// SaveResource 保存单个资源
func (r *ResourcePositionRepository) SaveResource(ctx context.Context, resource *domain.Resource) error {
	db, err := db.GetPaasDB()
	if err != nil {
		return err
	}

	if err := db.Save(resource).Error; err != nil {
		return err
	}

	// 删除资源缓存
	r.deleteResourceCache(ctx, resource.URI)

	// 同时需要删除对应位置的缓存，因为位置预加载了资源列表
	position := &domain.ResourcePosition{}
	if err := db.First(position, resource.PositionID).Error; err == nil {
		r.deletePositionCache(ctx, position.URI)
		r.deleteGamePositionCache(ctx, position.RelationURI)
	}

	return nil
}

// DeleteResource 删除资源
func (r *ResourcePositionRepository) DeleteResource(ctx context.Context, resource *domain.Resource) error {
	db, err := db.GetPaasDB()
	if err != nil {
		return err
	}

	// 先查询对应的位置，以便后续删除缓存
	position := &domain.ResourcePosition{}
	if err := db.First(position, resource.PositionID).Error; err != nil {
		return err
	}

	// 删除资源
	if err := db.Delete(resource).Error; err != nil {
		return err
	}

	// 删除资源缓存
	r.deleteResourceCache(ctx, resource.URI)

	// 删除位置缓存
	r.deletePositionCache(ctx, position.URI)

	r.deleteGamePositionCache(ctx, position.RelationURI)

	return nil
}

// GetByGameURI 根据游戏URI获取资源位置
func (r *ResourcePositionRepository) GetByGameURI(ctx context.Context, gameURI string) (*domain.ResourcePosition, error) {
	// 先从缓存获取
	rd, err := xredis.GetPaasRedis()
	if err != nil {
		return nil, err
	}

	// 为游戏URI关联的资源位创建特殊缓存键
	key := r.getGamePositionCacheKey(gameURI)
	rsBts, err := rd.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	// 处理缓存命中
	if string(rsBts) == notFoundPlaceHolder {
		return nil, constant.ErrRecordNotFound
	}

	if len(rsBts) > 0 {
		position := &domain.ResourcePosition{}
		if err := json.Unmarshal(rsBts, position); err != nil {
			return nil, err
		}
		return position, nil
	}

	// 缓存未命中，从数据库获取
	db, err := db.GetPaasDB()
	if err != nil {
		return nil, err
	}

	position := &domain.ResourcePosition{}
	if err := db.Where("relation_uri = ?", gameURI).
		Preload("Resources").
		First(position).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 缓存空值
			rd.Set(ctx, key, notFoundPlaceHolder, resourcePositionCacheExp)
			return nil, constant.ErrRecordNotFound
		}
		return nil, err
	}

	// 缓存结果
	positionBts, err := json.Marshal(position)
	if err != nil {
		return nil, err
	}

	rd.Set(ctx, key, string(positionBts), resourcePositionCacheExp)

	// 同时更新普通URI缓存，以保持一致性
	normalKey := r.getPositionCacheKey(position.URI)
	rd.Set(ctx, normalKey, string(positionBts), resourcePositionCacheExp)

	return position, nil
}

// 获取游戏URI关联的资源位置缓存键
func (r *ResourcePositionRepository) getGamePositionCacheKey(gameURI string) string {
	return fmt.Sprintf("paas:game_resource_position:%s", gameURI)
}

func (r *ResourcePositionRepository) deleteGamePositionCache(ctx context.Context, gameURI string) error {
	rd, err := xredis.GetPaasRedis()
	if err != nil {
		return err
	}

	key := r.getGamePositionCacheKey(gameURI)
	return rd.Del(ctx, key).Err()
}

// 删除资源位置缓存
func (r *ResourcePositionRepository) deletePositionCache(ctx context.Context, uri string) error {
	rd, err := xredis.GetPaasRedis()
	if err != nil {
		return err
	}

	key := r.getPositionCacheKey(uri)
	return rd.Del(ctx, key).Err()
}

// 删除资源缓存
func (r *ResourcePositionRepository) deleteResourceCache(ctx context.Context, uri string) error {
	rd, err := xredis.GetPaasRedis()
	if err != nil {
		return err
	}

	key := r.getResourceCacheKey(uri)
	return rd.Del(ctx, key).Err()
}

// 获取资源位置缓存键
func (r *ResourcePositionRepository) getPositionCacheKey(uri string) string {
	return fmt.Sprintf("paas:resource_position:%s", uri)
}

// 获取资源缓存键
func (r *ResourcePositionRepository) getResourceCacheKey(uri string) string {
	return fmt.Sprintf("paas:resource:%s", uri)
}
