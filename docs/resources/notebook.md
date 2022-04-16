---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "shoreline_notebook Resource - terraform-provider-shoreline"
subcategory: ""
description: |-
  Shoreline notebook. An interactive notebook of Op commands and user documentation .
  See the Shoreline Notebook Documentation https://docs.shoreline.io/ui/notebooks for more info.
---

# shoreline_notebook (Resource)

Shoreline notebook. An interactive notebook of Op commands and user documentation .

See the Shoreline [Notebook Documentation](https://docs.shoreline.io/ui/notebooks) for more info.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **cells** (String) The data cells inside a notebook.
- **name** (String) The name of the object (must be unique).

### Optional

- **description** (String) A user-friendly explanation of an object.
- **enabled** (Boolean) If the object is currently enabled or disabled. Defaults to `false`.
- **id** (String) The ID of this resource.
- **timeout_ms** (Number) Defaults to `60000`.

### Read-Only

- **type** (String) The type of object (i.e., Alarm, Action, Bot, Metric, Resource, or File).

