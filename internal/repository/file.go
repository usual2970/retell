package repository

import (
	"context"

	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/pkg/db"
)

type FileRepository struct{}

func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

func (r *FileRepository) Save(ctx context.Context, file *domain.File) error {

	db, err := db.GetPaasDB()
	if err != nil {
		return err
	}

	return db.Create(file).Error
}
