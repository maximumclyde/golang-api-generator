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
	//#region create folder
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
		//#region empty name
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

	//#region migrator
	migPath := path.Join(config.Cwd, config.Paths.Migrator)
	_, err = os.Stat(path.Join(migPath, "main.go"))
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Print("Creating " + migPath + "... ")
			err = os.MkdirAll(migPath, os.ModePerm)
			if err != nil {
				fmt.Println("⚠️  WARNING: migrator folders could not be created. You'll need a custom handler to apply the migrations")
				return
			}
			fmt.Println("✅")

			fmt.Print("Creating migrator main.go... ")
			mig, err := efs.ReadFile("migrator/main.go")
			if err != nil {
				fmt.Println("⚠️  WARNING: migrator template data could not be read. You'll need a custom handler to apply the migrations")
				return
			}

			// check for the migration path and replace the migration path
			rg := regexp.MustCompile(`(?m)NewWithDatabaseInstance\(\"`)
			match := rg.FindIndex(mig)
			if match != nil {
				migrationFolder := path.Join(config.Cwd, config.Paths.Migrations)
				res := append([]byte{}, mig[:match[1]]...)
				res = append(res, []byte("file://"+migrationFolder)...)
				res = append(res, mig[match[1]:]...)
				mig = res
			}

			mig = contentReplace(mig, "main")

			err = os.WriteFile(path.Join(migPath, "main.go"), mig, os.ModePerm)
			if err != nil {
				fmt.Println("⚠️  WARNING: migrator file could not be created. You'll need a custom handler to apply the migrations")
				return
			}

			fmt.Println("✅")
			fmt.Println("⚠️  WARNING: run \"go mod tidy\" to automatically install all the required dependencies")

		} else {
			fmt.Println("⚠️  WARNING: migrator file could not be found. You'll need a custom handler to apply the migrations")
		}
	}

	if empty {
		//#region empty migration
		// if empty we simply need to create the files and return
		fmt.Print("Creating " + upFileName + "... ")
		upFile, err := os.Create(path.Join(config.Cwd, config.Paths.Migrations, upFileName))
		if err != nil {
			fmt.Println("\n❌ Could not create migration file")
			panic(err)
		}
		defer upFile.Close()
		fmt.Println("✅")

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

	//#region write migrations
	upFileToRead := "migrations/template.up.sql"
	downFileToRead := "migrations/template.down.sql"
	if *Custom {
		upFileToRead = "migrations/custom_template.up.sql"
		downFileToRead = "migrations/custom_template.down.sql"
	}

	tRg := regexp.MustCompile(`(?m)template`)

	fmt.Print("Creating " + upFileName + "... ")
	upData, err := efs.ReadFile(upFileToRead)
	if err != nil {
		fmt.Println("\n❌ Could not load sql template data")
		panic(err)
	}

	upData = tRg.ReplaceAll(upData, ([]byte)(snakeName))

	err = os.WriteFile(path.Join(config.Cwd, config.Paths.Migrations, upFileName), upData, os.ModePerm)
	if err != nil {
		fmt.Println("\n❌ Could not write " + upFileName)
		panic(err)
	}

	fmt.Println("✅")

	fmt.Print("Creating " + downFileName + "... ")
	downData, err := efs.ReadFile(downFileToRead)
	if err != nil {
		fmt.Println("\n❌ Could not load sql template data")
		panic(err)
	}

	downData = tRg.ReplaceAll(downData, ([]byte)(snakeName))

	err = os.WriteFile(path.Join(config.Cwd, config.Paths.Migrations, downFileName), downData, os.ModePerm)
	if err != nil {
		fmt.Println("\n❌ Could not write " + downFileName)
		panic(err)
	}

	fmt.Println("✅")
}
