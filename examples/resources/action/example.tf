
resource "shoreline_action" "full_action" {
  name                    = "full_action"
  command                 = "`ls $${dir}; export FOO='bar'`"
  description             = "List some files ..."
  resource_query          = "host"
  params                  = ["dir"]
  res_env_var             = "FOO"
  timeout                 = 60
  shell                   = "/bin/zsh"
  allowed_resources_query = "host"
  allowed_entities        = ["<user_1>", "<user_2>"]
  editors                 = ["<user_2>", "<user_3>"]
  resource_tags_to_export = ["<tag_1>", "<tag_2>"]
  # file_deps               = ["<file1>", "<file2>"]
  communication_workspace = "<workspace_name>"
  communication_channel   = "<workspace_channel>"

  start_title_template    = "JVM dump started"
  complete_title_template = "JVM dump completed"
  error_title_template    = "JVM dump failed"

  start_short_template    = "JVM dump short started"
  complete_short_template = "JVM dump short completed"
  error_short_template    = "JVM dump short failed"

  start_long_template    = "JVM dump long started"
  error_long_template    = "JVM dump long completed"
  complete_long_template = "JVM dump long failed"

  enabled = true
}


resource "shoreline_action" "minimal_action" {
  name    = "minimal_action"
  command = "`ls $${dir}; export FOO='bar'`"
}


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

