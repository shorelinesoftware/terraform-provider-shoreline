---
page_title: "shoreline_resource Resource - terraform-provider-shoreline"
subcategory: ""
description: |-
  Shoreline resource. A server or compute resource in the system (e.g. host, pod, container).
---

# shoreline_resource (Resource)

[Resources](/t/resource) are the infrastructure objects managed by Shoreline.  A [Resources](/t/resource) refers to a specific host, pod, or container throughout the Shoreline platform. It may also refer to a more distinct entity, such as a virtual machine or a database instance.

[Resources](/t/resource) types are platform and provider agnostic. So whether it's a fleet of pods and containers in Kubernetes, stand-alone hosts, or many virtual machines, Shoreline can control those [Resourcess](/t/resource) within AWS, GCP, or Azure.

-> Shoreline [Resources](/t/resource) are not to be confused with [Terraform resources](https://www.terraform.io/docs/language/resources/index.html).

## Required Properties

- name - A unique name for the [Resource](/t/resource) resource.
- value - A valid [Op](/t/op) statement that defines a valid [Resource query](/t/resource).

## Usage

Shoreline [Resources](/t/resource) are effectively aliases for more complex [Resource queries](/t/resource).  For example, the following definition creates a custom [Resource](/t/resource) called `az_k8s` as an alias for the full query to target Kubernetes pods on Azure:

```tf
resource "shoreline_resource" "az_k8s" {
  name        = "${var.namespace}_az_k8s"
  value       = "hosts | k8s=true | cloud_provider='azure' | pods | namespace=[\"${var.namespace}\"]"
  description = "All Shoreline Kubernetes pods on Azure"
}
```

-> See the Shoreline [Resources documentation](https://docs.shoreline.io/platform/resources) for more info.

{{ .SchemaMarkdown | trimspace }}