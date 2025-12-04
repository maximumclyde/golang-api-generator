package main

import (
	db "api-generator/db"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	//#region main
	database, err := db.OpenDb(nil)
	if err != nil {
		panic(err)
	}

	sql, err := database.DB()
	if err != nil {
		db.Close(database)
		log.Fatalf("Could not get db instance:\n%+v\n", err)
	}

	driver, err := postgres.WithInstance(sql, &postgres.Config{})
	if err != nil {
		db.Close(database)
		log.Fatalf("Could not get sql driver:\n%+v\n", err)
	}

	m, err := migrate.NewWithDatabaseInstance("", "postgres", driver)
	if err != nil {
		db.Close(database)
		log.Fatalf("Could not create migration:\n%+v\n", err)
	}

	// uncomment to rollback to a version forcefully
	// err = m.Force(...)
	// if err != nil {
	// 	db.Close(database)
	// 	log.Fatalf("Force error:\n%+v\n", err)
	// }

	err = m.Up()
	if err != nil {
		db.Close(database)
		log.Fatalf("Migration error:\n%+v\n", err)
	}

	db.Close(database)
}
