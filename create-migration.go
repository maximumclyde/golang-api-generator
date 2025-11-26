package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
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
	fileName := fmt.Sprintf("%v_%v.sql", time.Now().UnixMilli()/1000, snakeName)

	fmt.Print("Creating migration " + fileName + "... ")

	file, err := os.Create(path.Join(config.Cwd, config.Paths.Migrations, fileName))
	if err != nil {
		fmt.Println("\n❌ Could not create migration file")
		panic(err)
	}
	defer file.Close()

	if empty {
		fmt.Println("✅")
		return
	}

	templateData, err := efs.ReadFile("migrations/template.sql")
	if err != nil {
		fmt.Println("\n❌ Could not load sql template data")
		panic(err)
	}

	_, err = file.WriteString(
		strings.ReplaceAll(string(templateData), "template", snakeName),
	)
	if err != nil {
		fmt.Println("\n❌ Could not write sql migration file")
		panic(err)
	}

	fmt.Println("✅")
}
