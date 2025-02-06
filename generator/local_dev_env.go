package main

import (
	"log"
	"os"
	"text/template"
)

var (
	devLogPrefix          = "[Generator Local Dev]"
	LocalDevEnvFilePath   = "./env/dev/local.env"
	LocalDevTempalateFile = "./env/dev/local.tmpl"
)

// sets up the local dev environment variables
func GenerateLocalDevEnv(providerBrand string) {
	envData, err := getMergedEnvData(providerBrand)
	if err != nil {
		log.Println(devLogPrefix+" failed to get env data: ", err)
		return
	}

	tmpl, err := template.ParseFiles(LocalDevTempalateFile)
	if err != nil {
		log.Println(devLogPrefix+" failed to parse files: ", err)
		return
	}

	outputFile, err := os.Create(LocalDevEnvFilePath)
	if err != nil {
		log.Println(devLogPrefix+" failed to create output file: ", err)
		return
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, envData)
	if err != nil {
		log.Println(devLogPrefix+" failed to execute template: ", err)
		return
	}
}

func getMergedEnvData(providerBrand string) (envData map[string]string, err error) {
	releaseEnv, err := GetReleaseEnv(providerBrand)
	if err != nil {
		return nil, err
	}
	devEnv, err := GetVariablesEnv()
	if err != nil {
		return nil, err
	}
	return mergeEnvData(releaseEnv, devEnv), nil
}

func mergeEnvData(envData map[string]string, priorityEnvData map[string]string) map[string]string {
	for key, value := range priorityEnvData {
		envData[key] = value
	}

	return envData
}
