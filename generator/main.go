// Copyright 2021, Shoreline Software Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"os"
)

var (
	NggImportFile       string = "./env/release/ngg.env"
	ShorelineImportFile string = "./env/release/shoreline.env"
	LocalDevImportFile  string = "./env/dev/local.env"
	VariablesImportFile string = "./variables.env"
	mainLogPrefix       string = "[ MAIN ]"
)

func main() {

	providerBrand := os.Getenv("PROVIDER_BRAND")
	useLocal := os.Getenv("USE_LOCAL")

	envData, err := getEnvData(providerBrand, useLocal)
	if err != nil {
		log.Println(mainLogPrefix+" failed to get env data: ", err)
		return
	}

	GenerateProviderConf(providerBrand, envData)

}

func getEnvData(providerBrand string, useLocal string) (envData map[string]string, err error) {

	if useLocal != "" {
		GenerateLocalDevEnv(providerBrand)
		envData, err = GetLocalDevEnv()
	} else {
		envData, err = GetReleaseEnv(providerBrand)
	}

	return envData, err
}
