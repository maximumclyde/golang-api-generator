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

func (s *RestService[T]) Create(ctx context.Context, data *T) error {
	//#region create
	db := s.GetDb(ctx)
	return db.Create(data).Error
}

func (s *RestService[T]) GetById(ctx context.Context, id string) (*T, error) {
	//#region get by id
	db := s.GetDb(ctx)
	val, err := db.Get(id)
	if err {
		return nil, errors.New("not_found")
	}
	return val.(*T), nil
}

func (s *RestService[T]) Find(ctx context.Context, query ...any) (*[]T, error) {
	//#region find
	db := s.GetDb(ctx)

	data := new([]T)

	dbWConditions := utils.AttachQueryConditions(db, query)

	err := dbWConditions.Find(data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (s *RestService[T]) Patch(ctx context.Context, id string, data *T) error {
	//#region patch
	db := s.GetDb(ctx)
	return db.Where("id = ?", id).Save(data).Error
}

func (s *RestService[T]) Update(ctx context.Context, data T, query ...any) error {
	//#region update
	db := s.GetDb(ctx)

	dbWConditions := utils.AttachQueryConditions(db, query)

	return dbWConditions.UpdateColumns(data).Error
}

func (s *RestService[T]) Remove(ctx context.Context, id string) error {
	//#region remove
	db := s.GetDb(ctx)
	return db.Where("id = ?", id).Delete(ctx).Error
}

func (s *RestService[T]) Delete(ctx context.Context, query ...any) error {
	//#region delete
	db := s.GetDb(ctx)

	dbWConditions := utils.AttachQueryConditions(db, query)

	return dbWConditions.Delete(ctx).Error
}

func (s *RestService[T]) GetConfig() *models.RestProviderConfig {
	//#region get config
	return s.config
}

func (s *RestService[T]) GetDb(ctx context.Context) *gorm.DB {
	//#region get db
	txK := s.GetConfig().TxKey
	var db *gorm.DB = s.db
	txDb := ctx.Value(txK)
	if txDb != nil {
		db = txDb.(*gorm.DB)
	}

	return db.WithContext(ctx)
}
