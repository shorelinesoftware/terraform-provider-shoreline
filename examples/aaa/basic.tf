# Defines the required Shoreline provider and version
# @ref https://docs.shoreline.io/op/packs/tutorial#create-a-configuration-file
terraform {
  required_providers {
    shoreline = {
      source  = "shorelinesoftware/shoreline"
      version = ">= 1.0.7"
    }
  }
}

# Set provider-specific arguments
# @ref https://docs.shoreline.io/op/packs/tutorial#create-a-configuration-file
provider "shoreline" {
  # Set to the Shoreline cluster API URL, e.g. https://acme.us.api.shoreline-acme.io
  url = "<CLUSTER_API_ENDPOINT>"
}

# Create an Alarm that fires when host CPU usage exceeds 35% for 48 of the previous 60 seconds.
# This Alarm clears when CPU usage is below 35% for the previous 180 seconds.
# @ref https://docs.shoreline.io/op/packs/tutorial#create-an-alarm
# @ref https://docs.shoreline.io/alarms
resource "shoreline_alarm" "high_cpu_alarm" {
  name                   = "high_cpu_alarm"
  fire_query             = "(cpu_usage > 35 | sum(60)) >= 48"
  clear_query            = "(cpu_usage < 35 | sum(180)) >= 180"
  description            = "Watch CPU usage."
  resource_query         = "hosts"
  enabled                = true
  resolve_short_template = "high_cpu_alarm resolved"
}

# Create an Action that executes a Linux command to count the active background jobs.
# @ref https://docs.shoreline.io/op/packs/tutorial#create-an-action
# @ref https://docs.shoreline.io/actions
resource "shoreline_action" "background_jobs_action" {
  name                    = "background_jobs_action"
  command                 = "`top -b -n 1 | head -n 15`"
  description             = "Count background jobs"
  start_title_template    = "background_jobs_action started"
  complete_title_template = "background_jobs_action completed"
  error_title_template    = "background_jobs_action failed"
  enabled                 = true
}

# Create a Bot that triggers the 'background_jobs_action' when the 'high_cpu_alarm' fires.
# @ref https://docs.shoreline.io/op/packs/tutorial#create-a-bot
# @ref https://docs.shoreline.io/bots
resource "shoreline_bot" "cpu_bot" {
  name        = "cpu_bot"
  command     = "if ${shoreline_alarm.high_cpu_alarm.name} then ${shoreline_action.background_jobs_action.name} fi"
  description = "Count background jobs on high CPU usage."
  enabled     = true
}
