package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/iancoleman/strcase"
)

func createService() {
	var tableName string

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Input the db table name for the model: ")
	_, err := fmt.Fscan(reader, &tableName)
	if err != nil {
		fmt.Println("\n❌ Internal error while scanning for input, please retry")
		panic(err)
	}

	// resource name in capitalized camel case
	lowerResourceName := strcase.ToLowerCamel(tableName)
	resourceName := strings.ToUpper(string([]byte(lowerResourceName)[0])) + string([]byte(lowerResourceName)[1:])

	tableName = strcase.ToSnake(tableName)

	// file prefix to attach to the start of the file
	filePrefix := strings.ToLower(strcase.ToKebab(resourceName))

	//#region migration
	createMigration(resourceName, false)

	//#region creating model
	fmt.Print("Creating model... ")
	mdlFile := filePrefix + ".model.go"
	tmpModelData, err := efs.ReadFile("models/template.model.go")
	if err != nil {
		fmt.Println("\n❌ Could not load model template data")
		panic(err)
	}

	tmpModelData = contentReplace(tmpModelData, getPackageName(config.Paths.Models))
	tmpModelData = ([]byte)(strings.ReplaceAll((string)(tmpModelData), "Template", resourceName))
	tmpModelData = ([]byte)(strings.ReplaceAll((string)(tmpModelData), "templates", tableName))

	err = os.WriteFile(path.Join(config.Paths.Models, mdlFile), tmpModelData, os.ModePerm)
	if err != nil {
		fmt.Println("\n❌ Could not create models file")
		panic(err)
	}

	fmt.Println("✅")

	//#region service
	fmt.Print("Creating service... ")
	svcFile := filePrefix + ".service.go"
	tmpModelData, err = efs.ReadFile("services/template.service.go")
	if err != nil {
		fmt.Println("\n❌ Could not load service template data")
		panic(err)
	}

	tmpModelData = contentReplace(tmpModelData, getPackageName(config.Paths.Services))
	tmpModelData = ([]byte)(strings.ReplaceAll((string)(tmpModelData), "Template", resourceName))

	err = os.WriteFile(path.Join(config.Paths.Services, svcFile), tmpModelData, os.ModePerm)
	if err != nil {
		fmt.Println("\n❌ Could not create services file")
		panic(err)
	}

	fmt.Println("✅")

	//#region store
	fmt.Print("Creating/updating store... ")
	storeFile := "store.go"
	storePath := path.Join(config.Paths.Store, storeFile)

	_, err = os.Stat(storePath)
	if err != nil {
		if os.IsNotExist(err) {
			// the file simply does not exist so we create it
			err = os.MkdirAll(config.Paths.Store, os.ModePerm)
			if err != nil {
				fmt.Println("\n❌ Could not create store folder")
				panic(err)
			}

			tmpModelData, err = efs.ReadFile("store/store.go")
			if err != nil {
				fmt.Println("\n❌ Could not load store template data")
				panic(err)
			}

			tmpModelData = contentReplace(tmpModelData, getPackageName(config.Paths.Store))
			tmpModelData = ([]byte)(strings.ReplaceAll((string)(tmpModelData), "Template", resourceName))

			err = os.WriteFile(storePath, tmpModelData, os.ModePerm)
			if err != nil {
				fmt.Println("\n❌ Could not create store file")
				panic(err)
			}
		} else {
			fmt.Println("\n❌ Could not open store file")
			panic(err)
		}
	} else {
		// we need to update the store with the new service
		tmpModelData, err = os.ReadFile(storePath)
		if err != nil {
			fmt.Println("\n❌ Could not read store file")
			panic(err)
		}

		tmpModelData = ([]byte)(strings.ReplaceAll(string(tmpModelData), "StoreServices struct {", "StoreServices struct {\n"+resourceName+" *services."+resourceName+"Service "))
		tmpModelData = ([]byte)(strings.ReplaceAll(string(tmpModelData), "Services: StoreServices{", "Services: StoreServices{\n"+resourceName+": services.New"+resourceName+"Service(db, txk), "))
		tmpModelData = contentReplace(tmpModelData, getPackageName(config.Paths.Store))

		err = os.WriteFile(storePath, tmpModelData, os.ModePerm)
		if err != nil {
			fmt.Println("\n❌ Could not update store file")
			panic(err)
		}
	}

	fmt.Println("✅")

	//#region handler
	if *NoHandler {
		return
	}

	fmt.Print("Creating handler... ")
	handlerFile := filePrefix + ".handler.go"
	err = os.MkdirAll(path.Join(config.Cwd, config.Paths.Handlers), os.ModePerm)
	if err != nil {
		fmt.Println("\n❌ Could not create handler folder")
		panic(err)
	}

	tmpModelData, err = efs.ReadFile("handlers/template.handler.go")
	if err != nil {
		fmt.Println("\n❌ Could not load handler template data")
		panic(err)
	}

	tmpModelData = contentReplace(tmpModelData, getPackageName(config.Paths.Handlers))
	tmpModelData = ([]byte)(strings.ReplaceAll((string)(tmpModelData), "Template", resourceName))
	tmpModelData = ([]byte)(strings.ReplaceAll((string)(tmpModelData), "templates", filePrefix))

	err = os.WriteFile(path.Join(config.Paths.Handlers, handlerFile), tmpModelData, os.ModePerm)
	if err != nil {
		fmt.Println("\n❌ Could not create handler file")
		panic(err)
	}

	fmt.Println("✅")

	//#region router
	fmt.Print("Creating/updating router... ")
	routerFile := "router.go"
	routerPath := path.Join(config.Paths.Router, routerFile)

	_, err = os.Stat(routerPath)
	if err != nil {
		if os.IsNotExist(err) {
			// the file simply does not exist so we create it
			err = os.MkdirAll(config.Paths.Router, os.ModePerm)
			if err != nil {
				fmt.Println("\n❌ Could not create router folder")
				panic(err)
			}

			tmpModelData, err = efs.ReadFile("router/router.go")
			if err != nil {
				fmt.Println("\n❌ Could not load router template data")
				panic(err)
			}

			tmpModelData = contentReplace(tmpModelData, getPackageName(config.Paths.Router))
			tmpModelData = ([]byte)(strings.ReplaceAll((string)(tmpModelData), "templateHandler", lowerResourceName+"Handler"))
			tmpModelData = ([]byte)(strings.ReplaceAll((string)(tmpModelData), "Template", resourceName))

			err = os.WriteFile(routerPath, tmpModelData, os.ModePerm)
			if err != nil {
				fmt.Println("\n❌ Could not create router file")
				panic(err)
			}
		} else {
			fmt.Println("\n❌ Could not open router file")
			panic(err)
		}
	} else {
		// we need to update the store with the new handler
		tmpModelData, err = os.ReadFile(routerPath)
		if err != nil {
			fmt.Println("\n❌ Could not read router file")
			panic(err)
		}

		replaceStr := "\n" + lowerResourceName + "Handler := handlers.New" + resourceName + "Handler(s)\n"
		replaceStr = replaceStr + lowerResourceName + "Handler.RegisterRoutes(publicRoutes, protectedRoutes)\n"

		tmpModelData = ([]byte)(strings.ReplaceAll(string(tmpModelData), "return g", replaceStr+"\nreturn g"))
		tmpModelData = contentReplace(tmpModelData, getPackageName(config.Paths.Router))

		err = os.WriteFile(routerPath, tmpModelData, os.ModePerm)
		if err != nil {
			fmt.Println("\n❌ Could not update router file")
			panic(err)
		}
	}

	fmt.Println("✅")
}
