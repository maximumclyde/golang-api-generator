package main

func replaceTokens(data []byte) []byte {
	handlerImportToken := ([]byte)(getPackageName(config.Paths.Handlers) + ".")
	interfaceImportToken := ([]byte)(getPackageName(config.Paths.Interfaces) + ".")
	modelImportToken := ([]byte)(getPackageName(config.Paths.Models) + ".")
	routerImportToken := ([]byte)(getPackageName(config.Paths.Router) + ".")
	serviceImportToken := ([]byte)(getPackageName(config.Paths.Services) + ".")
	storeImportToken := ([]byte)(getPackageName(config.Paths.Store) + ".")
	utilImportToken := ([]byte)(getPackageName(config.Paths.Utils) + ".")

	data = HandlerTokens.ReplaceAll(data, handlerImportToken)
	data = InterfaceTokens.ReplaceAll(data, interfaceImportToken)
	data = ModelTokens.ReplaceAll(data, modelImportToken)
	data = RouterTokens.ReplaceAll(data, routerImportToken)
	data = ServiceTokens.ReplaceAll(data, serviceImportToken)
	data = StoreTokens.ReplaceAll(data, storeImportToken)
	data = UtilTokens.ReplaceAll(data, utilImportToken)

	return data
}
