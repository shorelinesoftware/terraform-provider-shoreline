
resource "shoreline_action" "delayed_action" {
  name = "delayed_action"
  # sleeps for a while and creates a txt file in the directory
  command = "`sleep 15; touch $${dir}/test_file.txt`"

  description = "Create a text file after some delay"
  enabled     = true
  params      = ["dir"]

  start_title_template    = "Create test file action started"
  complete_title_template = "Create test file action completed"
  error_title_template    = "Create test file action failed"

  start_short_template    = "Create test file action short started"
  complete_short_template = "Create test file action short completed"
  error_short_template    = "Create test file action short failed"
}