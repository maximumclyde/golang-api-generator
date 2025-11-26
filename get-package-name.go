package main

import "strings"

func getPackageName(path string) string {
	split := strings.Split(path, "/")
	return split[len(split)-1]
}
