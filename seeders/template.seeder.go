package main

import (
	"fmt"

	"github.com/go-faker/faker/v4"
	"github.com/maximumclyde/golang-api-generator/models"
	"gorm.io/gorm"
)

// #region type
type TemplateSeeder struct{}

func (t *TemplateSeeder) Seed(db *gorm.DB) ([]*models.Template, error) {
	//#region seed
	fmt.Print("Seeding Template... ")

	allData := make([]*models.Template, 50)
	for i := range 100 {
		data := &models.Template{}
		err := faker.FakeData(data)
		if err != nil {
			fmt.Println("\n❌ Fake data generator error")
			fmt.Println(err)
			return allData, err
		}
		allData[i] = data
	}

	err := db.Create(allData).Error
	if err != nil {
		fmt.Println("\n❌ Database creation error")
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("✅")

	return allData, nil
}

func (t *TemplateSeeder) Rollback(db *gorm.DB) error {
	//#region rollback
	fmt.Print("Rolling back Template... ")
	mdl := &models.Template{}
	err := db.Delete(mdl).Error
	if err != nil {
		fmt.Println("\n❌ Could not roll back Template data")
		fmt.Println(err)
	} else {
		fmt.Println("✅")
	}

	return err
}
