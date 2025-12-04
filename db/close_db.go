package db

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func CloseGorm(db *gorm.DB) error {

	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	err = sqlDB.Close()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	log.Println("GORM database connection closed")
	return nil
}

func Close(db *gorm.DB) {
	if db != nil {
		if err := CloseGorm(db); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
		log.Println("Database manager closed")
	}
}
