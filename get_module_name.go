package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

func getModuleName() string {
	mod, err := os.ReadFile(path.Join(config.Cwd, "./go.mod"))
	if err != nil {
		fmt.Println("\n❌ Could not read go.mod file")
		panic(err)
	}

	rg := regexp.MustCompile(`(?m)module [^\s]+`)
	moduleDef := rg.Find(mod)
	if moduleDef == nil {
		fmt.Println("\n❌ Module is not defined correctly")
		panic("module_not_defined")
	}

	return strings.Replace(string(moduleDef), "module ", "", 1)
}
