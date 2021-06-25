
resource "shoreline_resource" "books" {
  name        = "books"
  description = "Pods with books app."
  value       = "host | pod | app='bookstore'"
}

