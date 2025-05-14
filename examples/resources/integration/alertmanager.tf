
resource "shoreline_integration" "am0" {
  name             = "am0"
  serial_number    = "001"
  service_name     = "alertmanager"
  external_url     = "http://alertmanager-0-external-url.io"
  payload_paths    = ["$.alertmanager.payload.path.0"]
  enabled          = true
  permissions_user = "<permissions_user>"
}

resource "shoreline_integration" "am1" {
  name             = "am1"
  serial_number    = "002"
  service_name     = "alertmanager"
  external_url     = "http://alertmanager-1-external-url.io"
  payload_paths    = ["$.alertmanager.payload.path.1"]
  enabled          = true
  permissions_user = "<permissions_user>"
}

resource "shoreline_integration" "am2" {
  name             = "am2"
  serial_number    = "003"
  service_name     = "alertmanager"
  external_url     = "http://alertmanager-2-external-url.io"
  payload_paths    = ["$.alertmanager.payload.path.2"]
  enabled          = true
  permissions_user = "<permissions_user>"
}

resource "shoreline_integration" "am3" {
  name             = "am3"
  serial_number    = "004"
  service_name     = "alertmanager"
  external_url     = "http://alertmanager-3-external-url.io"
  payload_paths    = ["$.alertmanager.payload.path.3"]
  enabled          = true
  permissions_user = "<permissions_user>"
}

resource "shoreline_integration" "am4" {
  name             = "am4"
  serial_number    = "005"
  service_name     = "alertmanager"
  external_url     = "http://alertmanager-4-external-url.io"
  payload_paths    = ["$.alertmanager.payload.path.4"]
  enabled          = true
  permissions_user = "<permissions_user>"
}

resource "shoreline_integration" "am5" {
  name             = "am5"
  serial_number    = "006"
  service_name     = "alertmanager"
  external_url     = "http://alertmanager-5-external-url.io"
  payload_paths    = ["$.alertmanager.payload.path.5"]
  enabled          = true
  permissions_user = "<permissions_user>"
}

resource "shoreline_integration" "am6" {
  name             = "am6"
  serial_number    = "007"
  service_name     = "alertmanager"
  external_url     = "http://alertmanager-6-external-url.io"
  payload_paths    = ["$.alertmanager.payload.path.6"]
  enabled          = true
  permissions_user = "<permissions_user>"
}

resource "shoreline_integration" "am7" {
  name             = "am7"
  serial_number    = "008"
  service_name     = "alertmanager"
  external_url     = "http://alertmanager-7-external-url.io"
  payload_paths    = ["$.alertmanager.payload.path.7"]
  enabled          = true
  permissions_user = "<permissions_user>"
}

resource "shoreline_integration" "am8" {
  name             = "am8"
  serial_number    = "009"
  service_name     = "alertmanager"
  external_url     = "http://alertmanager-8-external-url.io"
  payload_paths    = ["$.alertmanager.payload.path.8"]
  enabled          = true
  permissions_user = "<permissions_user>"
}

resource "shoreline_integration" "am9" {
  name             = "am9"
  serial_number    = "010"
  service_name     = "alertmanager"
  external_url     = "http://alertmanager-9-external-url.io"
  payload_paths    = ["$.alertmanager.payload.path.9"]
  enabled          = true
  permissions_user = "<permissions_user>"
}
