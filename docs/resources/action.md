---
page_title: 'shoreline_action Resource - terraform-provider-shoreline'
subcategory: ''
description: |-
---

# shoreline_action (Resource)

Actions execute shell commands on associated [Resources](https://docs.shoreline.io/platform/resources). Whenever an [Alarm](https://docs.shoreline.io/alarms) fires the associated [Bot](https://docs.shoreline.io/bots) triggers the corresponding [Action](https://docs.shoreline.io/actions), closing the basic auto-remediation loop of
Shoreline.

## Required Properties

Each Action has many properties that determine its behavior. The required properties are:

- [name](https://docs.shoreline.io/actions/properties#name) - The name of the Action.
- [command](https://docs.shoreline.io/actions/properties#command) - The shell command executed when the Action triggers.

-> Check out [Action Properties](https://docs.shoreline.io/actions/properties) for details on all available properties and how to use them.

## Usage

The following [Action](https://docs.shoreline.io/actions) definition creates a `cpu_threshold_action` that compares host CPU usage against a `cpu_threshold` parameter value.

```terraform
resource "shoreline_action" "cpu_threshold_action" {
  name = "cpu_threshold_action"
  # Evaluates current CPU usage and compares it to a parameter value named $cpu_threshold
  command = "`if [ $[100-$(vmstat 1 2|tail -1|awk '{print $15}')] -gt $cpu_threshold ]; then exit 1; fi`"

  description    = "Check CPU usage"
  enabled        = true
  params         = ["cpu_threshold"]
  resource_query = "hosts"
  timeout        = 5000

  start_title_template    = "CPU threshold action started"
  complete_title_template = "CPU threshold action completed"
  error_title_template    = "CPU threshold action failed"

  start_short_template    = "CPU threshold action short started"
  complete_short_template = "CPU threshold action short completed"
  error_short_template    = "CPU threshold action short failed"
}
```

This Action can be executed via an [Alarm's](https://docs.shoreline.io/alarms) [clear_query](https://docs.shoreline.io/alarms/properties#clear_query)
/ [fire_query](https://docs.shoreline.io/alarms/properties#fire_query), or directly via an [Op](https://docs.shoreline.io/op) command.

For example, the following [Alarm](https://docs.shoreline.io/alarms) fires and clears based on the result of the previously-generated `cpu_threshold_action`:

```terraform
resource "shoreline_alarm" "cpu_threshold_alarm" {
  fire_query = "cpu_threshold_action(cpu_threshold=75) == 1"
  name       = "cpu_threshold_alarm"

  clear_query        = "cpu_threshold_action(cpu_threshold=75) == 0"
  description        = "High CPU usage alarm"
  enabled            = true
  resource_query     = "hosts"
  check_interval_sec = 10

  fire_short_template    = "High CPU Alarm fired"
  resolve_short_template = "High CPU Alarm resolved"
}
```

You can also define [Terraform Input Variables](https://www.terraform.io/docs/language/values/variables.html) and use them within your [Action](https://docs.shoreline.io/actions) definitions:

```terraform
variable "namespace" {
  type        = string
  description = "A namespace to isolate multiple instances of the module with different parameters."
}

variable "resource_query" {
  type        = string
  description = "The set of hosts/pods/containers monitored and affected by this module."
}

variable "jvm_process_regex" {
  type        = string
  description = "A regular expression to match and select the monitored Java processes."
}

variable "mem_threshold" {
  type        = number
  description = "The high-water-mark, in Mb, above which the JVM process stack-trace is dumped."
  default     = 2000
}

variable "check_interval" {
  type        = number
  description = "Frequency, in seconds, to check the memory usage."
  default     = 60
}

variable "script_path" {
  type        = string
  description = "Destination (on selected resources) for the check, and stack-dump scripts."
  default     = "/agent/scripts"
}

variable "s3_bucket" {
  type        = string
  description = "Destination in AWS S3 for stack-dump output files."
  default     = "shore-oppack-test"
}
```

```terraform
# Action to check the JVM heap usage on the selected resources and process.
resource "shoreline_action" "jvm_trace_check_heap" {
  name        = "${var.namespace}_jvm_check_heap"
  description = "Check heap utilization by process regex."
  # Parameters passed in: the regular expression to select process name.
  params = ["JVM_PROCESS_REGEX"]
  # Extract the heap used for the matching process and return 1 if above threshold.
  command = "`hm=$(jstat -gc $(jps | grep \"$${JVM_PROCESS_REGEX}\" | awk '{print $1}') | tail -n 1 | awk '{split($0,a,\" \"); sum=a[3]+a[4]+a[6]+a[8]; print sum/1024}'); hm=$${hm%.*}; if [ $hm -gt ${var.mem_threshold} ]; then echo \"heap memory $hm MB > threshold ${var.mem_threshold} MB\"; exit 1; fi`"

  # UI / CLI annotation informational messages:
  start_short_template    = "Checking JVM heap usage."
  error_short_template    = "Error checking JVM heap usage."
  complete_short_template = "Finished checking JVM heap usage."
  start_long_template     = "Checking JVM process ${var.jvm_process_regex} heap usage."
  error_long_template     = "Error checking JVM process ${var.jvm_process_regex} heap usage."
  complete_long_template  = "Finished checking JVM process ${var.jvm_process_regex} heap usage."

  enabled = true
}
```

-> See the Shoreline [Actions Documentation](https://docs.shoreline.io/actions) for more info.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `command` (String) A specific action to run.
- `name` (String) The name/symbol for the object within Shoreline and the op language (must be unique, only alphanumeric/underscore).

### Optional

- `allowed_entities` (List of String) The list of users who can run an action or notebook. Any user can run if left empty.
- `allowed_resources_query` (String) The list of resources on which an action or notebook can run. No restriction, if left empty. Defaults to ``.
- `communication_channel` (String) A string value denoting the slack channel where notifications related to the object should be sent to. Defaults to ``.
- `communication_workspace` (String) A string value denoting the slack workspace where notifications related to the object should be sent to. Defaults to ``.
- `complete_long_template` (String) The long description of the Action's completion. Defaults to ``.
- `complete_short_template` (String) The short description of the Action's completion. Defaults to ``.
- `complete_title_template` (String) UI title of the Action's completion. Defaults to ``.
- `description` (String) A user-friendly explanation of an object. Defaults to ``.
- `editors` (List of String) List of users who can edit the object (with configure permission). Empty maps to all users.
- `enabled` (Boolean) If the object is currently enabled or disabled. Defaults to `false`.
- `error_long_template` (String) The long description of the Action's error condition. Defaults to ``.
- `error_short_template` (String) The short description of the Action's error condition. Defaults to ``.
- `error_title_template` (String) UI title of the Action's error condition. Defaults to ``.
- `file_deps` (List of String) file object dependencies.
- `params` (List of String) Named variables to pass to an object (e.g. an Action).
- `res_env_var` (String) Result environment variable ... an environment variable used to output values through. Defaults to ``.
- `resource_query` (String) A set of Resources (e.g. host, pod, container), optionally filtered on tags or dynamic conditions. Defaults to ``.
- `resource_tags_to_export` (List of String)
- `shell` (String) The commandline shell to use (e.g. /bin/sh). Defaults to ``.
- `start_long_template` (String) The long description when starting the Action. Defaults to ``.
- `start_short_template` (String) The short description when starting the Action. Defaults to ``.
- `start_title_template` (String) UI title of the start of the Action. Defaults to ``.
- `timeout` (Number) Maximum time to wait, in milliseconds. Defaults to `60000`.

### Read-Only

- `id` (String) The ID of this resource.
- `type` (String) The type of object (i.e., Alarm, Action, Bot, Metric, Resource, File, Time Trigger, Integration, Runbook, Principal, Report Template, Dashboard, NVault Secret).
