---
page_title: 'shoreline_action Resource - terraform-provider-shoreline'
subcategory: ''
description: |-
---

# shoreline_action (Resource)

Actions execute shell commands on associated [Resources](/t/resource). Whenever an [Alarm](/t/alarm) fires the associated [Bot](/t/bot) triggers the corresponding [Action](/t/action), closing the basic auto-remediation loop of
Shoreline.

## Required Properties

Each Action has many properties that determine its behavior. The required properties are:

- [name](/t/action/property#name) - The name of the Action.
- [command](/t/action/property#command) - The shell command executed when the Action triggers.

-> Check out [Action Properties](/t/action/property) for details on all available properties and how to use them.

## Usage

The following [Action](/t/action) definition creates a `cpu_threshold_action` that compares host CPU usage against a `cpu_threshold` parameter value.

{{tffile "examples/resources/action/cpu_threshold_action.tf"}}

This Action can be executed via an [Alarm's](/t/alarm) [clear_query](/t/alarm/property#clear_query)
/ [fire_query](/t/alarm/property#fire_query), or directly via an [Op](/t/op) command.

For example, the following [Alarm](/t/alarm) fires and clears based on the result of the previously-generated `cpu_threshold_action`:

{{tffile "examples/resources/alarm/cpu_threshold_alarm.tf"}}

You can also define [Terraform Input Variables](https://www.terraform.io/docs/language/values/variables.html) and use them within your [Action](/t/action) definitions:

{{tffile "examples/op_packs/jvm_trace/variables.tf"}}

{{tffile "examples/op_packs/jvm_trace/actions.tf"}}

-> See the Shoreline [Actions Documentation](https://docs.shoreline.io/actions) for more info.

{{ .SchemaMarkdown | trimspace }}
