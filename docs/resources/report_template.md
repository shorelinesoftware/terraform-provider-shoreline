---
page_title: 'shoreline_report_template Resource - terraform-provider-shoreline'
subcategory: ''
description: |-
---

# shoreline_report_template (Resource)


## Properties


- <b>name</b> (String) The name of the Report Template.
- <b>display_name</b> (String) The human-friendly name displayed in the UI.
- <b>links</b> (List of Object) References to other related Report Templates. A list of objects with the following attributes:
    - <b>label</b> (String) A label for the link.
    - <b>report_template_name</b> (String) The name of the linked Report Template.
- <b>blocks</b> (String) A list of Report Template blocks in JSON format. Typically, this string should be set using the [jsonencode](https://developer.hashicorp.com/terraform/language/functions/jsonencode) function with a [Terraform Input Variable](https://developer.hashicorp.com/terraform/language/values/variables). All it's properties must be present to avoid Terraform diffs. It has the following properties:
    - <i><b>title</b></i> (String) The name of the report block.
    - <i><b>resource_query</b></i> (String) Specifies which resources to include in the chart.
    - <i><b>group_by_tag</b></i> (String) The resource tag used to group resources in the chart.
    - <i><b>breakdown_by_tag</b></i> (String) The tag within each group used to further break down resources.
    - <i><b>breakdown_tags_values</b></i> (List of Object) Specifies which values of the breakdown tag to display in the chart. It's a list of objects with the following attributes:
        - <i><b>color</b></i> (String) The hexadecimal color code (`#RRGGBB`).
        - <i><b>values</b></i> (List of String) Tag values.
        - <i><b>label</b></i> (String) A label.
    - <i><b>include_other_breakdown_tag_values</b></i> (Boolean) When set to `true`, resources that do not have a value set for the breakdown tag are included in a separate `other` section of the specific row.
    - <i><b>view_mode</b></i> (String) Determines the display format for the bar charts, either as a `COUNT` (numerical count) or `PERCENTAGE` (percentage of the whole).
    - <i><b>resources_breakdown</b></i> (List of Object) Contains all data necessary for building the chart. It's a list of objects with the following attributes:
        - <i><b>group_by_value</b></i> (String) Existing tag value or `__no_value__`.
        - <i><b>breakdown_values</b></i> (List of Object) A list of objects with the following attributes:
            - <i><b>value</b></i> (String) Existing tag value or `__no_value__`.
            - <i><b>count</b></i> (Number)
    - <i><b>other_tags_to_export</b></i> (List of String) Additional tags (besides the group and breakdown tags) to include when exporting the Report Template.
    - <i><b>include_resources_without_group_tag</b></i> (Boolean) When set to `true`, resources without a group tag value are included in the chart in an another row labeled `other`.
    - <i><b>group_by_tag_order</b></i> (Object) Defines the display order for the values of the group by tag in the chart. Has the following attributes:
        - <i><b>type</b></i> (String) Can be one of the following: `DEFAULT`, `BY_TOTAL_ASC`, `BY_TOTAL_DESC`, `CUSTOM`.
        - <i><b>values</b></i> (List of String) If <b>type</b> is `CUSTOM`, this list defines the order of tags.



## Example


```terraform
resource "shoreline_report_template" "my_report_1" {
  name         = "my_report_1"
  display_name = "My First Report"
  blocks = jsonencode([
    {
      "title" : "Block Name",
      "resource_query" : "host",
      "group_by_tag" : "tag_0",
      "breakdown_by_tag" : "tag_1",
      "breakdown_tags_values" : [
        {
          "color" : "#AAAAAA",
          "values" : [
            "passed",
            "skipped"
          ],
          "label" : "label_0"
        }
      ],
      "view_mode" : "PERCENTAGE",
      "include_other_breakdown_tag_values" : true,
      "other_tags_to_export" : ["other_tag_1", "other_tag_2"],
      "include_resources_without_group_tag" : false,
      "group_by_tag_order" : {
        "type" : "DEFAULT",
        "values" : []
      },
      "resources_breakdown" : [
        {
          "group_by_value" : "tag_0",
          "breakdown_values" : [
            {
              "value" : "value",
              "count" : 1
            }
          ]
        }
      ]
    }
  ])
  depends_on = [
    shoreline_report_template.my_report_2
  ]
  links = jsonencode([{
    "label" : "second-report",
    "report_template_name" : "my_report_2"
  }])
}


resource "shoreline_report_template" "my_report_2" {
  name         = "my_report_2"
  display_name = "My Second Report"
  blocks = jsonencode([
    {
      "title" : "Block Name",
      "resource_query" : "host",
      "group_by_tag" : "tag_0",
      "breakdown_by_tag" : "tag_1",
      "breakdown_tags_values" : [
        {
          "color" : "#AAAAAA",
          "values" : [
            "passed",
            "skipped"
          ],
          "label" : "label_0"
        }
      ],
      "view_mode" : "COUNT",
      "include_other_breakdown_tag_values" : true,
      "other_tags_to_export" : ["other_tag_1", "other_tag_2"],
      "include_resources_without_group_tag" : false,
      "group_by_tag_order" : {
        "type" : "DEFAULT",
        "values" : []
      },
      "resources_breakdown" : [
        {
          "group_by_value" : "tag_0",
          "breakdown_values" : [
            {
              "value" : "value",
              "count" : 1
            }
          ]
        }
      ]
    }
  ])
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **blocks** (String) The JSON encoded blocks of the report template.
- **name** (String) The name/symbol for the object within Shoreline and the op language (must be unique, only alphanumeric/underscore).

### Optional

- **display_name** (String) A user friendly name shown in the UI.
- **id** (String) The ID of this resource.
- **links** (String) The JSON encoded links of a report template with other report templates. Defaults to `[]`.

### Read-Only

- **type** (String) The type of object (i.e., Alarm, Action, Bot, Metric, Resource, or File).
