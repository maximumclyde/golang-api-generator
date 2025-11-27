package main

import (
	"fmt"
	"go/format"
)

func contentReplace(data []byte, packageName string) []byte {
	data = replacePackageName(data, packageName)
	data = replaceImports(data)
	data = replaceTokens(data)

	data, err := format.Source(data)
	if err != nil {
		fmt.Println("\n❌ Could not format file")
		panic(err)
	}

	return data
}
