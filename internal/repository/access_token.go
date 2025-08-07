package repository

import (
	"context"
	"time"

	"github.com/usual2970/retell/pkg/redis"
)

const AuthExpireDuration = time.Hour * 24 * 30

type AccessTokenRepository struct{}

func NewAccessTokenRepository() *AccessTokenRepository {
	return &AccessTokenRepository{}
}

func (r *AccessTokenRepository) GetAccessToken(ctx context.Context, accessToken string) (int64, error) {
	xredis, err := redis.GetPaasRedis()
	if err != nil {
		return 0, err
	}

	key := r.getAccessTokenKey(accessToken)
	userID, err := xredis.Get(ctx, key).Int64()
	return userID, err
}
func (r *AccessTokenRepository) SetAccessToken(ctx context.Context, accessToken string, userID int64) error {
	xredis, err := redis.GetPaasRedis()
	if err != nil {
		return err
	}

	key := r.getAccessTokenKey(accessToken)
	_, err = xredis.Set(ctx, key, userID, AuthExpireDuration).Result()
	return err
}

func (r *AccessTokenRepository) DelAccessToken(ctx context.Context, accessToken string) error {
	xredis, err := redis.GetPaasRedis()
	if err != nil {
		return err
	}

	key := r.getAccessTokenKey(accessToken)
	_, err = xredis.Del(ctx, key).Result()
	return err
}

func (r *AccessTokenRepository) getAccessTokenKey(accessToken string) string {
	return "paas:access_token:" + accessToken
}
