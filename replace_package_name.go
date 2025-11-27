package main

import "regexp"

func replacePackageName(data []byte, packageName string) []byte {
	rg := regexp.MustCompile(`(?m)(package \w+)`)
	data = rg.ReplaceAllFunc(data, func(b []byte) []byte {
		return []byte("package " + packageName)
	})

	return data
}
