
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