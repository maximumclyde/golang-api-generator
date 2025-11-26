package main

import (
	"path"
	"strings"
)

func replaceImports(data []byte) []byte {
	strData := string(data)

	strData = strings.Replace(strData, HandlersImport, path.Join(config.Module, config.Paths.Handlers), 1)
	strData = strings.Replace(strData, InterfacesImport, path.Join(config.Module, config.Paths.Interfaces), 1)
	strData = strings.Replace(strData, ModelsImport, path.Join(config.Module, config.Paths.Models), 1)
	strData = strings.Replace(strData, RouterImport, path.Join(config.Module, config.Paths.Router), 1)
	strData = strings.Replace(strData, ServicesImport, path.Join(config.Module, config.Paths.Services), 1)
	strData = strings.Replace(strData, StoreImport, path.Join(config.Module, config.Paths.Store), 1)
	strData = strings.Replace(strData, UtilsImport, path.Join(config.Module, config.Paths.Utils), 1)

	return ([]byte)(strData)
}
