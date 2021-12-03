---
page_title: "shoreline_bot Resource - terraform-provider-shoreline"
subcategory: ""
description: |-
  Alarms use Bots to execute automated Actions.
---

# shoreline_bot (Resource)

A Bot connects a single [Alarm](/t/alarm) to one or more [Actions](/t/action). When the [Alarm](/t/alarm) is raised the [Bot](/t/bot) fires all associated and enabled [Actions](/t/action) to close the auto-remediation loop.

## Required Properties

Each Bot has various configurable [properties](/bots/properties) that determine its behavior. The minimal required properties to [create a Bot](/bots#create-a-bot) are:

- [name](/bots/properties/name) - The name of the Bot
- [command](/bots/properties/command) - An `if-then-fi` statement containing the [Alarm](/t/alarm) name and [Action](/t/action) name associated with the Bot. Alternatively, the `command` property can be a custom Linux command.

## Usage

The following example creates a [Bot](/t/bot) named `cpu_bot` that executes the `restart_action` [Action](/t/action) when the `high_cpu_alarm` [Alarm's](/t/alarm) [fire_query](/alarms/properties#fire_query) is true:

```tf
resource "shoreline_bot" "cpu_bot" {
  name = "cpu_bot"
  command = "if ${shoreline_alarm.high_cpu_alarm.name} then ${shoreline_action.restart_action.name} fi"
  description = "Restart on high CPU usage."
  enabled = true
}
```

The `command` property specifies the [Alarm](/t/alarm) and [Action](/t/action) that are connected by this [Bot](/t/bot). It uses Terraform's [built-in string interpolation](https://www.terraform.io/docs/language/expressions/strings.html#interpolation) to evaluate the name of both the [Alarm](/t/alarm) and [Action](/t/action).

### Advanced Usage

Configuring a combination of an [Alarm](/t/alarm), [Action](/t/action), and [Bot](/t/bot) closes the fundamental auto-remediation loop provided by Shoreline.  Below we're using portions of Shoreline's JVM [Op Pack](/t/op/pack) to create a full incident automation loop when JVM memory usage gets too high.

First, the `jvm_trace_check_heap` [Action](/t/action) determines if JVM heap usage exceeds a variable-defined threshold:

{{tffile "examples/op_packs/jvm_trace/actions.tf"}}

The `jvm_trace_heap_alarm` [Alarm](/t/alarm) executes the `jvm_trace_check_heap` [Action](/t/action) as part of its [fire_query](/t/alarm/property#fire_query) and [clear_query](/t/alarm/property#clear_query):

{{tffile "examples/op_packs/jvm_trace/alarms.tf"}}

We define another [Action](/t/action) called `jvm_trace_jvm_debug` that executes a bash script that dumps JVM debug data to AWS S3 before restarting the JVM:

```tf
# Action to dump the JVM stack-trace on the selected resources and process.
resource "shoreline_action" "jvm_trace_jvm_debug" {
  name = "${var.namespace}_jvm_dump_stack"
  description = "Dump JVM process (by regex) heap, thread and GC info to s3, then kill the pod."
  # Parameters passed in: the regular expression to select process name, and destination AWS S3 bucket.
  params = [ "JVM_PROCESS_REGEX" , "S3_BUCKET"]
  # Extract process info, and kill the pod.
  command = "`cd ${var.script_path} && chmod +x jvm_dumps.sh && ./jvm_dumps.sh $${JVM_PROCESS_REGEX} $${S3_BUCKET} >>/tmp/dumps.log`"
  # Select the shell to run 'command' with.
  #shell = "/bin/sh"

  # UI / CLI annotation informational messages:
  start_short_template    = "Dumping JVM info."
  error_short_template    = "Error dumping JVM info."
  complete_short_template = "Finished dumping JVM info."
  start_long_template     = "Dumping JVM process ${var.jvm_process_regex} info."
  error_long_template     = "Error dumping JVM process ${var.jvm_process_regex} info."
  complete_long_template  = "Finished dumping JVM process ${var.jvm_process_regex} info."

  enabled = true
}
```

Lastly, we connect the `jvm_trace_heap_alarm` [Alarm](/t/alarm) and the `jvm_trace_check_heap` [Action](/t/action) with the `jvm_trace_dump_bot` [Bot](/t/bot):

{{tffile "examples/op_packs/jvm_trace/bots.tf"}}

Now, anytime JVM memory exceeds our defined threshold the JVM is automatically restarted and the debug data is exported for further analysis.

-> See the Shoreline [Bots Documentation](/t/bot) for more info.

{{ .SchemaMarkdown | trimspace }}