
resource "shoreline_file" "ex_file" {
  name             = "ex_file"
  input_file       = "${path.module}/../data/opcp_example.sh"
  destination_path = "/tmp/opcp_example.sh"
  resource_query   = "host"
  description      = "op_copy example script."
  enabled          = false
}

