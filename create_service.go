package main

import (
	"bufio"
	"errors"
	"fmt"
	"go/format"
	"os"
	"path"
	"regexp"
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
	filePrefix := strings.ToLower(strcase.ToSnake(resourceName))

	kebabResource := strcase.ToKebab(resourceName)

	rscRg := regexp.MustCompile(`(?m)Template`)
	cRscRg := regexp.MustCompile(`(?m)CustomTemplate`)
	tmptsRg := regexp.MustCompile(`(?m)templates`)
	templateRg := regexp.MustCompile(`(?m)template`)

	var tmpModelData []byte

	//#region migration
	createMigration(resourceName, false)

	//#region model
	fmt.Print("Creating model... ")
	mdlFile := filePrefix + ".model.go"
	efsFileName := "models/template.model.go"
	if *Custom {
		efsFileName = "models/custom_template.model.go"
	}
	tmpModelData, err = efs.ReadFile(efsFileName)
	if err != nil {
		fmt.Println("\n❌ Could not load model template data")
		panic(err)
	}

	tmpModelData = contentReplace(tmpModelData, getPackageName(config.Paths.Models))
	if *Custom {
		tmpModelData = cRscRg.ReplaceAll(tmpModelData, ([]byte)(resourceName))
	} else {
		tmpModelData = rscRg.ReplaceAll(tmpModelData, ([]byte)(resourceName))
	}
	tmpModelData = tmptsRg.ReplaceAll(tmpModelData, ([]byte)(tableName))

	err = os.WriteFile(path.Join(config.Paths.Models, mdlFile), tmpModelData, os.ModePerm)
	if err != nil {
		fmt.Println("\n❌ Could not create models file")
		panic(err)
	}

	fmt.Println("✅")

	//#region service
	fmt.Print("Creating service... ")
	svcFile := filePrefix + ".service.go"
	efsFileName = "services/template.service.go"
	if *Custom {
		efsFileName = "services/custom_template.service.go"
	}
	tmpModelData, err = efs.ReadFile(efsFileName)
	if err != nil {
		fmt.Println("\n❌ Could not load service template data")
		panic(err)
	}

	tmpModelData = contentReplace(tmpModelData, getPackageName(config.Paths.Services))
	if *Custom {
		tmpModelData = cRscRg.ReplaceAll(tmpModelData, ([]byte)(resourceName))
	} else {
		tmpModelData = rscRg.ReplaceAll(tmpModelData, ([]byte)(resourceName))
	}

	err = os.WriteFile(path.Join(config.Paths.Services, svcFile), tmpModelData, os.ModePerm)
	if err != nil {
		fmt.Println("\n❌ Could not create services file")
		panic(err)
	}

	fmt.Println("✅")

	//#region store
	storeFile := "store.go"
	storePath := path.Join(config.Paths.Store, storeFile)

	_, err = os.Stat(storePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Print("Creating store... ")
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
			tmpModelData = rscRg.ReplaceAll(tmpModelData, ([]byte)(resourceName))

			err = os.WriteFile(storePath, tmpModelData, os.ModePerm)
			if err != nil {
				fmt.Println("\n❌ Could not create store file")
				panic(err)
			}
		} else {
			fmt.Println("\n❌ Could not stat store file at " + storePath)
			panic(err)
		}
	} else {
		// we need to update the store with the new service
		fmt.Print("Updating store... ")
		tmpModelData, err = os.ReadFile(storePath)
		if err != nil {
			fmt.Println("\n❌ Could not read store file")
			panic(err)
		}

		tmpModelData = ([]byte)(strings.ReplaceAll(string(tmpModelData), "StoreServices struct {", "StoreServices struct {\n"+resourceName+" *services."+resourceName+"Service "))

		// checks if the services where imported or not
		// if the services were not imported it means that
		// this is the first service that's being created after init
		servicesPath := path.Join(config.Module, config.Paths.Services)
		impRg := regexp.MustCompile(`(?m)` + servicesPath)

		im := impRg.Find(tmpModelData)
		if im == nil {
			// the import was not found, so we need to add the import and the txK declaration
			tmpModelData = ([]byte)(strings.ReplaceAll(string(tmpModelData), "Services: StoreServices{", "Services: StoreServices{\n"+resourceName+": services.New"+resourceName+"Service(db, txk),\n"))
			tmpModelData = ([]byte)(strings.Replace(string(tmpModelData), "import (", "import ("+"\n"+"\""+servicesPath+"\" ", 1))
			tmpModelData = ([]byte)(strings.Replace(string(tmpModelData), "store := &Store", "txk := new(models.TxKey)\n\n"+"store := &Store", 1))
		} else {
			tmpModelData = ([]byte)(strings.ReplaceAll(string(tmpModelData), "Services: StoreServices{", "Services: StoreServices{\n"+resourceName+": services.New"+resourceName+"Service(db, txk), "))
		}

		tmpModelData = contentReplace(tmpModelData, getPackageName(config.Paths.Store))

		err = os.WriteFile(storePath, tmpModelData, os.ModePerm)
		if err != nil {
			fmt.Println("\n❌ Could not update store file")
			panic(err)
		}
	}

	fmt.Println("✅")

	//#region seeder
	seederFile := "main.go"
	seederPath := path.Join(config.Paths.Seeders, seederFile)

	_, err = os.Stat(seederPath)
	if err != nil {
		if os.IsNotExist(err) {
			// create the seeder file
			fmt.Print("Creating seeder... ")
			err = os.MkdirAll(path.Join(config.Cwd, config.Paths.Seeders), os.ModePerm)
			if err != nil {
				fmt.Println("\n❌ Could not create seeder folder")
				panic(err)
			}

			tmpModelData, err = efs.ReadFile("seeders/main.go")
			if err != nil {
				fmt.Println("\n❌ Could not load seeder template data")
				panic(err)
			}

			tmpModelData = contentReplace(tmpModelData, "main")
			tmpModelData = templateRg.ReplaceAll(tmpModelData, ([]byte)(lowerResourceName))
			tmpModelData = rscRg.ReplaceAll(tmpModelData, ([]byte)(resourceName))

			err = os.WriteFile(seederPath, tmpModelData, os.ModePerm)
			if err != nil {
				fmt.Println("\n❌ Could not create seeder file")
				panic(err)
			}
		} else {
			fmt.Println("\n❌ Could not stat seeder file at " + seederPath)
			panic(err)
		}
	} else {
		// update the seeder file
		fmt.Print("Updating seeder... ")

		tmpModelData, err = efs.ReadFile("seeders/main.go")
		if err != nil {
			fmt.Println("\n❌ Could not load seeder template data")
			panic(err)
		}

		regionRg := regexp.MustCompile(`(?m)\/\/\#region Template`)
		regionIndex := regionRg.FindIndex(tmpModelData)
		if regionIndex == nil {
			fmt.Println("\n❌ Could not find seeder template region")
			panic(errors.New("invalid_template"))
		}

		tmpModelData = tmpModelData[regionIndex[0] : len(tmpModelData)-3]
		tmpModelData = templateRg.ReplaceAll(tmpModelData, ([]byte)(lowerResourceName))
		tmpModelData = rscRg.ReplaceAll(tmpModelData, ([]byte)(resourceName))

		actualSeederData, err := os.ReadFile(seederPath)
		if err != nil {
			fmt.Println("\n❌ Could not read seeder file")
			panic(err)
		}

		actualSeederData = append(actualSeederData[:len(actualSeederData)-3], ([]byte)("\n\n")...)
		actualSeederData = append(actualSeederData, tmpModelData...)
		actualSeederData = append(actualSeederData, ([]byte)("\n}\n")...)
		actualSeederData, err = format.Source(actualSeederData)
		if err != nil {
			fmt.Println("\n❌ Could not format seeder file")
			panic(err)
		}

		err = os.WriteFile(seederPath, actualSeederData, os.ModePerm)
		if err != nil {
			fmt.Println("\n❌ Could not update seeder file")
			panic(err)
		}
	}

	fmt.Println("✅")

	//#region template seeder
	fmt.Print("Creating " + resourceName + " seeder... ")
	sdFile := filePrefix + ".seeder.go"
	tmpModelData, err = efs.ReadFile("seeders/template.seeder.go")
	if err != nil {
		fmt.Println("\n❌ Could not load seeder template data")
		panic(err)
	}

	tmpModelData = contentReplace(tmpModelData, "main")
	tmpModelData = rscRg.ReplaceAll(tmpModelData, ([]byte)(resourceName))

	err = os.WriteFile(path.Join(config.Paths.Seeders, sdFile), tmpModelData, os.ModePerm)
	if err != nil {
		fmt.Println("\n❌ Could not create seeder file")
		panic(err)
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

	handlerEfsFilename := "handlers/template.handler.go"
	if *Custom {
		handlerEfsFilename = "handlers/custom_template.handler.go"
	}

	tmpModelData, err = efs.ReadFile(handlerEfsFilename)
	if err != nil {
		fmt.Println("\n❌ Could not load handler template data")
		panic(err)
	}

	tmpModelData = contentReplace(tmpModelData, getPackageName(config.Paths.Handlers))
	if *Custom {
		tmpModelData = cRscRg.ReplaceAll(tmpModelData, ([]byte)(resourceName))
	} else {
		tmpModelData = rscRg.ReplaceAll(tmpModelData, ([]byte)(resourceName))
	}
	tmpModelData = tmptsRg.ReplaceAll(tmpModelData, []byte(kebabResource))

	err = os.WriteFile(path.Join(config.Paths.Handlers, handlerFile), tmpModelData, os.ModePerm)
	if err != nil {
		fmt.Println("\n❌ Could not create handler file")
		panic(err)
	}

	fmt.Println("✅")

	//#region router
	routerFile := "router.go"
	routerPath := path.Join(config.Paths.Router, routerFile)

	_, err = os.Stat(routerPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Print("Creating router... ")
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

			tmpModelData = rscRg.ReplaceAll(tmpModelData, ([]byte)(resourceName))
			tmpModelData = templateRg.ReplaceAll(tmpModelData, ([]byte)(lowerResourceName))

			err = os.WriteFile(routerPath, tmpModelData, os.ModePerm)
			if err != nil {
				fmt.Println("\n❌ Could not create router file")
				panic(err)
			}
		} else {
			fmt.Println("\n❌ Could not stat router file at " + routerPath)
			panic(err)
		}
	} else {
		// we need to update the store with the new handler
		fmt.Print("Updating router... ")
		tmpModelData, err = os.ReadFile(routerPath)
		if err != nil {
			fmt.Println("\n❌ Could not read router file")
			panic(err)
		}

		// check whether the handlers where imported
		// if the handlers are not imported, this means that
		// this is the first service after init
		handlersPath := path.Join(config.Module, config.Paths.Handlers)
		impRg := regexp.MustCompile(`(?m)` + handlersPath)
		im := impRg.Find(tmpModelData)
		if im == nil {
			// the import was not found so we add the import and the
			// public and protected routes declarations
			impRg = regexp.MustCompile(`(?m)import \(`)
			tmpModelData = impRg.ReplaceAll(tmpModelData, ([]byte)("import (\n"+"\""+handlersPath+"\""))

			groupsBytes := ([]byte)("publicRoutes := g.Group(\"\")\n" + "protectedRoutes := g.Group(\"\")\n\n")
			impRg = regexp.MustCompile(`(?m)return.*`)
			tmpModelData = impRg.ReplaceAllFunc(tmpModelData, func(b []byte) []byte {
				return append(groupsBytes, b...)
			})
		}

		replaceStr := "\n" + "//#region " + resourceName + "\n" + lowerResourceName + "Handler := handlers.New" + resourceName + "Handler(s)\n"
		replaceStr = replaceStr + lowerResourceName + "Handler.RegisterRoutes(publicRoutes, protectedRoutes)\n"

		rtRg := regexp.MustCompile(`(?m)return g`)
		tmpModelData = rtRg.ReplaceAll(tmpModelData, ([]byte)(replaceStr+"\nreturn g"))
		tmpModelData = contentReplace(tmpModelData, getPackageName(config.Paths.Router))

		err = os.WriteFile(routerPath, tmpModelData, os.ModePerm)
		if err != nil {
			fmt.Println("\n❌ Could not update router file")
			panic(err)
		}
	}

	fmt.Println("✅")
}
