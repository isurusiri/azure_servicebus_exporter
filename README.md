# Azure ServiceBus Exporter

My humble attempt to develop a custom [Prometheus](https://prometheus.io/) metrics exporter. This exporter exposes metrics of Azure ServuceBus.

### But why exporters?

Applications and services that we are useing may expose various metrics such as network I/O, disk utilization, and many more. However all of these services may not expose metrics or exposed metrics are not compatible with Prometheus. Therefore, this is where metrics exporters helps us.

Metrics exporters acts as an adapter, and convert various metrics of our applications and services in to a Prometheus compatible way. A complete list of Prometheus exporters are availble [in here](https://prometheus.io/docs/instrumenting/exporters/).

Sometimes, we will not be able to find a suitable metrics exporter for our needs. In such as a case, writing our own custom metrics exporter is the best option.

```
         Scrape       API Call
   ......                      ......
   :    :  -->   ----   -->    :    :
   :    :        |  |          :    :
   :    :  <--   ----   <--    :    :
   ......      exporter        ......
 Prometheus              Application / Service
```

Once we have exporters exposing metrics to Prometheus, we can configure prometheus to generate alerts based on the rules configured and craete dashboards to monitor our applications.

### What's in this repository?

Like I have previously metnioned I wasnted to write a custom Prometheus metrics exporter and get experince. I followed [this tutorial](https://www.skyrise.tech/blog/tech/custom-prometheus-exporter/) for learning this and I have made some changes in here.

Here I have built a custom exporter for exposing metrics of Azure ServiceBus.

### Folder structure

```
|- client
| |- client.go
|- collector
| |- collector.go
|-servicebus_exporter.go
```

The `client.go` is responsible for collecting metrics from Azure ServiceBus. Collected metrics are returned as a struct called Stats.

The `collector.go` transform exposes collected metrics in a Prometheus compatible way.

The `servicebus_exporter.go` initiates metrics collection and transforming of it. Then finally expose metrics by running the HTTP server.
