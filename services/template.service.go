package services

import (
	"api-generator/models"

	"gorm.io/gorm"
)

type TemplateService struct {
	RestService[
		models.Template,
		models.TemplateCreate,
		models.TemplateQuery,
		models.TemplatePatch,
	]
}

func NewTemplateService(db *gorm.DB, txK *models.TxKey) *TemplateService {
	//#region new
	rest := NewRestService[
		models.Template,
		models.TemplateCreate,
		models.TemplateQuery,
		models.TemplatePatch,
	](db, txK)
	service := &TemplateService{
		RestService: *rest,
	}
	return service
}
