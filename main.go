package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"

	schema "github.com/xeipuuv/gojsonschema"
)

// #region embed
//
//go:embed rest.schema.json generator.config.json handlers interfaces migrations models router services store utils seeders db server migrator
var efs embed.FS

// #region started files fs
//
//go:embed starter_store.go starter_router.go
var starterFs embed.FS

// #region config paths
var ConfigPath = "./generator.config.json"
var SchemaName = "rest.schema.json"

// #region flags
var (
	NeedsHelp        = flag.Bool("help", false, "Display help page")
	ManualConfigPath = flag.String("config", ConfigPath, "Manually sets the config path")
	NoHandler        = flag.Bool("no-handler", false, "Disables automatically creating handlers for the new service")
	Custom           = flag.Bool("custom", false, "Indicates that the new service should not be derived from the default predefined provider")
)

// #region tokens
var (
	HandlerTokens   = regexp.MustCompile(`(?m)handlers\.`)
	InterfaceTokens = regexp.MustCompile(`(?m)interfaces\.`)
	ModelTokens     = regexp.MustCompile(`(?m)models\.`)
	RouterTokens    = regexp.MustCompile(`(?m)router\.`)
	ServiceTokens   = regexp.MustCompile(`(?m)services\.`)
	StoreTokens     = regexp.MustCompile(`(?m)store\.`)
	UtilTokens      = regexp.MustCompile(`(?m)utils\.`)
)

// #region imports
var (
	HandlersImport   = `api-generator/handlers`
	InterfacesImport = `api-generator/interfaces`
	ModelsImport     = `api-generator/models`
	RouterImport     = `api-generator/router`
	ServicesImport   = `api-generator/services`
	StoreImport      = `api-generator/store`
	UtilsImport      = `api-generator/utils`
	DatabaseImport   = `api-generator/db`
)

// #region paths
type Paths struct {
	Handlers   string `json:"handlers"`
	Interfaces string `json:"interfaces"`
	Migrations string `json:"migrations"`
	Models     string `json:"models"`
	Router     string `json:"router"`
	Services   string `json:"services"`
	Store      string `json:"store"`
	Utils      string `json:"utils"`
	Seeders    string `json:"seeders"`
	Database   string `json:"database"`
	Server     string `json:"server"`
	Migrator   string `json:"migrator"`
}

// #region config
type Config struct {
	Cwd    string
	Module string `json:"module"`
	Paths  Paths  `json:"paths"`
}

var config Config = Config{
	Paths: Paths{
		Handlers:   "./internal/handlers",
		Interfaces: "./internal/interfaces",
		Migrations: "./migrations",
		Models:     "./internal/models",
		Router:     "./internal/router",
		Services:   "./internal/services",
		Store:      "./internal/store",
		Utils:      "./internal/utils",
		Seeders:    "./cmd/seeders",
		Database:   "./internal/db",
		Server:     "./cmd/server",
		Migrator:   "./cmd/migrator",
	},
}

func init() {
	//#region init
	flag.Parse()

	var (
		data []byte = []byte{}
		err  error  = nil
	)

	cwd, err := os.Getwd()
	if err != nil {
		panic("❌ Could not load cwd")
	}

	config.Cwd = cwd

	moduleName := getModuleName()
	config.Module = moduleName

	if *ManualConfigPath != "" && *ManualConfigPath != ConfigPath {
		ConfigPath = *ManualConfigPath
	}

	data, err = os.ReadFile(path.Join(cwd, ConfigPath))
	if err != nil {
		if *NeedsHelp || len(flag.Args()) == 0 || (flag.Arg(0) == "create" && flag.Arg(1) == "config") {
			return
		}

		fmt.Println("⚠️  WARNING: generator.config.json could not be found. Preceding with defaults...")
		return
	}

	schemaData, err := efs.ReadFile(SchemaName)
	if err != nil {
		panic("❌ Could not read schema definition")
	}

	referenceLoader := schema.NewBytesLoader(schemaData)
	configLoader := schema.NewBytesLoader(data)

	result, err := schema.Validate(referenceLoader, configLoader)
	if err != nil {
		panic("❌ Config is not valid")
	}

	if !result.Valid() {
		fmt.Println(result.Errors())
		panic("❌ Config is not valid")
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		panic("❌ Could not load config")
	}
}

func main() {
	//#region main
	command := flag.Arg(0)

	if command == "" || *NeedsHelp {
		help()
		return
	}

	switch command {
	case "init":
		initProvider()
	case "create":
		switch flag.Arg(1) {
		case "config":
			createConfig()
		case "service":
			createService()
		case "migration":
			createMigration("", true)
		default:
			fmt.Printf("⚠️  WARNING: Invalid use of command \"create\":\n\n")
			help()
			return
		}
	default:
		fmt.Printf("⚠️  WARNING: Invalid command \"%s\":\n\n", command)
		help()
		return
	}

	fmt.Println("Done! ✅")
}
