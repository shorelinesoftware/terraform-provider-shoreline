---
page_title: "shoreline_alarm Resource - terraform-provider-shoreline"
subcategory: ""
description: |-
  Alarms are fully-customizable Metric or status checks that automatically trigger remediation Actions.
---

# shoreline_alarm (Resource)

Alarms frequently check one or more [Metric](/t/metric) thresholds or custom [Resource](/t/resource) queries. The Alarm is raised based on custom thresholds or shell commands you define, which informs any connected [Bot](/t/bot) to trigger remedial [Actions](/t/action).

-> Shoreline includes dozens of built-in [Metrics](/t/metric) on which to base your Alarms. You can also combine multiple [Metrics](/t/metric) into [Metric Sets](/t/metric-set) for monitoring many [Metrics](/t/metric) at once. You can even create your own [Derived Metrics](/t/derived-metric) if none of the built-in options meet your needs.

## Required Properties

Each Alarm can define many [properties](/alarms/properties) to determine its behavior. The required properties when [creating an Alarm](#create-an-alarm) are:

- [name](/alarms/properties#name) - The name of the Alarm.
- [fire_query](/alarms/properties#fire_query) - The [Op](/t/op) statement that triggers the Alarm.
- [clear_query](/alarms/properties#clear_query) - The [Op](/t/op) statement that clears the Alarm.
- [resource_query](/alarms/properties#resource_query) - The [Op](/t/op) query that selects which [Resources](/t/resource) the Alarm triggers from.

## Usage

The following example creates an [Alarm](/t/alarm) named `my_cpu_alarm` that fires when at least 80% of a host [Resource's](/t/resource) CPU usage metric measurements are equal to or exceed `40%` over the previous minute.

```tf
resource "shoreline_alarm" "cpu_alarm" {
  name = "my_cpu_alarm"
  fire_query = "(cpu_usage > 40 | sum(60)) >= 48.0"
  clear_query = "(cpu_usage < 40 | sum(60)) >= 48.0"
  resource_query = "hosts"
}
```

-> [Metric](/t/metrics) data points are collected once per second for all [Shoreline Resources](/t/resource) (i.e. hosts, pods, and containers). Thus, a [Metric](/t/metric) query of `(cpu_usage > 40 | sum(60)) >= 48.0` determines if at least 48 of the last 60 `cpu_usage` data points exceeded `40%`.  You can learn more from the [Metrics documentation](/t/metric).

### Advanced Usage

You can also combine other Terraform resource blocks and variables to create complex [Alarms](/t/alarm).  In this example we're defining an [Action](/t/action) called `jvm_trace_check_heap` that determines if JVM heap usage exceeds a variable-defined threshold:

{{tffile "examples/op_packs/jvm_trace/actions.tf"}}

The `jvm_trace_heap_alarm` [Alarm](/t/alarm) executes the `jvm_trace_check_heap` [Action](/t/action) as part of its [fire_query](/t/alarm/property#fire_query) and [clear_query](/t/alarm/property#clear_query), to determine when the [Alarm](/t/alarm) is raised or resolved, respectively:

{{tffile "examples/op_packs/jvm_trace/alarms.tf"}}

-> See the Shoreline [Alarms Documentation](/t/alarm) for more info.

{{ .SchemaMarkdown | trimspace }}