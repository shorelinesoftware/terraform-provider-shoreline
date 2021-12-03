---
page_title: "shoreline_metric Resource - terraform-provider-shoreline"
subcategory: ""
description: |- Shoreline metric. A periodic measurement of a system property.
---

# shoreline_metric (Resource)

Shoreline [Metrics](/t/metric) are time-series data obtained through tools such as [Prometheus](https://prometheus.io/)
and the [Shoreline Metric Scraper](/installation/kubernetes/metric-scraper). In addition to
Shoreline's [standard Metrics](/t/metric/standard) you can define your own custom metrics using the `shoreline_metric`
resource.

## Required Properties

A [Metric](/t/metric) has only a few required properties:

- name - A unique name for the metric resource.
- value - A valid [Op](/t/op) statement that defines a valid [Metric query](/t/metric#creating-new-metrics).

## Usage

The following example creates a [Metric](/t/metric) named `avg_cpu_usage_prev_min` that gets the average CPU usage over
the last minute:

```tf
resource "shoreline_metric" "avg_cpu_usage_prev_min" {
  name = "avg_cpu_usage_prev_min"
  value = "cpu_usage | window(60s) | mean(60)"
}
```

### Advanced Usage

You can also create your own [Metric queries](/t/metric-query) using any ingested metric known to the
associated [Shoreline Agent](/t/agent). In the following example we're using Elasticsearch's cluster health status to
create a custom `elasticsearch_red_status` [Metric](/t/metric):

```tf
resource "shoreline_metric" "elasticsearch_red_status" {
  name = "${var.namespace}_elasticsearch_red_status"
  value = "metric_query(metric_names=\"elasticsearch_cluster_health_status\") | color=\"red\""
  description = "Returns 1 if resource has status red else 0"
}
```

The `elasticsearch_red_status` [Metric](/t/metric) can now be used as the basis for an [Alarm](/t/alarm):

```tf
resource "shoreline_alarm" "elasticsearch_status_alarm" {
  name = "${var.prefix}_elasticsearch_status_alarm"
  fire_query = "${shoreline_metric.elasticsearch_red_status.name} == 1"
  clear_query = "${shoreline_metric.elasticsearch_red_status.name} == 0"
  description = "Watch Elasticsearch health status."
  resource_query = "pods"
  enabled = true
}
```

-> See the Shoreline [Metrics](/t/metric) and [Metric Scraper](/installation/kubernetes/metric-scraper) documentation
for more info.

{{ .SchemaMarkdown | trimspace }}

