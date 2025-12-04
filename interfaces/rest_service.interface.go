package interfaces

import (
	"context"

	"api-generator/models"

	"gorm.io/gorm"
)

type IRestService[T IGormModel] interface {
	Create(ctx context.Context, data *T) error
	GetById(ctx context.Context, id string) (*T, error)
	Find(ctx context.Context, query ...any) ([]T, error)
	Patch(ctx context.Context, id string, data *T) error
	Update(ctx context.Context, data T, query ...any) error
	Remove(ctx context.Context, id string) error
	Delete(ctx context.Context, query ...any) error
	GetConfig() *models.RestProviderConfig
	GetDb(ctx context.Context) *gorm.DB
}
