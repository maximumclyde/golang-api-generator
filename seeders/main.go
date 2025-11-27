package main

import (
	"gorm.io/gorm"
)

func main() {
	//#region main
	var db *gorm.DB = nil
	var err error = nil

	//#region Template
	templateSeeder := &TemplateSeeder{}
	_, err = templateSeeder.Seed(db)
	if err != nil {
		rbe := templateSeeder.Rollback(db)
		if rbe != nil {
			panic(rbe)
		}
		panic(err)
	}
}
