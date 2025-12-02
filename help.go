package main

import (
	"fmt"
)

func help() {
	fmt.Println("Rest Service Tool")
	fmt.Println("A code generator for rest APIs using gorm, gin and go-faker")
	fmt.Printf("\n")

	fmt.Println("Usage:")
	fmt.Println("api-generator [--help]")
	fmt.Println("api-generator [...options] <command>")
	fmt.Printf("\n")

	fmt.Println("Commands:")
	fmt.Println("init			Initialize structure on an empty folder")
	fmt.Println("create")
	fmt.Println(" - service		Creates a service along with it's definitions ans migrations")
	fmt.Println(" - migration		Creates a new empty migration file")
	fmt.Println(" - config		Creates a default configuration file")
	fmt.Printf("\n")

	fmt.Println("Options:")
	fmt.Println("--help			Display this help page")
	fmt.Println("--config		Set the config folder")
	fmt.Println("--no-handler		Disables automatically creating handlers for the new service")
	fmt.Println("--custom		Indicates that the new service should not be derived from the default predefined provider")
}
