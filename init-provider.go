package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
)

func initProvider() {
	//#region create directories
	pathsToCreate := []string{
		config.Paths.Interfaces,
		config.Paths.Models,
		config.Paths.Services,
		config.Paths.Utils,
	}

	var err error = nil

	for _, p := range pathsToCreate {
		fmt.Print("Creating " + p + "... ")

		newPath := path.Join(config.Cwd, p)
		err = os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			fmt.Println("\n❌ Could not create " + newPath)
			panic(err)
		}
		fmt.Println("✅")

		efsName := "interfaces"
		switch p {
		case config.Paths.Models:
			efsName = "models"
		case config.Paths.Services:
			efsName = "services"
		case config.Paths.Utils:
			efsName = "utils"
		}

		//#region create files
		_ = fs.WalkDir(efs, efsName, func(efsPath string, d fs.DirEntry, initErr error) error {
			if efsPath == "." || efsPath == efsName {
				return nil
			}

			fileName := getPackageName(efsPath)
			packageName := getPackageName(newPath)

			filePath := path.Join(newPath, fileName)

			if strings.Contains(fileName, "template") {
				return nil
			}

			fmt.Print("Creating " + filePath + "... ")

			fileData, err := efs.ReadFile(efsPath)
			if err != nil {
				fmt.Println("\n❌ Could not load data for " + fileName)
				panic(err)
			}

			fileData = contentReplace(fileData, packageName)

			err = os.WriteFile(filePath, fileData, os.ModePerm)
			if err != nil {
				fmt.Println("\n❌ Could not create " + filePath)
				panic(err)
			}

			fmt.Println("✅")

			return nil
		})
	}
}
