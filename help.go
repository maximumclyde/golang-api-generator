package main

import (
	"fmt"
)

func help() {
	fmt.Println("Rest Service Tool")
	fmt.Println("This is a tool used to create models, migrations and a ready to go rest API using gorm, go migrate and gin")
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
	// fmt.Printf("\n")
}
