---
page_title: "shoreline_notebook Resource - terraform-provider-shoreline"
subcategory: ""
description: |- Notebooks replace static runbooks by capturing interactive debug and remediation sessions in a convenient UI.
---

# shoreline_notebook (Resource)

Through Shoreline's web-based UI, [Notebooks](https://docs.shoreline.io/notebooks) automatically capture an entire debug and remediation session -- which can optionally be associated with a specific [Alarm](https://docs.shoreline.io/alarms) -- and then be shared with other team members to streamline future incident response. With Notebooks you can:

- Create a series of interactive [Op statement cells](https://docs.shoreline.io/notebooks#op-statements) allowing you to execute [Op commands](https://docs.shoreline.io/op/commands) within your browser -- all without installing or configuring the local [CLI](https://docs.shoreline.io/cli).
- Define and use [dynamic parameters](https://docs.shoreline.io/notebooks/parameters) across Notebook Op cells.
- Memorialize Notebooks with [historical snapshots](https://docs.shoreline.io/notebooks/runs).
- Add [Markdown-based notes](https://docs.shoreline.io/notebooks#notes) to inform operators how to use the Notebook.
- [Associate](https://docs.shoreline.io/notebooks#alarm-association) existing [Alarms](https://docs.shoreline.io/alarms) and Notebooks, allowing on-call members to click through to an interactive debugging and remediation Notebook directly from the triggered [Alarm](https://docs.shoreline.io/alarms) UI.

## Required Properties

Each Notebook uses a variety of properties to determine its behavior. The required properties when [creating a Notebook](https://docs.shoreline.io/notebooks#create-a-notebook) are:

- `name`: string - A unique symbol name for the Notebook object.
- `cells`: list(object) - A list of cells represented by JSON objects. Cells may either be [Op statement cells](https://docs.shoreline.io/notebooks#op-statements) or [Markdown cells](https://docs.shoreline.io/notebooks#notes).

### Download a Notebook as a Terraform resource

You can download an entire Notebook directly as a Terraform resource. This will allow you to just plug in the TF code into your infrastructure and deploy the runbook immediately.

1. Click the **Actions** button on the right side of the active Notebook panel.
2. Select the **Download Notebook as Terraform** button to download the full configuration of the current Notebook as a Terraform resource.

## Defining a Notebook using the legacy `data` property

You can also export the Notebook's configuration as a JSON file and then freely modify, share, and [upload](https://docs.shoreline.io/notebooks#upload-a-notebook) this Notebook at any time.

Note: this way of defining is **deprecated**. Please refer to the above instructions using the new format.

The following example creates a [Notebook](https://docs.shoreline.io/notebooks) named `my_notebook`.

1. [Download a Notebook](https://docs.shoreline.io/notebooks#download-a-notebook) as JSON.
2. Only keep the `cells`, `params`, `external_params`, and `enabled` fields fron the JSON file. Note: `externalParams` needs to be renamed to `external_params` in the JSON file.
3. Save the Notebook JSON to local path within your Terraform project.
4. Define a new `shoreline_notebook` Terraform resource in your Terraform configuration that points the `data` property to the correct local module path.

   ```terraform
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
```

-> See the Shoreline [Notebooks Documentation](https://docs.shoreline.io/notebooks) for more info on creating and using Notebooks.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name/symbol for the object within Shoreline and the op language (must be unique, only alphanumeric/underscore).

### Optional

- `allowed_entities` (List of String) The list of users who can run an action or notebook. Any user can run if left empty.
- `allowed_resources_query` (String) The list of resources on which an action or notebook can run. No restriction, if left empty. Defaults to ``.
- `approvers` (List of String)
- `cells` (String) The data cells inside a notebook. Defined as a list of JSON objects. These may be either Markdown or Op commands. Defaults to ``.
- `communication_approval_notifications` (Boolean) Enables slack notifications for approvals operations. (Requires workspace and channel.) Defaults to `true`.
- `communication_channel` (String) A string value denoting the slack channel where notifications related to the object should be sent to. Defaults to ``.
- `communication_cud_notifications` (Boolean) Enables slack notifications for create/update/delete operations. (Requires workspace and channel.) Defaults to `true`.
- `communication_execution_notifications` (Boolean) Enables slack notifications for the object executions. (Requires workspace and channel.) Defaults to `true`.
- `communication_workspace` (String) A string value denoting the slack workspace where notifications related to the object should be sent to. Defaults to ``.
- `data` (String, Deprecated) **Deprecated** Field 'data' is obsolete. The JSON representation of a Notebook. If this field is used, then the JSON should only contain these four fields: cells, params, external_params and enabled. Defaults to ``.
- `description` (String) A user-friendly explanation of an object. Defaults to ``.
- `editors` (List of String) List of users who can edit the object (with configure permission). Empty maps to all users.
- `enabled` (Boolean) If the object is currently enabled or disabled. Defaults to `true`.
- `external_params` (String) Notebook parameters defined via with a JSON path used to extract the parameter's value from an external payload, such as an Alertmanager alert. Defaults to ``.
- `filter_resource_to_action` (Boolean) Determines whether parameters containing resources are exported to actions. Defaults to `false`.
- `is_run_output_persisted` (Boolean) A boolean value denoting whether or not cell outputs should be persisted when running a notebook Defaults to `true`.
- `labels` (List of String) A list of strings by which notebooks can be grouped.
- `params` (String) Named variables to pass to an object (e.g. an Action). Defaults to ``.
- `resource_query` (String, Deprecated) **Deprecated** Please use 'allowed_resources_query' instead. A set of Resources (e.g. host, pod, container), optionally filtered on tags or dynamic conditions. Defaults to ``.
- `secret_names` (List of String) A list of strings that contains the name of the secrets that are used in the runbook.
- `timeout_ms` (Number) Defaults to `60000`.

### Read-Only

- `id` (String) The ID of this resource.
- `type` (String) The type of object (i.e., Alarm, Action, Bot, Metric, Resource, or File).
