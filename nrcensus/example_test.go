package nrcensus_test

import (
	"github.com/newrelic/newrelic-opencensus-exporter-go/nrcensus"
	"github.com/newrelic/newrelic-telemetry-sdk-go/telemetry"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

// To use, simply instantiate a new Exporter using `nrcensus.NewExporter` with
// your service name and Insights API key and register it with the OpenCensus
// view and/or trace API.
func Example() {
	exporter, err := nrcensus.NewExporter("My-OpenCensus-App", "__YOUR_NEW_RELIC_INSIGHTS_API_KEY__")
	if err != nil {
		panic(err)
	}
	view.RegisterExporter(exporter)
	trace.RegisterExporter(exporter)

	// create stats, traces, etc
}

func ExampleNewExporter() {
	// To enable Infinite Tracing on the New Relic Edge, use the
	// telemetry.ConfigSpansURLOverride along with the URL for your Trace
	// Observer, including scheme and path.  See
	// https://docs.newrelic.com/docs/understand-dependencies/distributed-tracing/enable-configure/enable-distributed-tracing
	exporter, err := nrcensus.NewExporter(
		"My-OpenCensus-App", "__YOUR_NEW_RELIC_INSIGHTS_API_KEY__",
		telemetry.ConfigSpansURLOverride("https://nr-internal.aws-us-east-1.tracing.edge.nr-data.net/trace/v1"),
	)
	if err != nil {
		panic(err)
	}
	view.RegisterExporter(exporter)
	trace.RegisterExporter(exporter)

	// create stats, traces, etc
}
