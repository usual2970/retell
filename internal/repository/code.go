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
)

type CodeRepository struct{}

func NewCodeRepository() *CodeRepository {
	return &CodeRepository{}
}

func (r *CodeRepository) getReceiverKey(receiver string, purpose domain.CodePurpose) string {
	return fmt.Sprintf("paas:code:%s:%d", receiver, purpose)
}

func (r *CodeRepository) GetByReceiver(ctx context.Context, receiver string, purpose domain.CodePurpose) (*domain.Code, error) {
	rd, err := xredis.GetPaasRedis()
	if err != nil {
		return nil, err
	}

	key := r.getReceiverKey(receiver, purpose)

	rsBts, err := rd.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, constant.ErrRecordNotFound
		}
		return nil, err
	}

	code := &domain.Code{}
	if err := json.Unmarshal(rsBts, code); err != nil {
		return nil, err
	}

	return code, nil
}
func (r *CodeRepository) Save(ctx context.Context, code *domain.Code) error {
	// 保存到数据库
	db, err := db.GetPaasDB()
	if err != nil {
		return err
	}

	if err := db.Create(code).Error; err != nil {
		return err
	}

	r.saveToCache(ctx, code)

	return nil

}
func (r *CodeRepository) Update(ctx context.Context, code *domain.Code) error {
	// 更新数据库
	db, err := db.GetPaasDB()
	if err != nil {
		return err
	}

	if err := db.Save(code).Error; err != nil {
		return err
	}

	r.saveToCache(ctx, code)

	return nil
}
func (r *CodeRepository) GetByReceiverAndCode(ctx context.Context, receiver, code string, purpose domain.CodePurpose) (*domain.Code, error) {
	rs, err := r.GetByReceiver(ctx, receiver, purpose)
	if err != nil {
		return nil, err
	}

	if rs.Code != code {
		return nil, constant.ErrRecordNotFound
	}

	return rs, nil
}

func (r *CodeRepository) saveToCache(ctx context.Context, code *domain.Code) error {
	rd, err := xredis.GetPaasRedis()
	if err != nil {
		return err
	}

	key := r.getReceiverKey(code.Receiver, code.Purpose)

	codeBts, err := json.Marshal(code)
	if err != nil {
		return err
	}

	return rd.Set(ctx, key, string(codeBts), code.ExpiredAt.Sub(time.Now().UTC())).Err()
}
