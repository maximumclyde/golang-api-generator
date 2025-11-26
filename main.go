package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"

	schema "github.com/xeipuuv/gojsonschema"
)

// #region embed
//
//go:embed rest.schema.json handlers interfaces migrations models router services store utils
var efs embed.FS

// #region config paths
var ConfigPath = "./rest.config.json"
var SchemaName = "rest.schema.json"

// #region flags
var (
	NeedsHelp        = flag.Bool("help", false, "Display help page")
	ManualConfigPath = flag.String("config", ConfigPath, "Manually sets the config path")
	NoHandler        = flag.Bool("no-handler", false, "Disables automatically creating handlers for the new service")
)

// #region tokens
var (
	HandlerTokens   = `handlers\.`
	InterfaceTokens = `interfaces\.`
	ModelTokens     = `models\.`
	RouterTokens    = `router\.`
	ServiceTokens   = `services\.`
	StoreTokens     = `store\.`
	UtilTokens      = `utils\.`
)

// #region imports
var (
	HandlersImport   = `github.com/maximumclyde/golang-api-generator/handlers`
	InterfacesImport = `github.com/maximumclyde/golang-api-generator/interfaces`
	ModelsImport     = `github.com/maximumclyde/golang-api-generator/models`
	RouterImport     = `github.com/maximumclyde/golang-api-generator/router`
	ServicesImport   = `github.com/maximumclyde/golang-api-generator/services`
	StoreImport      = `github.com/maximumclyde/golang-api-generator/store`
	UtilsImport      = `github.com/maximumclyde/golang-api-generator/utils`
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
		fmt.Println("⚠️  WARNING: rest.config.json could not be found. Preceding with defaults...")

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

	fmt.Println("✅ Done!")
}
