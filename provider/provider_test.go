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
				Config: getAccResourceAction(pre),
				Check: resource.ComposeTestCheckFunc(
					//resource.TestMatchResourceAttr( "shoreline_action.ls_action", "name", regexp.MustCompile("^ba")),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "name", pre+"_ls_action"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "resource_query", "host"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "command", "`ls /tmp`"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "timeout", "20"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "enabled", "true"),
					//resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "params", "[\"foo\"]"),
					resource.TestCheckResourceAttr("shoreline_action."+pre+"_ls_action", "error_title_template", "JVM dump failed"),
				),
			},
		},
	})
}

func getAccResourceAction(prefix string) string {
	return `
		resource "shoreline_action" "` + prefix + `_ls_action" {
			name = "` + prefix + `_ls_action"
			command = "` + "`ls /tmp`" + `"
			description = "List some files"
			resource_query = "host"
			timeout = 20
			#params = ["foo"]
			start_title_template    = "JVM dump started"
			complete_title_template = "JVM dump completed"
			error_title_template    = "JVM dump failed"
			enabled = true
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
				Config: getAccResourceAlarm(pre),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "name", pre+"_cpu_alarm"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "resource_query", "host"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "fire_query", "( cpu_usage > 0 | sum ( 5 ) ) >= 2"),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "description", "Watch CPU usage."),
					resource.TestCheckResourceAttr("shoreline_alarm."+pre+"_cpu_alarm", "enabled", "true"),
				),
			},
		},
	})
}

func getAccResourceAlarm(prefix string) string {
	return `
		resource "shoreline_alarm" "` + prefix + `_cpu_alarm" {
			name = "` + prefix + `_cpu_alarm"
	    fire_query = "(cpu_usage > 0 | sum(5)) >= 2"
	    clear_query = "(cpu_usage < 0 | sum(5)) >= 2"
	    description = "Watch CPU usage."
	    resource_query = "host"
	    enabled = true
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
				Config: getAccResourceAction(pre) + getAccResourceAlarm(pre) + getAccResourceBot(pre),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("shoreline_bot."+pre+"_cpu_bot", "name", pre+"_cpu_bot"),
					//resource.TestCheckResourceAttr("shoreline_bot."+pre+"_cpu_bot", "alarm_statement", pre+"_cpu_alarm"),
					resource.TestCheckResourceAttr("shoreline_bot."+pre+"_cpu_bot", "command", "if "+pre+"_cpu_alarm then "+pre+"_ls_action fi"),
					resource.TestCheckResourceAttr("shoreline_bot."+pre+"_cpu_bot", "description", "Act on CPU usage."),
					resource.TestCheckResourceAttr("shoreline_bot."+pre+"_cpu_bot", "enabled", "true"),
				),
			},
		},
	})
}

func getAccResourceBot(prefix string) string {
	return `
		resource "shoreline_bot" "` + prefix + `_cpu_bot" {
			name = "` + prefix + `_cpu_bot"
      command = "if ${shoreline_alarm.` + prefix + `_cpu_alarm.name} then ${shoreline_action.` + prefix + `_ls_action.name} fi"
      description = "Act on CPU usage."
      enabled = true
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
				Config: getAccResourceMetric(pre),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("shoreline_metric."+pre+"_cpu_plus_one", "name", pre+"_cpu_plus_one"),
					resource.TestCheckResourceAttr("shoreline_metric."+pre+"_cpu_plus_one", "value", "cpu_usage + 2"),
					resource.TestCheckResourceAttr("shoreline_metric."+pre+"_cpu_plus_one", "description", "Erroneous CPU usage."),
					//resource.TestCheckResourceAttr("shoreline_metric."+pre+"_cpu_plus_one", "resource_query", "host | pod"),
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
      #resource_query = "host| pod"
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
				Config: getAccResourceResource(pre),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("shoreline_resource."+pre+"_books", "name", pre+"_books"),
					resource.TestCheckResourceAttr("shoreline_resource."+pre+"_books", "description", "Pods with books app."),
					resource.TestCheckResourceAttr("shoreline_resource."+pre+"_books", "value", "host | pod | app = 'bookstore'"),
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
				Config: getAccResourceFile(pre),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("shoreline_file."+pre+"_ex_file", "name", pre+"_ex_file"),
					resource.TestCheckResourceAttr("shoreline_file."+pre+"_ex_file", "destination_path", "/tmp/opcp_example.sh"),
					resource.TestCheckResourceAttr("shoreline_file."+pre+"_ex_file", "description", "op_copy example script."),
					resource.TestCheckResourceAttr("shoreline_file."+pre+"_ex_file", "resource_query", "host"),
					resource.TestCheckResourceAttr("shoreline_file."+pre+"_ex_file", "enabled", "false"),
					// TODO length and checksum
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
