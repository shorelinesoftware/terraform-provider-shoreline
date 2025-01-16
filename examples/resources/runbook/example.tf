resource "shoreline_notebook" "my_notebook" {
  name        = "my_notebook"
  description = "A sample notebook."
  data        = file("${path.module}/data/my_notebook.json")
  enabled     = true
}