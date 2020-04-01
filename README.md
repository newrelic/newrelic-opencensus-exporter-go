# New Relic Go OpenCensus exporter [![GoDoc](https://godoc.org/github.com/newrelic/newrelic-opencensus-exporter-go/nrcensus?status.svg)](https://godoc.org/github.com/newrelic/newrelic-opencensus-exporter-go/nrcensus)
The `nrcensus` package provides an exporter for sending OpenCensus stats and
traces to New Relic.


## Requirements
* [OpenCensus-Go](https://github.com/census-instrumentation/opencensus-go) v0.10.0 or newer
* Go v1.8 or newer


## Install
To install, just go get this package with

```
$ go get -u github.com/newrelic/newrelic-opencensus-exporter-go
```

## Using the exporter
To use the exporter, create a new Exporter and register it with OpenCensus.

```go
package main

import (
    "github.com/newrelic/newrelic-opencensus-exporter-go/nrcensus"
    "go.opencensus.io/stats/view"
    "go.opencensus.io/trace"
)

func main() {
    exporter, err := nrcensus.NewExporter("My-OpenCensus-App", "__YOUR_NEW_RELIC_INSIGHTS_API_KEY__")
    if err != nil {
        panic(err)
    }
    view.RegisterExporter(exporter)
    trace.RegisterExporter(exporter)

    // create stats, traces, etc
}
```

## Find and use your data

Tips on how to find and query your data:
- [Find metric data](https://docs.newrelic.com/docs/data-ingest-apis/get-data-new-relic/metric-api/introduction-metric-api#find-data)
- [Find trace/span data](https://docs.newrelic.com/docs/understand-dependencies/distributed-tracing/trace-api/introduction-trace-api#view-data)

For general querying information, see:
- [Query New Relic data](https://docs.newrelic.com/docs/using-new-relic/data/understand-data/query-new-relic-data)
- [Intro to NRQL](https://docs.newrelic.com/docs/query-data/nrql-new-relic-query-language/getting-started/introduction-nrql)


## Find and use your data

Tips on how to find and query your data:

- [Find metric data](https://docs.newrelic.com/docs/data-ingest-apis/get-data-new-relic/metric-api/introduction-metric-api#find-data)
- [Find trace/span data](https://docs.newrelic.com/docs/understand-dependencies/distributed-tracing/trace-api/introduction-trace-api#view-data)

For general querying information, see:

- [Query New Relic data](https://docs.newrelic.com/docs/using-new-relic/data/understand-data/query-new-relic-data)
- [Intro to NRQL](https://docs.newrelic.com/docs/query-data/nrql-new-relic-query-language/getting-started/nrql-syntax-clauses-functions)


## Licensing
The New Relic Go OpenCensus exporter is licensed under the Apache 2.0 License.
The New Relic Go OpenCensus exporter also uses source code from third party
libraries. Full details on which libraries are used and the terms under which
they are licensed can be found in the third party notices document.


## Contributing
Full details are available in our CONTRIBUTING.md file. We'd love to get your
contributions to improve the Go OpenCensus exporter! Keep in mind when you
submit your pull request, you'll need to sign the CLA via the click-through
using CLA-Assistant. You only have to sign the CLA one time per project. To
execute our corporate CLA, which is required if your contribution is on
behalf of a company, or if you have any questions, please drop us an email at
open-source@newrelic.com.


## Limitations
The New Relic Telemetry APIs are rate limited. Please reference the
documentation for [New Relic Metrics
API](https://docs.newrelic.com/docs/introduction-new-relic-metric-api) and [New
Relic Trace API Requirements and
Limits](https://docs.newrelic.com/docs/apm/distributed-tracing/trace-api/trace-api-general-requirements-limits)
on the specifics of the rate limits.