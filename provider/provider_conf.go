// Copyright 2021, Shoreline Software Inc.
// SPDX-License-Identifier: Apache-2.0

package provider

var ObjectConfigJsonStr = `
{
	"action": {
		"attributes": {
			"type":                    { "type": "string",     "computed": true, "value": "ACTION" },
			"name":                    { "type": "label",      "required": true, "forcenew": true, "skip": true },
			"command":                 { "type": "command",    "required": true, "primary": true, "refs": {"action":1} },
			"description":             { "type": "string",     "optional": true },
			"enabled":                 { "type": "intbool",    "optional": true, "default": false },
			"params":                  { "type": "string[]",   "optional": true },
			"resource_tags_to_export": { "type": "string_set", "optional": true },
			"res_env_var":             { "type": "string",     "optional": true },
			"resource_query":          { "type": "command",    "optional": true },
			"shell":                   { "type": "string",     "optional": true },
			"timeout":                 { "type": "int",        "optional": true, "default": 60000 },
			"file_deps":               { "type": "string_set", "optional": true, "refs": {"file":1} },
			"start_short_template":    { "type": "string",     "optional": true, "step": "start_step_class.short_template" },
			"start_long_template":     { "type": "string",     "optional": true, "step": "start_step_class.long_template" },
			"start_title_template":    { "type": "string",     "optional": true, "step": "start_step_class.title_template", "suppress_null_regex": "^started \\w*$" },
			"error_short_template":    { "type": "string",     "optional": true, "step": "error_step_class.short_template" },
			"error_long_template":     { "type": "string",     "optional": true, "step": "error_step_class.long_template" },
			"error_title_template":    { "type": "string",     "optional": true, "step": "error_step_class.title_template", "suppress_null_regex": "^failed \\w*$" },
			"complete_short_template": { "type": "string",     "optional": true, "step": "complete_step_class.short_template" },
			"complete_long_template":  { "type": "string",     "optional": true, "step": "complete_step_class.long_template" },
			"complete_title_template": { "type": "string",     "optional": true, "step": "complete_step_class.title_template", "suppress_null_regex": "^completed \\w*$" },
			"allowed_entities":        { "type": "string_set", "optional": true },
			"allowed_resources_query": { "type": "command",    "optional": true },
			"communication_workspace": { "type": "string",     "optional": true, "min_ver": "14.1.0", "step": "communication.workspace"},
			"communication_channel":   { "type": "string",     "optional": true, "min_ver": "14.1.0", "step": "communication.channel"},
			"editors":                 { "type": "string_set", "optional": true, "min_ver": "18.0.0" }
		}
	},

	"alarm": {
		"attributes": {
			"type":                   { "type": "string",   "computed": true, "value": "ALARM" },
			"name":                   { "type": "label",    "required": true, "forcenew": true, "skip": true },
			"fire_query":             { "type": "command",  "required": true, "primary": true, "refs": {"action":1} },
			"clear_query":            { "type": "command",  "optional": true, "refs": {"action":1} },
			"description":            { "type": "string",   "optional": true },
			"resource_query":         { "type": "command",  "optional": true },
			"enabled":                { "type": "intbool",  "optional": true, "default": false },
			"mute_query":             { "type": "string",   "optional": true },
			"resolve_short_template": { "type": "string",   "optional": true, "step": "clear_step_class.short_template" },
			"resolve_long_template":  { "type": "string",   "optional": true, "step": "clear_step_class.long_template" },
			"resolve_title_template": { "type": "string",   "optional": true, "step": "clear_step_class.title_template", "suppress_null_regex": "^cleared \\w*$" },
			"fire_short_template":    { "type": "string",   "optional": true, "step": "fire_step_class.short_template" },
			"fire_long_template":     { "type": "string",   "optional": true, "step": "fire_step_class.long_template" },
			"fire_title_template":    { "type": "string",   "optional": true, "step": "fire_step_class.title_template", "suppress_null_regex": "^fired \\w*$" },
			"condition_type":         { "type": "command",  "optional": true, "step": "condition_details.[0].condition_type" },
			"condition_value":        { "type": "string",   "optional": true, "step": "condition_details.[0].condition_value", "match_null": "0", "outtype": "float" },
			"metric_name":            { "type": "string",   "optional": true, "step": "condition_details.[0].metric_name" },
			"raise_for":              { "type": "command",  "optional": true, "step": "condition_details.[0].raise_for", "default": "local" },
			"check_interval_sec":     { "type": "command",  "optional": true, "step": "check_interval_sec", "default": 1, "outtype": "int" },
			"resource_type":          { "type": "resource", "optional": true, "step": "resource_type" },
			"family":                 { "type": "command",  "optional": true, "step": "config_data.family", "default": "custom" }
		}
	},

	"time_trigger": {
		"attributes": {
			"type":                   { "type": "string",   "computed": true, "value": "TIME_TRIGGER" },
			"name":                   { "type": "label",    "required": true, "forcenew": true, "skip": true },
			"fire_query":             { "type": "command",  "required": true, "primary": true },
			"start_date":             { "type": "string", 	"optional": true, "suppress_null_regex": "^[-+:0-9T]*$" },
			"end_date":               { "type": "string", 	"optional": true },
			"enabled":                { "type": "intbool",  "optional": true, "default": false }
		}
	},

	"bot": {
		"attributes": {
			"type":                    { "type": "string",  "computed": true, "value": "BOT" },
			"name":                    { "type": "label",   "required": true, "forcenew": true, "skip": true },
			"command":                 { "type": "command", "required": true, "primary": true, "refs": {"action":1, "alarm":1},
				"compound_in": "^\\s*if\\s*(?P<alarm_statement>.*?)\\s*then\\s*(?P<action_statement>.*?)\\s*fi\\s*$",
				"compound_out": "if ${alarm_statement} then ${action_statement} fi"
			},
			"description":             { "type": "string",  "optional": true },
			"enabled":                 { "type": "intbool", "optional": true, "default": false },
			"family":                  { "type": "command", "optional": true, "step": "config_data.family", "default": "custom" },
			"action_statement":        { "type": "command", "internal": true },
			"alarm_statement":         { "type": "command", "internal": true },
			"event_type":              { "type": "string",  "optional": true, "step": "event_type", "alias_out": "trigger_source", "match_null": "shoreline" },
			"monitor_id":              { "type": "string",  "optional": true, "step": "monitor_id", "alias_out": "external_trigger_id" },
			"alarm_resource_query":    { "type": "command", "optional": true },
			"#trigger_source":         { "type": "string",  "optional": true, "preferred_alias": "event_type", "step": "trigger_source", "default": "shoreline" },
			"#external_trigger_id":    { "type": "string",  "optional": true, "preferred_alias": "monitor_id", "step": "external_trigger_id", "default": "" },
			"communication_workspace": { "type": "string",  "optional": true, "min_ver": "14.1.0", "step": "communication.workspace"},
			"communication_channel":   { "type": "string",  "optional": true, "min_ver": "14.1.0", "step": "communication.channel"},
			"integration_name":        { "type": "string",  "optional": true, "min_ver": "15.0.0", "step": "integration_name"}
		}
	},

	"circuit_breaker": {
		"attributes": {
			"type":                    { "type": "string",  "computed": true, "value": "CIRCUIT_BREAKER" },
			"name":                    { "type": "label",   "required": true, "forcenew": true, "skip": true },
			"command":                 { "type": "command", "required": true, "primary": true, "forcenew": true, "refs": {"action":1},
				"compound_in": "^\\s*(?P<resource_query>.+)\\s*\\|\\s*(?P<action_name>[a-zA-Z_][a-zA-Z_]*)\\s*$",
				"compound_out": "${resource_query} | ${action_name}"
			},
			"breaker_type":            { "type": "string",  "optional": true },
			"hard_limit":              { "type": "int",     "required": true },
			"soft_limit":              { "type": "int",     "optional": true, "default": -1 },
			"duration":                { "type": "time_s",  "required": true },
			"fail_over":               { "type": "string",  "optional": true, "default": "safe"},
			"enabled":                 { "type": "bool",    "optional": true, "default": false },
			"action_name":             { "type": "command", "internal": true },
			"resource_query":          { "type": "command", "internal": true },
			"communication_workspace": { "type": "string",  "optional": true, "min_ver": "14.1.0", "step": "communication.workspace"},
			"communication_channel":   { "type": "string",  "optional": true, "min_ver": "14.1.0", "step": "communication.channel"}
		}
	},

	"file": {
		"attributes": {
			"type":             { "type": "string",   "computed": true, "value": "FILE" },
			"name":             { "type": "label",    "required": true, "forcenew": true, "skip": true },
			"destination_path": { "type": "string",   "required": true, "primary": true },
			"description":      { "type": "string",   "optional": true },
			"resource_query":   { "type": "string",   "required": true },
			"enabled":          { "type": "intbool",  "optional": true, "default": false },
			"input_file":       { "type": "string",   "optional": true, "skip": true, "not_stored": true, "conflicts": [ "inline_data" ] },
			"inline_data":      { "type": "string",   "optional": true, "skip": true, "not_stored": true, "conflicts": [ "input_file" ] },
			"file_data":        { "type": "string",   "computed": true, "outtype": "file" },
			"file_length":      { "type": "int",      "computed": true },
			"checksum":         { "type": "string",   "computed": true },
			"md5":              { "type": "string",   "optional": true, "proxy": "file_length,checksum,file_data" },
			"mode":             { "type": "string",   "optional": true, "min_ver": "23.0.0" },
			"owner":            { "type": "string",   "optional": true, "min_ver": "23.0.0" }
		}
	},

	"integration": {
		"internal": {
			"alias": {
				"key": "service_name",
				"map": {
					"newrelic": {
						"api_key": { "type": "string",   "optional": true, "step": "params_unpack.incident_management_api_key", "alias_out": "incident_management_api_key" },
						"api_url": { "type": "string",   "optional": true, "step": "params_unpack.incident_management_url",     "alias_out": "incident_management_url" }
					},
					"elastic": {
						"api_key": { "type": "string",   "optional": true, "step": "params_unpack.api_token", "alias_out": "api_token" },
						"api_url": { "type": "string",   "optional": true, "step": "params_unpack.url",       "alias_out": "url" }
					},
					"fluentbit_elastic": {
						"api_key": { "type": "string",   "optional": true, "step": "params_unpack.api_token", "alias_out": "api_token" },
						"api_url": { "type": "string",   "optional": true, "step": "params_unpack.url",       "alias_out": "url" }
					},
					"okta": {
						"api_key": { "type": "string",   "optional": true, "step": "params_unpack.api_token", "alias_out": "api_token" },
						"api_url": { "type": "string",   "optional": true, "step": "params_unpack.url",       "alias_out": "url" }
					}
				}
			}
		},
		"attributes": {
			"type":                        { "type": "string",   "computed": true, "value": "INTEGRATION" },
			"name":                        { "type": "label",    "required": true, "forcenew": true, "skip": true },
			"service_name":                { "type": "command",  "required": true, "primary": true, "forcenew": true, "skip": true },
			"serial_number":               { "type": "string",   "required": true },
			"permissions_user":            { "type": "string",   "optional": true, "match_null": "Shoreline" },
			"api_url":                     { "type": "string",   "optional": true, "step": "params_unpack.api_url" },
			"site_url":                    { "type": "string",   "optional": true, "step": "params_unpack.site_url", "min_ver": "19.0.0"},
			"api_key":                     { "type": "string",   "optional": true, "step": "params_unpack.api_key" },
			"app_key":                     { "type": "string",   "optional": true, "step": "params_unpack.app_key" },
			"dashboard_name":              { "type": "string",   "optional": true, "step": "params_unpack.dashboard_name", "max_ver": "25.99.999", "deprecated": true },
			"webhook_name":                { "type": "string",   "optional": true, "step": "params_unpack.webhook_name" },
			"##description":               { "type": "string",   "optional": true },

			"account_id":                  { "type": "string",   "optional": true, "step": "params_unpack.account_id" },
			"insights_collector_url":      { "type": "string",   "optional": true, "step": "params_unpack.insights_collector_url" },
			"insights_collector_api_key":  { "type": "string",   "optional": true, "step": "params_unpack.insights_collector_api_key" },
			"#incident_management_url":     { "type": "string",   "optional": true, "step": "params_unpack.incident_management_url" },
			"#incident_management_api_key": { "type": "string",   "optional": true, "step": "params_unpack.incident_management_api_key" },

			"cache_ttl":                   { "type": "int", "optional": true, "step": "params_unpack.cache_ttl" },
			"api_rate_limit":              { "type": "int", "optional": true, "step": "params_unpack.api_rate_limit" },
			"enabled":                     { "type": "intbool",  "optional": true, "default": false },

			"external_url":                { "type": "string",   "optional": true, "min_ver": "17.0.0", "step": "params_unpack.external_url" },
			"payload_paths":               { "type": "string_set", "optional": true, "min_ver": "28.4.0", "step": "params_unpack.payload_paths" },

			"cache_ttl_ms":                { "type": "int",    "optional": true, "min_ver": "18.0.0", "step": "params_unpack.cache_ttl_ms" },
			"subject":                     { "type": "string", "optional": true, "min_ver": "18.0.0", "step": "params_unpack.subject" },
			"credentials":                 { "type": "string", "optional": true, "min_ver": "18.0.0", "step": "params_unpack.credentials" },

			"tenant_id":                   { "type": "string", "optional": true, "min_ver": "18.0.0", "step": "params_unpack.tenant_id" },
			"client_id":                   { "type": "string", "optional": true, "min_ver": "18.0.0", "step": "params_unpack.client_id" },
			"client_secret":               { "type": "string", "optional": true, "min_ver": "18.0.0", "step": "params_unpack.client_secret" },

			"idp_name":                    { "type": "string", "optional": true, "min_ver": "22.0.0", "step": "params_unpack.idp_name" },

			"api_certificate":             { "type": "string", "optional": true, "min_ver": "28.1.0", "step": "params_unpack.api_certificate" }
		}
	},

	"metric": {
		"attributes": {
			"type":           { "type": "string",   "computed": true, "value": "METRIC" },
			"name":           { "type": "label",    "required": true, "forcenew": true, "skip": true },
			"value":          { "type": "command",  "required": true, "primary": true, "alias_out": "val" },
			"description":    { "type": "string",   "optional": true },
			"units":          { "type": "string",   "optional": true },
			"resource_type":  { "type": "resource", "optional": true }
		}
	},

	"notebook": {
		"attributes": {
			"type":                    { "type": "string",     "computed": true, "value": "NOTEBOOK" },
			"name":                    { "type": "label",      "required": true, "forcenew": true, "skip": true },
			"data":                    { "type": "b64json",    "optional": true, "step": ".", "primary": true, "deprecated": true,
																	 "omit":       { "cells": "dynamic_cell_fields", ".": "dynamic_fields" },
																	 "omit_items": { "external_params": "dynamic_params" },
																	 "cast":       { "params": "string[]", "params_values": "string[]" },
																	 "force_set":  [ "allowed_entities", "approvers", "is_run_output_persisted",
																	                 "communication_workspace", "communication_channel" ],
																	 "skip_diff":  [ "allowedUsers", "isRunOutputPersisted", "approvers", "communication",
																	                 "interactive_state", "name", "allowedResourcesQuery", "timeoutMs" ],
																	 "outtype": "json"
			                           },
			"#cells":                   { "type": "b64json",    "optional": true, "step": "cells", "outtype": "json", "conflicts": ["data"]},
			"#params":                  { "type": "b64json",    "optional": true, "step": "params", "outtype": "json", "conflicts": ["data"]},
			"#external_params":         { "type": "b64json",    "optional": true, "step": "external_params", "outtype": "json", "conflicts": ["data"]},
			"#enabled":                 { "type": "bool",       "optional": true, "step": "enabled", "default": true, "conflicts": ["data"]},
			"cells":                   { "type": "b64json",    "optional": true, "step": "skip.param", "outtype": "json", "conflicts": ["data"]},
			"params":                  { "type": "b64json",    "optional": true, "outtype": "json", "conflicts": ["data"], "match_null": "[]"},
			"external_params":         { "type": "b64json",    "optional": true, "outtype": "json", "conflicts": ["data"], "match_null": "[]"},
			"enabled":                 { "type": "bool",       "optional": true, "default": true, "conflicts": ["data"]},
			"description":             { "type": "string",     "optional": true },
			"timeout_ms":              { "type": "unsigned",   "optional": true, "default": 60000 },
			"allowed_entities":        { "type": "string_set", "optional": true },
			"approvers":               { "type": "string_set", "optional": true },
			"resource_query":          { "type": "string",     "optional": true, "deprecated_for": "allowed_resources_query" },
			"is_run_output_persisted": { "type": "bool",       "optional": true, "step": "is_run_output_persisted", "default": true, "min_ver": "12.3.0" },
			"allowed_resources_query": { "type": "command",    "optional": true, "replaces": "resource_query", "min_ver": "12.3.0" },
			"communication_workspace": { "type": "string",     "optional": true, "min_ver": "12.5.0", "step": "communication_workspace" },
			"communication_channel":   { "type": "string",     "optional": true, "min_ver": "12.5.0", "step": "communication_channel" },
			"labels":                  { "type": "string_set", "optional": true, "min_ver": "16.0", "step": "labels" },
			"editors":                 { "type": "string_set", "optional": true, "min_ver": "15.1.0" },
			"communication_cud_notifications":       { "type": "bool", "default": true, "optional": true, "min_ver": "17.0.0", "step": "communication_cud_notifications" },
			"communication_approval_notifications":  { "type": "bool", "default": true, "optional": true, "min_ver": "17.0.0", "step": "communication_approval_notifications" },
			"communication_execution_notifications": { "type": "bool", "default": true, "optional": true, "min_ver": "17.0.0", "step": "communication_execution_notifications" },
			"filter_resource_to_action": { "type": "bool", "default": false, "optional": true, "min_ver": "28.0.0", "force_update": true },
			"secret_names": { "type": "string_set", "optional": true, "min_ver": "28.1.0" }
		}
	},

	"resource": {
		"attributes": {
			"type":            { "type": "string",   "computed": true, "value": "RESOURCE" },
			"name":            { "type": "label",    "required": true, "forcenew": true, "skip": true },
			"value":           { "type": "command",  "required": true, "primary": true },
			"description":     { "type": "string",   "optional": true },
			"params":          { "type": "string[]", "optional": true },
			"#units":          { "type": "string",   "optional": true },
			"#resource_type":  { "type": "resource", "optional": true },
			"#user":           { "type": "string",   "optional": true },
			"#read_only":      { "type": "bool",     "optional": true }
		}
	},

	"principal": {
		"attributes": {
			"type":                  { "type": "string",   "computed": true, "value": "PRINCIPAL" },
			"name":                  { "type": "label",    "required": true, "forcenew": true, "skip": true },
			"identity":              { "type": "string",   "required": true, "primary": true },
			"view_limit":            { "type": "int",      "optional": true },
			"action_limit":          { "type": "int",      "optional": true },
			"execute_limit":         { "type": "int",      "optional": true },
			"configure_permission":  { "type": "intbool",  "optional": true },
			"administer_permission": { "type": "intbool",  "optional": true },
			"idp_name":              { "type": "string",   "optional": true, "min_ver": "22.0.0" }
		}
	},

	"system_settings": {
		"internal": {
			"singleton": "system_settings",
			"read_single_attr": true,
			"no_create": true,
			"no_delete": true
		},
		"attributes": {
			"type":                                             { "type": "string",   "computed": true, "value": "SYSTEM_SETTINGS" },
			"name":                                             { "type": "label",    "required": true, "forcenew": true, "skip": true },
			"administrator_grants_create_user":                 { "type": "bool",     "optional": true, "default": true },
			"administrator_grants_create_user_token":           { "type": "bool",     "optional": true, "default": true },
			"administrator_grants_regenerate_user_token":       { "type": "bool",     "optional": true, "default": true },
			"administrator_grants_read_user_token":             { "type": "bool",     "optional": true, "default": true },
			"approval_feature_enabled":                         { "type": "bool",     "optional": true, "default": true },
			"notebook_ad_hoc_approval_request_enabled":         { "type": "bool",     "optional": true, "default": true, "deprecated_for": "runbook_ad_hoc_approval_request_enabled" },
			"runbook_ad_hoc_approval_request_enabled":          { "type": "bool",     "optional": true, "default": true, "min_ver": "25.0.0", "replaces": "notebook_ad_hoc_approval_request_enabled" },
			"notebook_approval_request_expiry_time":            { "type": "int",      "optional": true, "default": 60, "deprecated_for": "runbook_approval_request_expiry_time"  },
			"runbook_approval_request_expiry_time":             { "type": "int",      "optional": true, "default": 60, "min_ver": "25.0.0", "replaces": "notebook_approval_request_expiry_time" },
			"notebook_run_approval_expiry_time":                { "type": "int",      "optional": true, "default": 60, "deprecated_for": "run_approval_expiry_time"  },
			"run_approval_expiry_time":                         { "type": "int",      "optional": true, "default": 60, "min_ver": "25.0.0", "replaces": "notebook_run_approval_expiry_time" },
			"approval_editable_allowed_resource_query_enabled": { "type": "bool",     "optional": true, "default": true },
			"approval_allow_individual_notification":           { "type": "bool",     "optional": true, "min_ver": "17.0.0", "default": true },
			"approval_optional_request_ticket_url":             { "type": "bool",     "optional": true, "min_ver": "17.0.0", "default": false },
			"time_trigger_permissions_user":                    { "type": "string",   "optional": true, "min_ver": "19.1.0", "default": "Shoreline"},
			"external_audit_storage_enabled":                   { "type": "bool",     "optional": true, "default": false },
			"external_audit_storage_type":                      { "type": "command",  "optional": true, "default": "ELASTIC" },
			"#external_audit_storage_url":                      { "type": "string",   "optional": true },
			"#external_audit_storage_api_token":                { "type": "string",   "optional": true },
			"external_audit_storage_batch_period_sec":          { "type": "int",      "optional": true, "default": 5 },
			"environment_name":                                 { "type": "string",   "optional": true, "default": "" },
			"environment_name_background":                      { "type": "string",   "optional": true, "min_ver": "18.0.0", "default": "#EF5350" },
			"param_value_max_length":                           { "type": "int",      "optional": true, "min_ver": "19.0.0", "default": 10000 },
			"parallel_notebook_runs_fired_by_time_triggers":    { "type": "int",      "optional": true, "default": 10, "min_ver": "20.0.2", "deprecated_for": "parallel_runs_fired_by_time_triggers"  },
			"parallel_runs_fired_by_time_triggers":             { "type": "int",      "optional": true, "default": 10, "min_ver": "25.0.0", "replaces": "parallel_notebook_runs_fired_by_time_triggers" },
			"maintenance_mode_enabled":                         { "type": "bool",     "optional": true, "min_ver": "25.1.0", "default": false },
			"allowed_tags":                                     { "type": "string_set", "optional": true, "min_ver": "27.2.0" },
			"skipped_tags":                                     { "type": "string_set", "optional": true, "min_ver": "27.2.0" },
			"managed_secrets": 									{ "type": "string", "optional": true, "min_ver": "28.1.0", "default": "LOCAL" }
		}
	},

	"report_template": {
       "attributes": {
           "type":            { "type": "string",   "computed": true, "value": "REPORT_TEMPLATE" },
           "name":            { "type": "label",    "required": true, "forcenew": true, "skip": true},
           "blocks":          { "type": "b64json",  "required": true, "outtype": "json", "primary": true},
		   "links":           { "type": "b64json",  "optional": true, "outtype": "json", "default": "[]", "step": "links"}
       }
    },

	"dashboard": {
       "attributes": {
           "type":               { "type": "string",     "computed": true, "value": "DASHBOARD" },
           "name":               { "type": "label",      "required": true, "forcenew": true, "skip": true},
           "dashboard_type":     { "type": "string",     "required": true, "primary": true },
           "resource_query":     { "type": "string",     "optional": true, "step": "dashboard_configuration.resource_query" },
           "groups":             { "type": "b64json",    "optional": true, "step": "dashboard_configuration.groups", "outtype": "json" },
           "values":             { "type": "b64json",    "optional": true, "step": "dashboard_configuration.values", "outtype": "json" },
           "other_tags":         { "type": "string_set", "optional": true, "step": "dashboard_configuration.other_tags" },
           "identifiers":        { "type": "string_set", "optional": true, "step": "dashboard_configuration.identifiers" }
       }
    },

	"docs": {
		"objects": {
			"action":    "A command that can be run.\n\nSee the Shoreline [Actions Documentation](https://docs.shoreline.io/actions) for more info.",
			"alarm":     "A condition that triggers Alerts or Actions.\n\nSee the Shoreline [Alarms Documentation](https://docs.shoreline.io/alarms) for more info.",
			"time_trigger": "A condition that triggers Notebooks.",
			"bot":       "An automation that ties an Action to an Alert.\n\nSee the Shoreline [Bots Documentation](https://docs.shoreline.io/bots) for more info.",
			"circuit_breaker": "An automatic rate limit on actions.\n\nSee the Shoreline [CircuitBreakers Documentation](https://docs.shoreline.io/circuit_breakers) for more info.",
			"file":      "A datafile that is automatically copied/distributed to defined Resources.\n\nSee the Shoreline [OpCp Documentation](https://docs.shoreline.io/op/commands/cp) for more info.",
			"integration":  "A third-party integration (e.g. DataDog, NewRelic, etc) .\n\nSee the Shoreline [Metrics Documentation](https://docs.shoreline.io/integrations) for more info.",
			"metric":    "A periodic measurement of a system property.\n\nSee the Shoreline [Metrics Documentation](https://docs.shoreline.io/metrics) for more info.",
			"notebook":  "An interactive notebook of Op commands and user documentation .\n\nSee the Shoreline [Notebook Documentation](https://docs.shoreline.io/ui/notebooks) for more info.",
			"principal": "An authorization group (e.g. Okta groups). Note: Admin privilege (in Shoreline) to create principal objects.",
			"resource":  "A server or compute resource in the system (e.g. host, pod, container).\n\nSee the Shoreline [Resources Documentation](https://docs.shoreline.io/platform/resources) for more info.",
			"system_settings":  "System-level settings. Note: there must only be one instance of this terraform resource named 'system_settings'.\n\nSee the Shoreline [Settings Documentation](https://docs.shoreline.io/platform/settings) for more info.",
			"report_template":  "A resource report template. Note: Configure privilege (in Shoreline) to create report template objects.",
			"dashboard": "A platform for visualizing resources and their associated tags."
		},

		"attributes": {
			"name":                    "The name/symbol for the object within Shoreline and the op language (must be unique, only alphanumeric/underscore).",
			"type":                    "The type of object (i.e., Alarm, Action, Bot, Metric, Resource, or File).",
			"action_limit":            "The number of simultaneous actions allowed for a permissions group.",
			"administer_permission":   "If a permissions group is allowed to perform \"administer\" actions.",
			"allowed_entities":        "The list of users who can run an action or notebook. Any user can run if left empty.",
			"allowed_resources_query": "The list of resources on which an action or notebook can run. No restriction, if left empty.",
			"cells":                   "The data cells inside a notebook. Defined as a list of JSON objects. These may be either Markdown or Op commands.",
			"check_interval":          "Interval (in seconds) between Alarm evaluations.",
			"checksum":                "Cryptographic hash (e.g. md5) of a File Resource.",
			"clear_query":             "The Alarm's resolution condition.",
			"command":                 "A specific action to run.",
			"complete_long_template":  "The long description of the Action's completion.",
			"complete_short_template": "The short description of the Action's completion.",
			"complete_title_template": "UI title of the Action's completion.",
			"condition_type":          "Kind of check in an Alarm (e.g. above or below) vs a threshold for a Metric.",
			"condition_value":         "Switching value (threshold) for a Metric in an Alarm.",
			"configure_permission":    "If a permissions group is allowed to perform \"configure\" actions.",
			"data":                    "The JSON representation of a Notebook. If this field is used, then the JSON should only contain these four fields: cells, params, external_params and enabled.",
			"description":             "A user-friendly explanation of an object.",
			"destination_path":        "Target location for a copied distributed File object.  See [Op: cp](https://docs.shoreline.io/op/commands/cp).",
			"enabled":                 "If the object is currently enabled or disabled.",
			"end_date":                "When the trigger condition stops firing. (defaults to unset, e.g. no stop date). The accepted format is ISO8601, e.g. '2029-02-17T08:08:01'.",
			"error_long_template":     "The long description of the Action's error condition.",
			"error_short_template":    "The short description of the Action's error condition.",
			"error_title_template":    "UI title of the Action's error condition.",
			"event_type":              "Used to tag 'datadog' monitor triggers vs 'shoreline' alarms (default).",
			"execute_limit":           "The number of simultaneous linux (shell) commands allowed for a permissions group.",
			"family":                  "General class for an Action or Bot (e.g., custom, standard, metric, or system check).",
			"file_data":               "Internal representation of a distributed File object's data (computed).",
			"file_deps":               "file object dependencies.",
			"file_length":             "Length, in bytes, of a distributed File object (computed)",
			"fire_long_template":      "The long description of the Alarm's triggering condition.",
			"fire_query":              "The trigger condition for an Alarm (general expression) or the TimeTrigger (e.g. 'every 5m').",
			"fire_short_template":     "The short description of the Alarm's triggering condition.",
			"fire_title_template":     "UI title of the Alarm's triggering condition.",
			"identity":                "The email address or provider's (e.g. Okta) group-name for a permissions group.",
			"identifiers":             "A list of additional tags that will be used to identify certain resources. They will be displayed before the tags_sequence column.",
			"idp_name":                "The Identity Provider's name.",
			"input_file":              "The local source of a distributed File object. (conflicts with inline_data)",
			"inline_data":             "The inline file data of a distributed File object. (conflicts with input_file)",
			"is_run_output_persisted": "A boolean value denoting whether or not cell outputs should be persisted when running a notebook",
			"labels":                  "A list of strings by which notebooks can be grouped.",
			"md5":                     "The md5 checksum of a file, e.g. filemd5(\"${path.module}/data/example-file.txt\")",
			"metric_name":             "The Alarm's triggering Metric.",
			"mode":                    "The File's permissions, like 'chmod', in octal (e.g. '0644').",
			"owner":                   "The File's ownership, like 'chown' (e.g. 'user:group').",
			"monitor_id":              "For 'datadog' monitor triggered bots, the DD monitor identifier.",
			"mute_query":              "The Alarm's mute condition.",
			"params":                  "Named variables to pass to an object (e.g. an Action).",
			"external_params":         "Notebook parameters defined via with a JSON path used to extract the parameter's value from an external payload, such as an Alertmanager alert.",
			"raise_for":               "Where an Alarm is raised (e.g., local to a resource, or global to the system).",
			"res_env_var":             "Result environment variable ... an environment variable used to output values through.",
			"resolve_long_template":   "The long description of the Alarm's resolution.",
			"resolve_short_template":  "The short description of the Alarm's resolution.",
			"resolve_title_template":  "UI title of the Alarm's' resolution.",
			"resource_query":          "A set of Resources (e.g. host, pod, container), optionally filtered on tags or dynamic conditions.",
			"shell":                   "The commandline shell to use (e.g. /bin/sh).",
			"start_date":              "When the trigger condition starts firing (defaults to creation/update time of the trigger). The accepted format is ISO8601, e.g. '2024-02-17T08:08:01'.",
			"start_long_template":     "The long description when starting the Action.",
			"start_short_template":    "The short description when starting the Action.",
			"start_title_template":    "UI title of the start of the Action.",
			"timeout":                 "Maximum time to wait, in milliseconds.",
			"units":                   "Units of a Metric (e.g., bytes, blocks, packets, percent).",
			"value":                   "The Op statement that defines a Metric or Resource.",
			"view_limit":              "The number of simultaneous metrics allowed for a permissions group.",
			"is_run_output_persisted": "A boolean value denoting whether or not cell outputs should be persisted when running a notebook",
			"communication_workspace": "A string value denoting the slack workspace where notifications related to the object should be sent to.",
			"communication_channel":   "A string value denoting the slack channel where notifications related to the object should be sent to.",
			"service_name":            "The name of a 3rd-party service to integrate with (e.g. 'datadog', or 'newrelic').",
			"account_id":              "Account ID for a 3rd-party service integration.",
			"api_key":                 "API key for a 3rd-party service integration.",
			"api_url":                 "API url for a 3rd-party service integration.",
			"site_url":                "Site/Application url for a 3rd-party service integration.",
			"app_key":                 "Application key for a 3rd-party service integration.",
			"insights_collector_url":  "Insights url for a 3rd-party service integration.",
			"insights_collector_api_key": "Insights key for a 3rd-party service integration.",
			"permissions_user":        "The user which 3rd-party service integration remediations run as (default 'Shoreline').",
			"dashboard_name":          "The name of a dashboard for 3rd-party service integration (datadog).",
			"webhook_name":            "The name of a webhook for 3rd-party service integration (datadog).",
			"cache_ttl":               "The amount of time group memberships will be cached (in milliseconds).",
			"api_rate_limit":          "The number of API calls a client is able to make in a minute.",
			"administrator_grants_create_user":                 "System setting controlling if administrators can create users.",
			"administrator_grants_create_user_token":           "System setting controlling if administrators can create user access tokens.",
			"administrator_grants_regenerate_user_token":       "System setting controlling if administrators can update user access tokens.",
			"administrator_grants_read_user_token":             "System setting controlling if administrators can view user access tokens.",
			"approval_feature_enabled":                         "System setting controlling if notebook approvals are enabled.",
			"notebook_ad_hoc_approval_request_enabled":         "System setting controlling if approvals are enabled for ad-hoc notebook execution.",
			"approval_editable_allowed_resource_query_enabled": "System setting controlling if notebook resource queries can be modified on approved executions.",
			"approval_allow_individual_notification":           "System setting controlling if approvals notifications are sent to individual users, in case no specific notebook communication setting is defined.",
			"approval_optional_request_ticket_url":	           	"System setting controlling if the ticket url is optional when creating an approval request.",
			"time_trigger_permissions_user":                    "System setting for the user that time-triggered notebooks run as.",
			"external_audit_storage_enabled":                   "System setting controlling if audit information is stored in an alternate location.",
			"external_audit_storage_url":                       "System setting for alternate audit storage URL.",
			"external_audit_storage_type":                      "System setting for alternate audit storage type (e.g. 'ELASTIC').",
			"external_audit_storage_api_token":                 "System setting for alternate audit storage API access token.",
			"external_audit_storage_batch_period_sec":          "System setting for alternate audit storage batching interval (in seconds).",
			"notebook_approval_request_expiry_time":            "System setting for maximum wait for approval after request (in minutes).",
			"notebook_run_approval_expiry_time":                "System setting for maximum wait for execution after approval (in minutes).",
			"environment_name":                                 "System setting for the name of the environment.",
			"environment_name_background":                      "System setting for the background colour of the environment name. The format is #<6-digit hex>",
			"param_value_max_length":                           "System setting controlling the maximum allowable length for a notebook's parameter",
			"parallel_notebook_runs_fired_by_time_triggers":    "System setting controlling the maximum number of different parallel notebook runs initiated via time triggers",
			"maintenance_mode_enabled":                        	"System setting that when enabled, rejects new runs, allowing ongoing tasks to complete before stopping.",
			"allowed_tags":                                     "Defines a list of tags that are allowed on agent tag ingestion",
			"skipped_tags":                                     "Defines a list of tags that are skipped on agent tag ingestion",
			"managed_secrets": 									"System setting that discriminates between usage of external vaults and the built in one.",
			"integration_name":                                 "The name/symbol of a Shoreline integration involved in triggering the bot.",
			"editors":                                          "List of users who can edit the object (with configure permission). Empty maps to all users.",
			"communication_cud_notifications":                  "Enables slack notifications for create/update/delete operations. (Requires workspace and channel.)",
			"communication_approval_notifications":             "Enables slack notifications for approvals operations. (Requires workspace and channel.)",
			"communication_execution_notifications":            "Enables slack notifications for the object executions. (Requires workspace and channel.)",
			"filter_resource_to_action":                        "Determines whether parameters containing resources are exported to actions.",
			"external_url": 									"External url for a 3rd-party service integration.",
			"cache_ttl_ms":            "The amount of time group memberships will be cached (in milliseconds).",
			"subject":                 "The subject whose authentication details is used for a 3rd-party service integration (google cloud identity).",
			"credentials":             "The credentials used for a 3rd-party service integration (google cloud identity), encoded in base64.",
			"tenant_id":               "Tenant id for a 3rd-party service integration (Microsoft Entra ID).",
			"client_id":               "Application id for a 3rd-party service integration (Microsoft Entra ID).",
			"client_secret":           "Client secret for a 3rd-party service integration (Microsoft Entra ID).",
			"blocks":           	   "The JSON encoded blocks of the report template.",
			"links":           	   	   "The JSON encoded links of a report template with other report templates.",
			"dashboard_type":          "Specifies the type of the dashboard configuration. Currently, only 'TAGS_SEQUENCE' is supported.",
			"secret_names":            "A list of strings that contains the name of the secrets that are used in the runbook.",
			"groups":                  "A JSON-encoded list of groups in the dashboard configuration. Each group is an object with 'name' (the group's name) and 'tags' (a list of tag names belonging to the group).",
			"values":                  "A JSON-encoded list of objects defining the values and their associated colors in the dashboard configuration. Each object contains: 'color' (the color associated with the values) and 'values' (a list of values corresponding to specific tags).",
			"other_tags":              "A list of additional tags that will be displayed for the resources.",
			"api_certificate":         "API certificate for a 3rd-party service integration.",
			"payload_paths":           "A list of JSON paths to extract values from the payload of an alert."
		}
	}
}
`
