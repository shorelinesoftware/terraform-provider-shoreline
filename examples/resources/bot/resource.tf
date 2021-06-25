
resource "shoreline_bot" "cpu_bot" {
  name = "cpu_bot"
  command = "if ${shoreline_alarm.cpu_alarm.name} then ${shoreline_action.ls_action.name}('/tmp') fi"
  description = "Act on CPU usage."
  enabled = true
}

