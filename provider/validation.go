package provider

import "errors"

func validateShorelineNotebookDataField(data interface{}) error {
	allowedKeys := map[string]bool{"cells": true, "params": true, "external_params": true, "enabled": true}
	for k := range CastToObject(data).(map[string]interface{}) {
		if ok := allowedKeys[k]; !ok {
			return errors.New("shoreline_notebook.data field is expected to only have the following keys: cells, params, external_params and enabled, but got: " + k)
		}
	}
	return nil
}
