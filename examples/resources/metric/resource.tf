
resource "shoreline_metric" "cpu_plus_one" {
  name = "cpu_plus"
  value = "cpu_usage + 2"
  description = "Erroneous CPU usage."
  resource_query = "host| pod"
}

