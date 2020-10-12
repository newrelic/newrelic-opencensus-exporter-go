[![Community Project header](https://github.com/newrelic/open-source-office/raw/master/examples/categories/images/Community_Project.png)](https://github.com/newrelic/open-source-office/blob/master/examples/categories/index.md#community-project)

# New Relic Go OpenCensus exporter [![GoDoc](https://godoc.org/github.com/newrelic/newrelic-opencensus-exporter-go/nrcensus?status.svg)](https://godoc.org/github.com/newrelic/newrelic-opencensus-exporter-go/nrcensus)
The `nrcensus` package provides an exporter for sending OpenCensus stats and
traces to New Relic.

## Installation

Requirements:
* [OpenCensus-Go](https://github.com/census-instrumentation/opencensus-go) v0.10.0 or newer
* Go v1.8 or newer

To install, just go get this package with

```
$ go get -u github.com/newrelic/newrelic-opencensus-exporter-go
```

## Getting started
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

## Support

Should you need assistance with New Relic products, you are in good hands with several [**Optional** support diagnostic tools and] support channels.

If the issue has been confirmed as a bug or is a feature request, file a GitHub issue.

**Support Channels**
* [New Relic Documentation](LINK to specific docs page): Comprehensive guidance for using our platform
* [New Relic Community](LINK to specific community page): The best place to engage in troubleshooting questions
* [New Relic Developer](https://developer.newrelic.com/): Resources for building a custom observability applications
* [New Relic University](https://learn.newrelic.com/): A range of online training for New Relic users of every level

## Privacy
At New Relic we take your privacy and the security of your information seriously, and are committed to protecting your information. We must emphasize the importance of not sharing personal data in public forums, and ask all users to scrub logs and diagnostic information for sensitive information, whether personal, proprietary, or otherwise.

We define “Personal Data” as any information relating to an identified or identifiable individual, including, for example, your name, phone number, post code or zip code, Device ID, IP address, and email address.

For more information, review [New Relic’s General Data Privacy Notice](https://newrelic.com/termsandconditions/privacy).


## Contributing
We encourage your contributions to improve our OpenCensus Exporter! Keep in mind that when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. You only have to sign the CLA one time per project.

If you have any questions, or to execute our corporate CLA (which is required if your contribution is on behalf of a company), drop us an email at opensource@newrelic.com.

**A note about vulnerabilities**

As noted in our [security policy](../../security/policy), New Relic is committed to the privacy and security of our customers and their data. We believe that providing coordinated disclosure by security researchers and engaging with the security community are important means to achieve our security goals.

If you believe you have found a security vulnerability in this project or any of New Relic's products or websites, we welcome and greatly appreciate you reporting it to New Relic through [HackerOne](https://hackerone.com/newrelic).

If you would like to contribute to this project, review [these guidelines](./CONTRIBUTING.md).

To [all contributors](<LINK TO contributors>), we thank you!  Without your contribution, this project would not be what it is today.  We also host a community project page dedicated to [Project Name](<LINK TO https://opensource.newrelic.com/projects/... PAGE>).

## Licensing
The New Relic Go OpenCensus exporter is licensed under the Apache 2.0 License.
The New Relic Go OpenCensus exporter also uses source code from third party
libraries. Full details on which libraries are used and the terms under which
they are licensed can be found in the third party notices document.

## Limitations
The New Relic Telemetry APIs are rate limited. Please reference the
documentation for [New Relic Metrics
API](https://docs.newrelic.com/docs/introduction-new-relic-metric-api) and [New
Relic Trace API Requirements and
Limits](https://docs.newrelic.com/docs/apm/distributed-tracing/trace-api/trace-api-general-requirements-limits)
on the specifics of the rate limits.
