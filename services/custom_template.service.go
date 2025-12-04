package services

import (
	"context"

	"api-generator/models"

	"gorm.io/gorm"
)

type CustomTemplateService struct {
	db     *gorm.DB
	config *models.RestProviderConfig
}

func NewCustomTemplateService(db *gorm.DB, txK *models.TxKey) *CustomTemplateService {
	//#region new
	mdl := &models.CustomTemplate{}

	service := &CustomTemplateService{
		db: db,
		config: &models.RestProviderConfig{
			Table: mdl.TableName(),
			TxKey: txK,
		},
	}

	return service
}

func (s *CustomTemplateService) GetConfig() *models.RestProviderConfig {
	//#region get config
	return s.config
}

func (s *CustomTemplateService) GetDb(ctx context.Context) *gorm.DB {
	//#region get db
	txK := s.GetConfig().TxKey
	var db *gorm.DB = s.db
	txDb := ctx.Value(txK)
	if txDb != nil {
		db = txDb.(*gorm.DB).Model(&models.CustomTemplate{})
	}

	return db.WithContext(ctx)
}
