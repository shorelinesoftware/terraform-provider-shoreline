// Copyright 2021, Shoreline Software Inc.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	//"regexp"
	"testing"

	//"github.com/hashicorp/terraform-plugin-sdk/acctest"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"math/rand"
)

func getProviderConfigString() string {
	return `
	provider "shoreline" {
		url = "https://opsstage.us.api.shoreline-stage.io"
		retries = 2
		debug = true
	}
`
}

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"shoreline": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

//func TestMain(m *testing.M) {
//  acctest.UseBinaryDriver("shoreline", New("dev"))
//  resource.TestMain(m)
//}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

func RandStringFromCharSet(strlen int, charSet string) string {
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(result)
}

func RandomAlphaPrefix(slen int) string {
	alpha := "abcdefghijklmnopqrstuvwxyz"
	return RandStringFromCharSet(slen, alpha)
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
// Action

func TestAccResourceAction(t *testing.T) {
	//t.Skip("resource not yet implemented, remove this once you add your own code")
	pre := RandomAlphaPrefix(5)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfigString() + getAccResourceAction(pre, false),
				Check: resource.ComposeTestCheckFunc(
					//resource.TestMatchResourceAttr( "shoreline_action.ls_action", "name", regexp.MustCompile("^ba")),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "name", pre+"_ls_action"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "command", "`ls ${dir}; export FOO='bar'`"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "description", "List some files"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "enabled", "true"),
					//resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "params", "[\"dir\"]"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "start_title_template", "my_action started"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "complete_title_template", "my_action completed"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "error_title_template", "my_action failed"),
				),
			},
			{
				Config: getProviderConfigString() + getAccResourceAction(pre, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "name", pre+"_ls_action"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "command", "`ls ${dir}; export FOO='bar'`"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "enabled", "true"),
					//resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "params", "[\"dir\"]"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "start_title_template", "my_action started"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "complete_title_template", "my_action completed"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "error_title_template", "my_action failed"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "resource_query", "host"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "timeout", "20"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "shell", "/bin/bash"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "res_env_var", "FOO"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "start_short_template", "started"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "start_long_template", "started..."),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "complete_short_template", "completed"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "complete_long_template", "completed..."),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "error_short_template", "failed"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "error_long_template", "failed..."),
				),
			},
		},
	})

	//resource.UnitTest(t, resource.TestCase{
	//	PreCheck:          func() { testAccPreCheck(t) },
	//	ProviderFactories: providerFactories,
	//})
}

func getAccResourceAction(prefix string, full bool) string {
	extra := `
			resource_query = "host"
			timeout = 20
			shell = "/bin/bash"
			res_env_var = "FOO"
			start_short_template    = "started"
			start_long_template    = "started..."
			complete_short_template = "completed"
			complete_long_template = "completed..."
			error_short_template    = "failed"
			error_long_template    = "failed..."
`
	if !full {
		extra = ""
	}
	return `
		resource "shoreline_action" "` + prefix + `_ls_action" {
			name = "` + prefix + `_ls_action"
			command = "` + "`ls $${dir}; export FOO='bar'`" + `"
			description = "List some files"
			params = ["dir"]
			enabled = true
			start_title_template    = "my_action started"
			complete_title_template = "my_action completed"
			error_title_template    = "my_action failed"
			` + extra + `
		}
`
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
// Alarm

func TestAccResourceAlarm(t *testing.T) {
	pre := RandomAlphaPrefix(5)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfigString() + getAccResourceAlarm(pre, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "name", pre+"_cpu_alarm"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "fire_query", "( cpu_usage > 0 | sum ( 5 ) ) >= 2"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "clear_query", "( cpu_usage < 0 | sum ( 5 ) ) >= 2"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "description", "Watch CPU usage."),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "resource_query", "host"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "enabled", "true"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "fire_title_template", "alarm fired"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "resolve_title_template", "alarm resolved"),
				),
			},
			{
				Config: getProviderConfigString() + getAccResourceAlarm(pre, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "name", pre+"_cpu_alarm"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "fire_query", "( cpu_usage > 0 | sum ( 5 ) ) >= 2"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "clear_query", "( cpu_usage < 0 | sum ( 5 ) ) >= 2"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "description", "Watch CPU usage."),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "resource_query", "host"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "enabled", "true"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "fire_short_template", "fired"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "fire_long_template", "fired..."),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "resolve_short_template", "resolved"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "resolve_long_template", "resolved..."),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "raise_for", "local"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "check_interval", "50"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "compile_eligible", "false"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "resource_type", "host"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "family", "custom"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "metric_name", "cpu_usage"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "condition_type", "above"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "condition_value", "0"),
				),
			},
		},
	})
}

func getAccResourceAlarm(prefix string, full bool) string {
	extra := `
			fire_short_template     = "fired"
			fire_long_template      = "fired..."
			resolve_short_template  = "resolved"
			resolve_long_template   = "resolved..."
			raise_for               = "local"
			check_interval          = 50
			compile_eligible        = false
			resource_type           = "host"
			family                  = "custom"
			metric_name             = "cpu_usage"
			condition_type          = "above"
			condition_value         = "0"
`
	if !full {
		extra = ""
	}
	return `
		resource "shoreline_alarm" "` + prefix + `_cpu_alarm" {
			name = "` + prefix + `_cpu_alarm"
	    fire_query = "(cpu_usage > 0 | sum(5)) >= 2"
	    clear_query = "(cpu_usage < 0 | sum(5)) >= 2"
	    description = "Watch CPU usage."
	    resource_query = "host"
	    enabled = true
			fire_title_template     = "alarm fired"
			resolve_title_template  = "alarm resolved"
			` + extra + `
		}
`
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
// Bot

func TestAccResourceBot(t *testing.T) {
	pre := RandomAlphaPrefix(5)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfigString() + getAccResourceAction(pre, false) + getAccResourceAlarm(pre, false) + getAccResourceBot(pre),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("shoreline_bot."+pre+"_cpu_bot", "name", pre+"_cpu_bot"),
					//resource.TestCheckResourceAttr("shoreline_bot."+pre+"_cpu_bot", "alarm_statement", pre+"_cpu_alarm"),
					resource.TestCheckResourceAttr("shoreline_bot."+pre+"_cpu_bot", "command", "if "+pre+"_cpu_alarm then "+pre+"_ls_action ( \"/tmp\" ) fi"),
					resource.TestCheckResourceAttr("shoreline_bot."+pre+"_cpu_bot", "description", "Act on CPU usage."),
					resource.TestCheckResourceAttr("shoreline_bot."+pre+"_cpu_bot", "enabled", "true"),
					resource.TestCheckResourceAttr("shoreline_bot."+pre+"_cpu_bot", "family", "custom"),
				),
			},
		},
	})
}

func getAccResourceBot(prefix string) string {
	return `
		resource "shoreline_bot" "` + prefix + `_cpu_bot" {
			name        = "` + prefix + `_cpu_bot"
      command     = "if ${shoreline_alarm.` + prefix + `_cpu_alarm.name} then ${shoreline_action.` + prefix + `_ls_action.name}(\"/tmp\") fi"
      description = "Act on CPU usage."
      enabled     = true
			family      = "custom"
		}
`
}

//func getAccResourceBot(prefix string) string {
//	return `
//		resource "shoreline_bot" "` + prefix + `_cpu_bot" {
//			name = "` + prefix + `_cpu_bot"
//      alarm_statement = "${shoreline_alarm.` + prefix + `_cpu_alarm.name}"
//      action_statement = "${shoreline_action.` + prefix + `_ls_action.name}"
//      description = "Act on CPU usage."
//      enabled = true
//		}
//`
//}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
// Metric

func TestAccResourceMetric(t *testing.T) {
	pre := RandomAlphaPrefix(5)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfigString() + getAccResourceMetric(pre),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("shoreline_metric."+pre+"_cpu_plus_one", "name", pre+"_cpu_plus_one"),
					resource.TestCheckResourceAttr("shoreline_metric."+pre+"_cpu_plus_one", "value", "cpu_usage + 2"),
					resource.TestCheckResourceAttr("shoreline_metric."+pre+"_cpu_plus_one", "description", "Erroneous CPU usage."),
					resource.TestCheckResourceAttr("shoreline_metric."+pre+"_cpu_plus_one", "resource_query", "host | pod"),
					resource.TestCheckResourceAttr("shoreline_metric."+pre+"_cpu_plus_one", "units", "cores"),
					resource.TestCheckResourceAttr("shoreline_metric."+pre+"_cpu_plus_one", "shell", "/bin/bash"),
					resource.TestCheckResourceAttr("shoreline_metric."+pre+"_cpu_plus_one", "timeout", "5"),
				),
			},
		},
	})
}

func getAccResourceMetric(prefix string) string {
	return `
		resource "shoreline_metric" "` + prefix + `_cpu_plus_one" {
			name = "` + prefix + `_cpu_plus_one"
      value = "cpu_usage + 2"
      description = "Erroneous CPU usage."
      resource_query = "host | pod"
			units = "cores"
			shell = "/bin/bash"
			timeout = "5"
		}
`
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
// Resource

func TestAccResourceResource(t *testing.T) {
	pre := RandomAlphaPrefix(5)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfigString() + getAccResourceResource(pre),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("shoreline_resource."+pre+"_books", "name", pre+"_books"),
					resource.TestCheckResourceAttr("shoreline_resource."+pre+"_books", "description", "Pods with books app."),
					resource.TestCheckResourceAttr("shoreline_resource."+pre+"_books", "value", "host | pod | app = 'bookstore'"),
					resource.TestCheckResourceAttr("shoreline_resource."+pre+"_books", "shell", "/bin/bash"),
					resource.TestCheckResourceAttr("shoreline_resource."+pre+"_books", "timeout", "5"),
					resource.TestCheckResourceAttr("shoreline_resource."+pre+"_books", "res_env_var", "FOO"),
				),
			},
		},
	})
}

func getAccResourceResource(prefix string) string {
	return `
		resource "shoreline_resource" "` + prefix + `_books" {
			name = "` + prefix + `_books"
      description = "Pods with books app."
      value = "host | pod | app='bookstore'"
			shell = "/bin/bash"
			res_env_var = "FOO"
			timeout = "5"
		}
`
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
// File

func TestAccResourceFile(t *testing.T) {
	pre := RandomAlphaPrefix(5)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: getProviderConfigString() + getAccResourceFile(pre),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("shoreline_file."+pre+"_ex_file", "name", pre+"_ex_file"),
					resource.TestCheckResourceAttr("shoreline_file."+pre+"_ex_file", "destination_path", "/tmp/opcp_example.sh"),
					resource.TestCheckResourceAttr("shoreline_file."+pre+"_ex_file", "description", "op_copy example script."),
					resource.TestCheckResourceAttr("shoreline_file."+pre+"_ex_file", "resource_query", "host"),
					resource.TestCheckResourceAttr("shoreline_file."+pre+"_ex_file", "enabled", "false"),
					// computed values...
					resource.TestCheckResourceAttr("shoreline_file."+pre+"_ex_file", "file_length", "58"),
					resource.TestCheckResourceAttr("shoreline_file."+pre+"_ex_file", "checksum", "dbfb2a7d8176bd6e3dde256824421de3"),
					// just check that it's set
					resource.TestCheckResourceAttrSet("shoreline_file."+pre+"_ex_file", "file_length"),
				),
			},
		},
	})
}

func getAccResourceFile(prefix string) string {
	return `
		resource "shoreline_file" "` + prefix + `_ex_file" {
			name = "` + prefix + `_ex_file"
      input_file = "${path.module}/../data/opcp_example.sh"
      destination_path = "/tmp/opcp_example.sh"
      resource_query = "host"
      description = "op_copy example script."
      enabled = false
		}
`
}
