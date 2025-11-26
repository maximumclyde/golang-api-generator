package main

import "strings"

func replaceTokens(data []byte) []byte {
	handlerImportToken := getPackageName(config.Paths.Handlers)
	interfaceImportToken := getPackageName(config.Paths.Interfaces)
	modelImportToken := getPackageName(config.Paths.Models)
	routerImportToken := getPackageName(config.Paths.Router)
	serviceImportToken := getPackageName(config.Paths.Services)
	storeImportToken := getPackageName(config.Paths.Store)
	utilImportToken := getPackageName(config.Paths.Utils)

	str := string(data)
	str = strings.ReplaceAll(str, HandlerTokens, handlerImportToken)
	str = strings.ReplaceAll(str, InterfaceTokens, interfaceImportToken)
	str = strings.ReplaceAll(str, ModelTokens, modelImportToken)
	str = strings.ReplaceAll(str, RouterTokens, routerImportToken)
	str = strings.ReplaceAll(str, ServiceTokens, serviceImportToken)
	str = strings.ReplaceAll(str, StoreTokens, storeImportToken)
	str = strings.ReplaceAll(str, UtilTokens, utilImportToken)

	return ([]byte)(str)
}
