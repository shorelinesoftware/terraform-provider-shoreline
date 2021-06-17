
terraform {
  required_providers {
    shoreline = {
      source  = "shoreline.io/terraform/shoreline"
      version = ">= 1.0"
    }
  }
}

provider "shoreline" {
  # provider configuration here
  url = "https://test.us.api.shoreline-vm1.io"
  #url = "https://test.us.api.shoreline-test5.io"
  #url = "https://test.us.api.shoreline-test4.io"
  #token = "xyz1.asdfj.asd3fas..."
}

resource "shoreline_bot" "cpu_bot" {
  name = "cpu_bot"
  alarm_statement = "${shoreline_alarm.cpu_alarm.name}"
  action_statement = "${shoreline_action.ls_action.name}"
  description = "Act on CPU usage."
  enabled = true
}

resource "shoreline_action" "ls_action" {
  name = "ls_action"
  command = "`ls /tmp`"
  description = "List some files"
  resource_query = "host"
  #params = [ "foo", "bar", "blah" ]
  #timeout = 60
  start_title_template    = "JVM dump started"
  complete_title_template = "JVM dump completed"
  error_title_template    = "JVM dump failed"

  enabled = true
}

resource "shoreline_alarm" "cpu_alarm" {
  name = "cpu_alarm"
  fire_query = "(cpu_usage > 0 | sum(5)) >= 2.75"
  clear_query = "(cpu_usage < 0 | sum(5)) >= 2.75"
  description = "Watch CPU usage."
  resource_query = "host"
  enabled = true
  resolve_short_template = "cleared blah123"
}

resource "shoreline_metric" "cpu_plus_one" {
  name = "cpu_plus"
  value = "cpu_usage + 2"
  description = "Erroneous CPU usage."
  resource_query = "host| pod"
}

resource "shoreline_resource" "books" {
  name = "books"
  description = "Pods with books app."
  value = "host | pod | app='bookstore'"
}

resource "shoreline_file" "ex_file" {
  name = "ex_file"
  input_file = "${path.module}/../data/opcp_example.sh"
  destination_path = "/tmp/opcp_example.sh"
  resource_query = "host"
  description = "op_copy example script."
  enabled = false
}


