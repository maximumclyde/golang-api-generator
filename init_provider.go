package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"
)

func initProvider() {
	var err error = nil

	//#region init questions
	reader1 := bufio.NewReader(os.Stdin)
	reader2 := bufio.NewReader(os.Stdin)
	hostname := ""
	connectionStr := ""
	fmt.Printf("Input the hostname (optional): ")
	fmt.Fscanf(reader1, "%s", &hostname)

	fmt.Printf("Input the db connection string (optional): ")
	fmt.Fscanf(reader2, "%s", &connectionStr)

	//#region starter store
	newStoreFolder := path.Join(config.Cwd, config.Paths.Store)
	fmt.Print("Creating " + newStoreFolder + "... ")
	err = os.MkdirAll(newStoreFolder, os.ModePerm)
	if err != nil {
		fmt.Println("\n❌ Could not create " + newStoreFolder)
		panic(err)
	}
	fmt.Println("✅")

	fmt.Print("Creating store.go... ")
	starterStoreData, err := starterFs.ReadFile("starter_store.go")
	if err != nil {
		fmt.Println("\n❌ Could not read store template data")
		panic(err)
	}
	starterStoreData = contentReplace(starterStoreData, getPackageName(newStoreFolder))
	err = os.WriteFile(path.Join(newStoreFolder, "store.go"), starterStoreData, os.ModePerm)
	if err != nil {
		fmt.Println("\n❌ Could not create store.go")
		panic(err)
	}
	fmt.Println("✅")

	//#region starter router
	newRouterFolder := path.Join(config.Cwd, config.Paths.Router)
	fmt.Print("Creating " + newRouterFolder + "... ")
	err = os.MkdirAll(newRouterFolder, os.ModePerm)
	if err != nil {
		fmt.Println("\n❌ Could not create " + newRouterFolder)
		panic(err)
	}
	fmt.Println("✅")

	fmt.Print("Creating router.go... ")
	starterRouterData, err := starterFs.ReadFile("starter_router.go")
	if err != nil {
		fmt.Println("\n❌ Could not read router template data")
		panic(err)
	}
	starterRouterData = contentReplace(starterRouterData, getPackageName(newRouterFolder))
	err = os.WriteFile(path.Join(newRouterFolder, "router.go"), starterRouterData, os.ModePerm)
	if err != nil {
		fmt.Println("\n❌ Could not create router.go")
		panic(err)
	}
	fmt.Println("✅")

	//#region create directories
	pathsToCreate := []string{
		config.Paths.Interfaces,
		config.Paths.Models,
		config.Paths.Services,
		config.Paths.Utils,
		config.Paths.Database,
		config.Paths.Server,
	}

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
		case config.Paths.Database:
			efsName = "db"
		case config.Paths.Server:
			efsName = "server"
		}

		//#region create files
		_ = fs.WalkDir(efs, efsName, func(efsPath string, d fs.DirEntry, initErr error) error {
			if efsPath == "." || efsPath == efsName {
				return nil
			}

			fileName := getPackageName(efsPath)
			packageName := getPackageName(newPath)
			if efsName == "server" {
				packageName = "main"
			}

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

			// replace the hostname in the appropriate places
			if efsName == "server" && fileName == "main.go" {
				rg := regexp.MustCompile(`(?m)Addr:\s*\"`)
				match := rg.FindIndex(fileData)
				if match != nil {
					res := []byte{}
					res = append(res, fileData[:match[1]]...)
					res = append(res, ([]byte)(hostname)...)
					res = append(res, fileData[match[1]:]...)
					fileData = res
				}
			}

			// replace the connection string in the appropriate places
			if efsName == "db" && fileName == "open_db.go" {
				rg := regexp.MustCompile(`(?m)DSN\:\s*\"`)
				match := rg.FindIndex(fileData)
				if match != nil {
					res := []byte{}
					res = append(res, fileData[:match[1]]...)
					res = append(res, ([]byte)(connectionStr)...)
					res = append(res, fileData[match[1]:]...)
					fileData = res
				}
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

	fmt.Println("⚠️  WARNING: run \"go mod tidy\" to automatically install all the required dependencies")
}
