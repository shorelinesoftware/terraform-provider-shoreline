package main

import (
	"log"
	"os"
	"text/template"
)

var (
	confLogPrefix         = "[Generator Provider Conf]"
	LocalConfTemplatePath = "./provider/conf/provider_conf_template.json"
	LocalConfPath         = "./provider/conf/provider_conf.json"
)

func GenerateProviderConf(providerBrand string, envData map[string]string) {

	tmpl, err := template.ParseFiles(LocalConfTemplatePath)
	if err != nil {
		log.Println(confLogPrefix+" failed to parse files: ", err)
		return
	}

	outputFile, err := os.Create(LocalConfPath)
	if err != nil {
		log.Println(confLogPrefix+" failed to create output file: ", err)
		return
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, envData)
	if err != nil {
		log.Println(confLogPrefix+" failed to execute template: ", err)
		return
	}
}
