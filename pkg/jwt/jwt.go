package jwt

import (
	"context"

	"github.com/usual2970/retell/domain/constant"

	"strconv"

	"github.com/golang-jwt/jwt/v4"
)

type authKey struct{}

func WithContext(ctx context.Context, token *jwt.Token) context.Context {
	return context.WithValue(ctx, authKey{}, token)
}

func FromContext(ctx context.Context) (*jwt.Token, bool) {
	rs, ok := ctx.Value(authKey{}).(*jwt.Token)
	return rs, ok
}

func GetUserID(ctx context.Context) (int64, error) {
	token, ok := FromContext(ctx)
	if !ok {
		return 0, constant.ErrNotLogin
	}

	mapClaim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, constant.ErrNotLogin
	}

	id, ok := mapClaim["jti"]
	if !ok {
		return 0, constant.ErrNotLogin
	}

	uid, err := strconv.Atoi(id.(string))
	if err != nil {
		return 0, err
	}
	return int64(uid), nil
}

func GetAccessToken(ctx context.Context) (string, error) {
	token, ok := FromContext(ctx)
	if !ok {
		return "", constant.ErrNotLogin
	}

	return token.Raw, nil
}

func GetDeveloprID(ctx context.Context) (int64, error) {
	// 开发阶段先返回一个固定的用户ID
	return 1, nil
}
