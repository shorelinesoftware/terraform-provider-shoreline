
resource "shoreline_time_trigger" "tt" {
  name       = "tt"
  fire_query = "every 5m"
  start_date = "2024-02-29T08:00:00"
  end_date   = "2026-02-28T08:00:00"
  enabled    = true
}

