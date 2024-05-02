
resource "shoreline_circuit_breaker" "delayed_circuit_breaker" {
  name       = "delayed_circuit_breaker"
  command    = shoreline_action.delayed_action.name
  soft_limit = 1
  enabled    = true
}