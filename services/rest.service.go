package services

import (
	"context"

	"github.com/maximumclyde/golang-api-generator/interfaces"
	"github.com/maximumclyde/golang-api-generator/models"
	"github.com/maximumclyde/golang-api-generator/utils"

	"gorm.io/gorm"
)

type RestService[T interfaces.IGormModel, C any, Q any, P any] struct {
	db     *gorm.DB
	config *models.RestProviderConfig
}

func NewRestService[T interfaces.IGormModel, C any, Q any, P any](db *gorm.DB, txK *models.TxKey) *RestService[T, C, Q, P] {
	mdl := new(T)

	provider := &RestService[T, C, Q, P]{
		db: db.Model(mdl),
		config: &models.RestProviderConfig{
			TxKey: txK,
			Table: (*mdl).TableName(),
		},
	}

	return provider
}

func (s *RestService[T, C, Q, P]) Create(ctx context.Context, data *C) error {
	//#region create
	db := s.GetDb(ctx)
	return db.Create(data).Error
}

func (s *RestService[T, C, Q, P]) GetById(ctx context.Context, id string) (*T, error) {
	//#region get by id
	db := s.GetDb(ctx)
	data := new(T)
	err := db.Where("id = ?", id).First(data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *RestService[T, C, Q, P]) Find(ctx context.Context, query *Q) (*[]T, error) {
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

func (s *RestService[T, C, Q, P]) Patch(ctx context.Context, id string, data *P) error {
	//#region patch
	db := s.GetDb(ctx)

	return db.Where("id = ?", id).Updates(data).Error
}

func (s *RestService[T, C, Q, P]) Update(ctx context.Context, data T, query *Q) error {
	//#region update
	db := s.GetDb(ctx)

	dbWConditions := utils.AttachQueryConditions(db, query)

	return dbWConditions.UpdateColumns(data).Error
}

func (s *RestService[T, C, Q, P]) Remove(ctx context.Context, id string) error {
	//#region remove
	db := s.GetDb(ctx)
	mdl := new(T)
	return db.Delete(mdl, "id = ?", id).Error
}

func (s *RestService[T, C, Q, P]) Delete(ctx context.Context, query *Q) error {
	//#region delete
	db := s.GetDb(ctx)

	dbWConditions := utils.AttachQueryConditions(db, query)

	return dbWConditions.Delete(ctx).Error
}

func (s *RestService[T, C, Q, P]) GetConfig() *models.RestProviderConfig {
	//#region get config
	return s.config
}

func (s *RestService[T, C, Q, P]) GetDb(ctx context.Context) *gorm.DB {
	//#region get db
	txK := s.GetConfig().TxKey
	var db *gorm.DB = s.db

	// in case we are in a transaction, we put the provider in the context and work the same way
	txDb := ctx.Value(txK)
	if txDb != nil {
		mdl := new(T)
		db = txDb.(*gorm.DB).Model(mdl)
	}

	return db.WithContext(ctx)
}
