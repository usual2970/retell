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

type MessageRepository struct{}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{}
}

func (m *MessageRepository) getTaskByUriKey(uri string) string {
	return "paas:message:task:uri:" + uri
}

func (m *MessageRepository) GetTaskByUri(ctx context.Context, uri string) (*domain.MessageTask, error) {
	// 先从缓存中获取
	rd, err := xredis.GetPaasRedis()
	if err != nil {
		return nil, err
	}

	key := m.getTaskByUriKey(uri)

	rsBts, err := rd.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if string(rsBts) == notFoundPlaceHolder {
		return nil, constant.ErrRecordNotFound
	}

	if len(rsBts) > 0 {
		task := &domain.MessageTask{}
		if err := json.Unmarshal(rsBts, task); err != nil {
			return nil, err
		}
		return task, nil
	}

	// 从数据库中获取
	db, err := db.GetPaasDB()
	if err != nil {
		return nil, err
	}

	task := &domain.MessageTask{}
	if err := db.Preload("MessaageTpls").Where("uri = ?", uri).First(task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 缓存空值
			rd.Set(ctx, key, notFoundPlaceHolder, time.Hour*24)
			return nil, constant.ErrRecordNotFound
		}
		return nil, err
	}

	// 缓存结果
	taskBts, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	rd.Set(ctx, key, string(taskBts), time.Hour*24)

	return task, nil
}
func (m *MessageRepository) BatchInsertHistory(ctx context.Context, historys []*domain.MessageHistory) error {

	db, err := db.GetPaasDB()
	if err != nil {
		return err
	}

	if err := db.Create(historys).Error; err != nil {
		return err
	}

	return nil
}

func (m *MessageRepository) Total(ctx context.Context, filters map[string]any) (int64, error) {
	db, err := db.GetPaasDB()
	if err != nil {
		return 0, err
	}

	var total int64
	query := db.Model(&domain.MessageHistory{})

	sql, vals, err := utils.WhereBuild(filters)
	if err != nil {
		return 0, err
	}

	if err := query.Where(sql, vals...).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (m *MessageRepository) List(ctx context.Context, filters map[string]any, pagination *constant.Pagination, orderBy string) ([]domain.MessageHistory, error) {
	db, err := db.GetPaasDB()
	if err != nil {
		return nil, err
	}

	query := db.Model(&domain.MessageHistory{})

	sql, vals, err := utils.WhereBuild(filters)
	if err != nil {
		return nil, err
	}

	query = query.Where(sql, vals...)

	if pagination != nil {
		query = query.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}

	if orderBy != "" {
		query = query.Order(orderBy)
	}

	var historys []domain.MessageHistory
	if err := query.Find(&historys).Error; err != nil {
		return nil, err
	}

	return historys, nil
}

func (m *MessageRepository) Updates(ctx context.Context, where map[string]any, updates map[string]any) error {
	db, err := db.GetPaasDB()
	if err != nil {
		return err
	}

	sql, vals, err := utils.WhereBuild(where)
	if err != nil {
		return err
	}

	query := db.Model(&domain.MessageHistory{}).Where(sql, vals...)

	if err := query.Updates(updates).Error; err != nil {
		return err
	}

	return nil
}
