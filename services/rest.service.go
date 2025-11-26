package services

import (
	"context"
	"errors"

	"github.com/maximumclyde/golang-api-generator/interfaces"
	"github.com/maximumclyde/golang-api-generator/models"
	"github.com/maximumclyde/golang-api-generator/utils"

	"gorm.io/gorm"
)

type RestService[T interfaces.IGormModel] struct {
	db     *gorm.DB
	config *models.RestProviderConfig
}

func NewRestService[T interfaces.IGormModel](db *gorm.DB, txK *models.TxKey) *RestService[T] {
	mdl := new(T)

	provider := &RestService[T]{
		db: db.Model(mdl),
		config: &models.RestProviderConfig{
			TxKey: txK,
			Table: (*mdl).TableName(),
		},
	}

	return provider
}

func (p *RestService[T]) Create(ctx context.Context, data *T) error {
	//#region create
	db := p.GetDb(ctx)
	return db.Create(data).Error
}

func (p *RestService[T]) GetById(ctx context.Context, id string) (*T, error) {
	//#region get by id
	db := p.GetDb(ctx)
	val, err := db.Get(id)
	if err {
		return nil, errors.New("not_found")
	}
	return val.(*T), nil
}

func (p *RestService[T]) Find(ctx context.Context, query ...any) ([]T, error) {
	//#region find
	db := p.GetDb(ctx)

	data := []T{}

	dbWConditions := utils.AttachQueryConditions(db, query)

	err := dbWConditions.Find(data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (p *RestService[T]) Patch(ctx context.Context, id string, data *T) error {
	//#region patch
	db := p.GetDb(ctx)
	return db.Where("id = ?", id).Save(data).Error
}

func (p *RestService[T]) Update(ctx context.Context, data T, query ...any) error {
	//#region update
	db := p.GetDb(ctx)

	dbWConditions := utils.AttachQueryConditions(db, query)

	return dbWConditions.UpdateColumns(data).Error
}

func (p *RestService[T]) Remove(ctx context.Context, id string) error {
	//#region remove
	db := p.GetDb(ctx)
	return db.Where("id = ?", id).Delete(ctx).Error
}

func (p *RestService[T]) Delete(ctx context.Context, query ...any) error {
	//#region delete
	db := p.GetDb(ctx)

	dbWConditions := utils.AttachQueryConditions(db, query)

	return dbWConditions.Delete(ctx).Error
}

func (p *RestService[T]) GetConfig() *models.RestProviderConfig {
	//#region get config
	return p.config
}

func (p *RestService[T]) GetDb(ctx context.Context) *gorm.DB {
	//#region get db
	txK := p.GetConfig().TxKey
	var db *gorm.DB = ctx.Value(txK).(*gorm.DB)
	if db == nil {
		db = p.db
	}

	return db.WithContext(ctx)
}
