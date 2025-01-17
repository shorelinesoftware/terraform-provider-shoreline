
resource "shoreline_bot" "full_bot" {
  name                    = "full_bot"
  command                 = "if ${var.full_alarm_name} then ${var.full_action_name}('/tmp') fi"
  description             = "Act on CPU usage."
  event_type              = "shoreline"
  monitor_id              = "<external_trigger_id>"
  alarm_resource_query    = "host"
  communication_workspace = "<communication_workspace>"
  communication_channel   = "<communication_channel>"
  integration_name        = "<integration_name>"
  enabled                 = true
}


resource "shoreline_bot" "minimal_bot" {
  name    = "minimal_bot"
  command = "if ${var.minimal_alarm_name} then ${var.minimal_action_name} fi"
}

