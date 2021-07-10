
terraform {
  required_providers {
    shoreline = {
      source  = "shoreline.io/terraform/shoreline"
      version = ">= 1.0.4"
    }
  }
}

provider "shoreline" {
  # provider configuration here
  #token = "xyz1.asdfj.asd3fas..."
  url = "https://test.us.api.shoreline-vm1.io"
  #url = "https://test.us.api.shoreline-test6.io"
  retries = 2
  debug = true
}

resource "shoreline_bot" "cpu_bot" {
  name = "cpu_bot"
  command = "if ${shoreline_alarm.cpu_alarm.name} then ${shoreline_action.ls_action.name}(dir='/tmp') fi"
  #command = "if ${shoreline_alarm.cpu_alarm.name} then ${shoreline_action.ls_action.name}('/tmp', 'blah') fi"
  description = "Act on CPU usage."
  enabled = true
}

resource "shoreline_action" "ls_action" {
  name = "ls_action"
  description = "List some files ..."

  command = "`ls $${dir}; export FOO='bar'`"
  params = [ "dir" ]
  #command = "`ls $${dir}; echo $${msg}; export FOO='bar'`"
  #params = [ "dir" , "msg" ]

  res_env_var = "FOO"
  resource_query = "host"
  #timeout = 60
  start_title_template    = "JVM dump started"
  complete_title_template = "JVM dump  completed"
  error_title_template    = "JVM dump failed"

  start_short_template    = "JVM dump short started"
  complete_short_template = "JVM dump short completed"
  error_short_template    = "JVM dump short failed"

  enabled = true
}

resource "shoreline_alarm" "cpu_alarm" {
  name = "cpu_alarm"
  fire_query = "(cpu_usage > 0 | sum(5)) >= 2.75"
  clear_query = "(cpu_usage < 0 | sum(5)) >= 2.75"
  description = "Watch CPU usage."
  resource_query = "host"

  fire_short_template = "fired blah123"
  resolve_short_template = "cleared blah123"
  check_interval_sec = 50
  compile_eligible = false
  condition_type = "above"
  condition_value = 10.1

  enabled = true
}

resource "shoreline_metric" "cpu_plus_one" {
  name = "cpu_plus"
  value = "cpu_usage + 2"
  description = "Erroneous CPU usage."
  #resource_query = "host| pod"
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


