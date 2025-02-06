package main

import "github.com/joho/godotenv"

func GetVariablesEnv() (envData map[string]string, err error) {
	return godotenv.Read(VariablesImportFile)
}

func GetLocalDevEnv() (envData map[string]string, err error) {
	return godotenv.Read(LocalDevImportFile)
}

func GetReleaseEnv(providerBrand string) (envData map[string]string, err error) {
	return godotenv.Read(GetImportFileName(providerBrand))
}

func GetImportFileName(providerBrand string) (importFileName string) {
	if providerBrand == "ngg" {
		importFileName = NggImportFile
	} else {
		importFileName = ShorelineImportFile
	}
	return importFileName
}
