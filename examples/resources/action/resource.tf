
resource "shoreline_action" "ls_action" {
  name           = "ls_action"
  command        = "`ls $${dir}; export FOO='bar'`"
  description    = "List some files ..."
  resource_query = "host"
  params         = ["dir"]
  res_env_var    = "FOO"
  #timeout = 60
  start_title_template    = "JVM dump started"
  complete_title_template = "JVM dump completed"
  error_title_template    = "JVM dump failed"

  start_short_template    = "JVM dump short started"
  complete_short_template = "JVM dump short completed"
  error_short_template    = "JVM dump short failed"

  enabled = true
}

