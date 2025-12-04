package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbConfig struct {
	DSN             string
	MaxConns        int
	MinConns        int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	ConnectTimeout  time.Duration
	AutoMigrate     bool
}

var defaultConfig = &DbConfig{
	DSN:             "",
	MaxConns:        25,
	MinConns:        5,
	MaxConnLifetime: 5 * time.Minute,
	MaxConnIdleTime: 30 * time.Second,
	ConnectTimeout:  10 * time.Second,
	AutoMigrate:     true,
}

func OpenDb(config *DbConfig) (*gorm.DB, error) {
	//#region open
	dbConf := defaultConfig
	if config != nil {
		dbConf = config
	}

	db, err := gorm.Open(postgres.Open(dbConf.DSN))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(dbConf.MaxConns)
	sqlDB.SetMaxIdleConns(dbConf.MinConns)
	sqlDB.SetConnMaxLifetime(dbConf.MaxConnLifetime)
	sqlDB.SetConnMaxIdleTime(dbConf.MaxConnIdleTime)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), dbConf.ConnectTimeout)
	defer cancel()

	if err = sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("GORM database connection established successfully - Max: %d, Min: %d", dbConf.MaxConns, dbConf.MinConns)

	return db, nil
}
