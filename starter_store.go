package main

import (
	"api-generator/models"

	"gorm.io/gorm"
)

// #region models
type StoreServices struct {
}

type Store struct {
	Services StoreServices
	TxKey    *models.TxKey
}

func NewStore(db *gorm.DB) *Store {
	//#region new
	store := &Store{
		Services: StoreServices{},
		TxKey:    new(models.TxKey),
	}

	return store
}
