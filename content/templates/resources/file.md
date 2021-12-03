---
page_title: "shoreline_file Resource - terraform-provider-shoreline"
subcategory: ""
description: |-
  Automatically distribute artifacts across your Shoreline Resources.
---

# shoreline_file (Resource)

The [File](/t/file) allows you to easily transmit critical files throughout your fleet, even to ephemeral [Resources](/t/resource) such as Kubernetes (k8s) containers. This technique is potent when you need to distribute and execute custom bash scripts or other critical files without the need for manual intervention.

## Required Properties

A [File](/t/file) must know where to copy from and where to distribute to.  The required properties are:

- name - A unique object name.
- destination_path - The relative, destination file path.
- input_file - The relative, local file path of the source artifact.
- [resource_query](/t/resource) - The target Shoreline [Resources](/t/resource) to distribute the artifact to.

## Usage

The following example distributes the local `<terraform_module_directory>/data/jvm_dumps.sh` file the target Shoreline [Resources](/t/resource) defined by the `resource_query` Terraform variable:

{{tffile "examples/op_packs/jvm_trace/files.tf"}}

{{tffile "examples/op_packs/jvm_trace/variables.tf"}}

-> See the Shoreline [Op: `cp` documentation](https://docs.shoreline.io/op/commands/cp) for more info.

{{ .SchemaMarkdown | trimspace }}