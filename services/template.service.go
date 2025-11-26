package services

import (
	"github.com/maximumclyde/golang-api-generator/models"

	"gorm.io/gorm"
)

type TemplateService struct {
	RestService[models.Template]
}

func NewTemplateService(db *gorm.DB, txK *models.TxKey) *TemplateService {
	//#region new
	rest := NewRestService[models.Template](db, txK)
	service := &TemplateService{
		RestService: *rest,
	}
	return service
}
