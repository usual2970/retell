package secret

import (
	"context"
	"sync"

	"github.com/usual2970/retell/internal/domain"
	"github.com/usual2970/retell/internal/util/app"
)

var once sync.Once
var instance domain.ISecretRepository

type repository struct{}

func NewRepository() domain.ISecretRepository {
	once.Do(func() {
		instance = &repository{}
	})
	return instance
}

func (r *repository) Get(ctx context.Context, filter string) (*domain.Secret, error) {
	record, err := app.Get().FindFirstRecordByFilter("secrets",
		filter,
	)
	if err != nil {
		return nil, err
	}

	ext := make(map[string]string)

	record.UnmarshalJSONField("ext", &ext)
	meta := domain.Meta{
		Id:      record.Id,
		Created: record.GetDateTime("created").Time(),
		Updated: record.GetDateTime("updated").Time(),
	}
	rs := &domain.Secret{

		Uri:         record.GetString("uri"),
		ApiKey:      record.GetString("api_key"),
		SecretKey:   record.GetString("secret_key"),
		Description: record.GetString("description"),
		Ext:         ext,

		Meta: meta,
	}

	return rs, nil
}
