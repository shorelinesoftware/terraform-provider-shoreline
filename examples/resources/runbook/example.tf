# DEPRECATED: Use the `cells` field instead
resource "shoreline_runbook" "data_runbook" {
  name        = "data_runbook"
  description = "A sample runbook defined using the data field, which loads the runbook configuration from a separate JSON file."
  data        = file("${path.module}/data.json")
}


resource "shoreline_runbook" "full_runbook" {
  cells = jsonencode([
    {
      "md" : "CREATE"
    },
    {
      "op" : "action success = `echo SUCCESS`"
    },
    {
      "op" : "enable success"
    },
    {
      "op" : "success",
      "enabled" : false
    },
    {
      "md" : "CLEANUP"
    },
    {
      "op" : "delete success"
    }
  ])
  params = jsonencode([
    {
      "name" : "param_1",
      "value" : "<default_value>"
    },
    {
      "name" : "param_2",
      "value" : "<default_value>",
      "required" : false,
      "export" : true
    },
    {
      "name" : "param_3",
      "value" : "<default_value>",
      "export" : true
    },
    {
      "name" : "param_4",
      "required" : false
    }
  ])
  external_params = jsonencode([
    {
      "name" : "external_param_1",
      "source" : "alertmanager",
      "json_path" : "$.<path>",
      "export" : true,
      "value" : "<default_value>"
    },
    {
      "name" : "external_param_2",
      "source" : "alertmanager",
      "json_path" : "$.<path>",
      "value" : "<default_value>"
    },
    {
      "name" : "external_param_3",
      "source" : "alertmanager",
      "json_path" : "$.<path>",
      "export" : true
    },
    {
      "name" : "external_param_4",
      "source" : "alertmanager",
      "json_path" : "$.<path>"
    }
  ])
  name                                  = "full_runbook"
  description                           = "A sample runbook."
  timeout_ms                            = 5000
  allowed_entities                      = ["<user_1>", "<user_2>"]
  approvers                             = ["<user_2>", "<user_3>"]
  editors                               = ["<user_2>", "<user_4>"]
  is_run_output_persisted               = true
  allowed_resources_query               = "host"
  communication_workspace               = "<workspace_name>"
  communication_channel                 = "<channel_name>"
  labels                                = ["label1", "label2"]
  communication_cud_notifications       = true
  communication_approval_notifications  = false
  communication_execution_notifications = true
  filter_resource_to_action             = true
  enabled                               = true
  secret_names                          = ["<secret_name_1>", "<secret_name_2>"]
}


resource "shoreline_runbook" "minimal_runbook" {
  name  = "minimal_runbook"
  cells = jsonencode([])
}
