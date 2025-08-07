package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/domain/constant"
	"github.com/usual2970/retell/pkg/db"
	xredis "github.com/usual2970/retell/pkg/redis"
	"github.com/usual2970/retell/pkg/utils"
	"gorm.io/gorm"
)

const categoryCacheExp = 1 * time.Hour // 24小时缓存

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (c *CategoryRepository) GetByURI(ctx context.Context, uri string) (*domain.Category, error) {
	rd, err := xredis.GetPaasRedis()
	if err != nil {
		return nil, err
	}

	key := c.getByURICacheKey(uri)

	rsBts, err := rd.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if string(rsBts) == notFoundPlaceHolder {
		return nil, constant.ErrRecordNotFound
	}

	if len(rsBts) > 0 {
		category := &domain.Category{}
		if err := json.Unmarshal(rsBts, category); err != nil {
			return nil, err
		}
	}
	db, err := db.GetPaasDB()
	if err != nil {
		return nil, err
	}
	category := &domain.Category{}
	if err := db.WithContext(ctx).Where("uri = ?", uri).First(category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 缓存空值
			rd.Set(ctx, key, notFoundPlaceHolder, categoryCacheExp)
			return nil, constant.ErrRecordNotFound
		}
		return nil, err
	}
	categoryBts, err := json.Marshal(category)
	if err != nil {
		return nil, err
	}
	rd.Set(ctx, key, string(categoryBts), categoryCacheExp)
	return category, nil

}
func (c *CategoryRepository) List(ctx context.Context, where map[string]any, pagination *constant.Pagination, order string) ([]domain.Category, error) {
	db, err := db.GetPaasDB()
	if err != nil {
		return nil, err
	}

	var categories []domain.Category
	query := db.WithContext(ctx).Model(&domain.Category{})

	sql, vals, err := utils.WhereBuild(where)
	if err != nil {
		return nil, err
	}

	query = query.Where(sql, vals...)

	if pagination != nil {
		query = query.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}

	if err := query.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *CategoryRepository) getByURICacheKey(uri string) string {
	return "paas:category:" + uri
}
