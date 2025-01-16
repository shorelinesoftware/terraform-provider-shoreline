
resource "shoreline_alarm" "full_alarm" {
  name               = "full_alarm"
  fire_query         = "(cpu_usage > 1 | sum(5)) >= 2.75"
  mute_query         = "(cpu_usage < 0 | sum(5)) >= 2.75"
  clear_query        = "(cpu_usage < 0 | sum(5)) >= 2.75"
  description        = "Watch CPU usage."
  resource_query     = "host"
  check_interval_sec = 50
  compile_eligible   = false
  condition_type     = "above"
  condition_value    = "10"
  metric_name        = "<metric_name>"
  raise_for          = "local"
  resource_type      = "HOST"

  fire_title_template    = "fired cpu_alarm title"
  resolve_title_template = "cleared cpu_alarm title"

  fire_short_template    = "fired cpu_alarm short"
  resolve_short_template = "cleared cpu_alarm short"

  fire_long_template    = "fired cpu_alarm long"
  resolve_long_template = "cleared cpu_alarm long"

  enabled = true
}


resource "shoreline_alarm" "minimal_alarm" {
  name               = "minimal_alarm"
  fire_query         = "(cpu_usage > 1 | sum(5)) >= 2.75"
}


resource "shoreline_alarm" "cpu_threshold_alarm" {
  fire_query = "${var.cpu_threshold_action_name}(cpu_threshold=75) == 1"
  name       = "cpu_threshold_alarm"

  clear_query        = "${var.cpu_threshold_action_name}(cpu_threshold=75) == 0"
  description        = "High CPU usage alarm"
  enabled            = true
  resource_query     = "hosts"
  check_interval_sec = 10

  fire_short_template    = "High CPU Alarm fired"
  resolve_short_template = "High CPU Alarm resolved"
}

