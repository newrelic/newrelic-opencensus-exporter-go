// Copyright 2019 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/newrelic/newrelic-opencensus-exporter-go/nrcensus"
	"github.com/newrelic/newrelic-telemetry-sdk-go/telemetry"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

func mustGetEnv(v string) string {
	val := os.Getenv(v)
	if val == "" {
		panic(fmt.Sprintf("%s unset", v))
	}
	return val
}

func main() {
	exporter, err := nrcensus.NewExporter("Example App",
		mustGetEnv("NEW_RELIC_INSIGHTS_INSERT_API_KEY"),
		telemetry.ConfigBasicErrorLogger(os.Stderr),
		telemetry.ConfigBasicDebugLogger(os.Stdout),
	)
	if nil != err {
		panic(err)
	}
	trace.RegisterExporter(exporter)
	view.RegisterExporter(exporter)

	// Always trace for this demo. In a production application, you should
	// configure this to a trace.ProbabilitySampler set at the desired
	// probability.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	client := &http.Client{Transport: &ochttp.Transport{}}

	m := stats.Float64("myMetric", "description", "inches")
	keyFirst, _ := tag.NewKey("first")
	keySecond, _ := tag.NewKey("second")
	countView := &view.View{
		Measure:     m,
		Name:        "MyCount",
		Description: "a count of the inches",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{keyFirst, keySecond},
	}
	if err := view.Register(countView); err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "hello world")

		ctx, _ := tag.New(context.Background(),
			tag.Insert(keyFirst, "firstValue"),
			tag.Insert(keySecond, "secondValue"))
		err := stats.RecordWithTags(ctx, nil, m.M(1.234))
		if nil != err {
			fmt.Println("error using RecordWithTags", err.Error())
		}

		// Provide an example of how spans can be annotated with metadata
		_, span := trace.StartSpan(req.Context(), "child")
		defer span.End()
		span.Annotate([]trace.Attribute{trace.StringAttribute("key", "value")}, "something happened")
		span.AddAttributes(trace.StringAttribute("hello", "world"))
		time.Sleep(time.Millisecond * 125)

		r, _ := http.NewRequest("GET", "https://example.com", nil)

		// Propagate the trace header info in the outgoing requests.
		r = r.WithContext(req.Context())
		resp, err := client.Do(r)
		if err != nil {
			log.Println(err)
		} else {
			// TODO: handle response
			resp.Body.Close()
		}
	})
	log.Fatal(http.ListenAndServe(":50030", &ochttp.Handler{}))
}
