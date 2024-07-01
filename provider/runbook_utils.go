// Copyright 2021, Shoreline Software Inc.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// buildRunbookDataObject builds a JSON containing the core runbook data. expects at least "cells" to be present
func buildRunbookDataObject(d *schema.ResourceData, cells interface{}) (interface{}, error) {
	runbookData := map[string]interface{}{}

	cellsData, err := buildCellsData(cells)
	if err != nil {
		return nil, err
	}
	runbookData["cells"] = cellsData

	parametersData, err := buildParametersData(d)
	if err != nil {
		return nil, err
	}
	runbookData["params"] = parametersData

	externalParametersData, err := buildExternalParametersData(d)
	if err != nil {
		return nil, err
	}
	runbookData["external_params"] = externalParametersData

	enabled, exists := d.GetOk("enabled")
	if !exists {
		enabled = true
	}
	runbookData["enabled"] = enabled

	// return the json encoded runbookData
	encodedRunbookData, err := json.Marshal(runbookData)
	if err != nil {
		return nil, fmt.Errorf("error encoding runbook data: %v", err)
	}
	return string(encodedRunbookData), nil
}

func buildCellsData(cells interface{}) (interface{}, error) {
	var decodedCells []interface{}
	var cellContent map[string]interface{}
	cellsData := []interface{}{}

	appendActionLog(fmt.Sprintf("building runbook cells from: %v\n", cells))

	err := json.Unmarshal([]byte(cells.(string)), &decodedCells)
	if err != nil {
		return nil, fmt.Errorf("error decoding cells: %v", err)
	}

	for _, cell := range decodedCells {
		markdownContent := cell.(map[string]interface{})["md"]
		oplangContent := cell.(map[string]interface{})["op"]
		if markdownContent == nil && oplangContent == nil {
			return nil, fmt.Errorf(`runbook cell must specify either an oplang command or markdown content using the "md" or "op" fields`)
		}
		if markdownContent != nil && oplangContent != nil {
			return nil, fmt.Errorf("runbook cell cannot have both markdown and oplang content")
		}

		enabled, exists := cell.(map[string]interface{})["enabled"]
		if !exists {
			enabled = true
		}

		if markdownContent != nil {
			cellContent = map[string]interface{}{
				"content": markdownContent,
				"enabled": enabled,
				"type":    "MARKDOWN",
				"name":    "unnamed",
			}
		} else {
			cellContent = map[string]interface{}{
				"content": oplangContent,
				"enabled": enabled,
				"type":    "OP_LANG",
				"name":    "unnamed",
			}
		}

		cellsData = append(cellsData, cellContent)
	}

	return cellsData, nil
}

func buildParametersData(d *schema.ResourceData) ([]interface{}, error) {
	var decodedParameters []interface{}
	parametersData := []interface{}{}
	parameters, exists := d.GetOk("params")

	appendActionLog(fmt.Sprintf("building runbook parameters from: %v\n", parameters))

	if !exists {
		return []interface{}{}, nil
	}

	err := json.Unmarshal([]byte(parameters.(string)), &decodedParameters)
	if err != nil {
		return nil, fmt.Errorf("error decoding parameters: %v", err)
	}

	for _, parameter := range decodedParameters {
		name, exists := parameter.(map[string]interface{})["name"]
		if !exists {
			return nil, fmt.Errorf("parameter name is required")
		}
		required, exists := parameter.(map[string]interface{})["required"]
		if !exists {
			required = true
		}
		value, exists := parameter.(map[string]interface{})["value"]
		if !exists {
			value = ""
		}
		export, exists := parameter.(map[string]interface{})["export"]
		if !exists {
			export = false
		}

		parameterData := map[string]interface{}{
			"export":   export, // false by default
			"name":     name,
			"required": required, // true by default
			"value":    value,    // empty string by default
		}

		parametersData = append(parametersData, parameterData)
	}

	return parametersData, nil
}

func buildExternalParametersData(d *schema.ResourceData) ([]interface{}, error) {
	var decodedExternalParameters []interface{}
	externalParametersData := []interface{}{}
	externalParameters, exists := d.GetOk("external_params")

	appendActionLog(fmt.Sprintf("building runbook external parameters from: %v\n", externalParameters))

	if !exists {
		return []interface{}{}, nil
	}

	err := json.Unmarshal([]byte(externalParameters.(string)), &decodedExternalParameters)
	if err != nil {
		return nil, fmt.Errorf("error decoding external parameters: %v", err)
	}

	for _, externalParameter := range decodedExternalParameters {
		name, exists := externalParameter.(map[string]interface{})["name"]
		if !exists {
			return nil, fmt.Errorf("external parameter name is required")
		}
		source, exists := externalParameter.(map[string]interface{})["source"]
		if !exists {
			return nil, fmt.Errorf("external parameter source is required")
		}
		value, exists := externalParameter.(map[string]interface{})["value"]
		if !exists {
			value = ""
		}

		externalParameterData := map[string]interface{}{
			"name":   name,
			"source": source,
			"value":  value,
		}

		externalParametersData = append(externalParametersData, externalParameterData)
	}

	return externalParametersData, nil
}
