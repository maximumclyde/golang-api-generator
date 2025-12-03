package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/iancoleman/strcase"
)

func createMigration(name string, empty bool) {
	fmt.Print("Creating migrations folder... ")

	err := os.MkdirAll(config.Paths.Migrations, os.ModePerm)
	if err != nil {
		fmt.Println("\n❌ Could not create migrations folder")
		panic(err)
	}

	fmt.Println("✅")

	migrationName := name
	// if the migration name is not set, then ask the user for the name
	if migrationName == "" {
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("Input the migration name: ")
		_, err := fmt.Fscan(reader, &migrationName)
		if err != nil {
			fmt.Println("\n❌ Internal error while scanning for input, please retry")
			panic(err)
		}
		if migrationName == "" {
			fmt.Println("\n❌ Invalid migration name")
			panic("invalid_name")
		}
	}

	snakeName := strcase.ToSnake(migrationName)
	sqlVersion := time.Now().UnixMilli() / 1000

	upFileName := fmt.Sprintf("%v_%v.up.sql", sqlVersion, snakeName)
	downFileName := fmt.Sprintf("%v_%v.down.sql", sqlVersion, snakeName)

	if empty {
		// if empty we simply need to create the files and return
		fmt.Print("Creating " + upFileName + "... ")
		upFile, err := os.Create(path.Join(config.Cwd, config.Paths.Migrations, upFileName))
		if err != nil {
			fmt.Println("\n❌ Could not create migration file")
			panic(err)
		}
		defer upFile.Close()

		fmt.Print("Creating " + downFileName + "... ")
		downFile, err := os.Create(path.Join(config.Cwd, config.Paths.Migrations, downFileName))
		if err != nil {
			fmt.Println("\n❌ Could not create migration file")
			panic(err)
		}
		defer downFile.Close()
		fmt.Println("✅")

		return
	}

	upFileToRead := "migrations/template.up.sql"
	downFileToRead := "migrations/template.down.sql"
	if *Custom {
		upFileToRead = "migrations/custom_template.up.sql"
		downFileToRead = "migrations/custom_template.down.sql"
	}

	tRg := regexp.MustCompile(`(?m)template`)

	upData, err := efs.ReadFile(upFileToRead)
	if err != nil {
		fmt.Println("\n❌ Could not load sql template data")
		panic(err)
	}

	upData = tRg.ReplaceAll(upData, ([]byte)(snakeName))

	downData, err := efs.ReadFile(downFileToRead)
	if err != nil {
		fmt.Println("\n❌ Could not load sql template data")
		panic(err)
	}

	downData = tRg.ReplaceAll(downData, ([]byte)(snakeName))

	err = os.WriteFile(path.Join(config.Cwd, config.Paths.Migrations, upFileName), upData, os.ModePerm)
	if err != nil {
		fmt.Println("\n❌ Could not write " + upFileName)
		panic(err)
	}

	err = os.WriteFile(path.Join(config.Cwd, config.Paths.Migrations, downFileName), downData, os.ModePerm)
	if err != nil {
		fmt.Println("\n❌ Could not write " + downFileName)
		panic(err)
	}

	fmt.Println("✅")
}
