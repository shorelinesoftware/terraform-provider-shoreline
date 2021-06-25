
resource "shoreline_alarm" "cpu_alarm" {
  name           = "cpu_alarm"
  fire_query     = "(cpu_usage > 0 | sum(5)) >= 2.75"
  clear_query    = "(cpu_usage < 0 | sum(5)) >= 2.75"
  description    = "Watch CPU usage."
  resource_query = "host"

  fire_short_template    = "fired blah123"
  resolve_short_template = "cleared blah123"
  check_interval         = 50
  compile_eligible       = false
  condition_type         = "above"
  condition_value        = "10"

  enabled = true
}

