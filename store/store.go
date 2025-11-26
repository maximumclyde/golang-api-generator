package store

import (
	"github.com/maximumclyde/golang-api-generator/models"
	"github.com/maximumclyde/golang-api-generator/services"

	"gorm.io/gorm"
)

// #region models
type StoreServices struct {
	Template *services.TemplateService
}

type Store struct {
	Services StoreServices
	TxKey    *models.TxKey
}

func NewStore(db *gorm.DB) *Store {
	//#region new
	txk := new(models.TxKey)

	store := &Store{
		Services: StoreServices{
			Template: services.NewTemplateService(db, txk),
		},
		TxKey: new(models.TxKey),
	}

	return store
}
