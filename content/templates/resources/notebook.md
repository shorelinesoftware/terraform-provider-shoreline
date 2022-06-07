---
page_title: "shoreline_notebook Resource - terraform-provider-shoreline"
subcategory: ""
description: |- Notebooks replace static runbooks by capturing interactive debug and remediation sessions in a convenient UI.
---

# shoreline_notebook (Resource)

Through Shoreline's web-based UI, [Notebooks](/t/notebook) automatically capture an entire debug and remediation session -- which can optionally be associated with a specific [Alarm](/t/alarm) -- and then be shared with other team members to streamline future incident response. With Notebooks you can:

- Create a series of interactive [Op statement cells](/t/notebook#op-statements) allowing you to execute [Op commands](/t/op/command) within your browser -- all without installing or configuring the local [CLI](/t/cli).
- Define and use [dynamic parameters](/t/notebook-param) across Notebook Op cells.
- Memorialize Notebooks with [historical snapshots](/t/notebook-run).
- Add [Markdown-based notes](/t/notebook#notes) to inform operators how to use the Notebook.
- [Associate](/t/notebook#alarm-association) existing [Alarms](/t/alarm) and Notebooks, allowing on-call members to click through to an interactive debugging and remediation Notebook directly from the triggered [Alarm](/t/alarm) UI.

## Required Properties

Each Notebook uses a variety of properties to determine its behavior. The required properties when [creating a Notebook](/t/notebook#create-a-notebook) are:

- `name`: string - A unique symbol name for the Notebook object.
- `data`: string - A Notebook in JSON format.  Typically, this string should be generated using the Terraform [file](https://www.terraform.io/language/functions/file) function while referencing a local file `path`.

### Download a Notebook

You can download an entire Notebook in JSON format from the Shoreline Notebooks dashboard.

1. Click the **Actions** button at the top-right of the active Notebook panel.
2. Select the **Download Notebook** button to download the full configuration of the current Notebook to a local JSON file.

   You can also freely modify, share, and [upload](/t/notebook#upload-a-notebook) this Notebook at any time.

## Usage

The following example creates an [Notebook](/t/notebook) named `my_notebook`.

1. [Download a Notebook](#download-a-notebook) as JSON.
2. Save the Notebook JSON to local path within your Terraform project.
3. Define a new `shoreline_notebook` Terraform resource in your Terraform configuration that points the `data` property to the correct local module path.

   {{tffile "examples/resources/notebook/resource.tf"}}

-> See the Shoreline [Notebooks Documentation](/t/notebook) for more info on creating and using Notebooks.

{{ .SchemaMarkdown | trimspace }}
