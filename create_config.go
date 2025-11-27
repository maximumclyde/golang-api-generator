package main

import (
	"fmt"
	"os"
)

func createConfig() {
	fmt.Printf("Creating %s... ", ConfigPath)

	configData, err := efs.ReadFile("generator.config.json")
	if err != nil {
		fmt.Println("\n❌ Could not read config template data")
		panic(err)
	}

	err = os.WriteFile(ConfigPath, configData, os.ModePerm)
	if err != nil {
		fmt.Println("\n❌ Could not create " + ConfigPath)
		panic(err)
	}

	fmt.Println("✅")
}
