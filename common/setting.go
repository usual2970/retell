package common

import (
	"context"
	"encoding/json"

	"github.com/usual2970/retell/internal/repository"
)

func GetSettingByKey[T any](ctx context.Context, key string) (T, error) {
	var rs T
	repo := repository.NewSettingRepository()
	setting, err := repo.GetSetting(ctx, key)
	if err != nil {
		return rs, err
	}

	if setting.Value == "" {
		return rs, nil
	}

	if err := json.Unmarshal([]byte(setting.Value), &rs); err != nil {
		return rs, err
	}
	return rs, nil
}
