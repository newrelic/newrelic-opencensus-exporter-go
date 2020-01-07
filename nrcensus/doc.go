// Copyright 2019 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// Package nrcensus provides an exporter for sending OpenCensus stats and
// traces to New Relic.
//
// To use, simply instantiate a new Exporter using `nrcensus.NewExporter` with
// your service name and Insights API key and register it with the OpenCensus
// view and/or trace API.
//
//    exporter, err := nrcensus.NewExporter("My-OpenCensus-App", "__YOUR_NEW_RELIC_INSIGHTS_API_KEY__")
//    if err != nil {
//        panic(err)
//    }
//    view.RegisterExporter(exporter)
//    trace.RegisterExporter(exporter)
//
//    // create stats, traces, etc
package nrcensus
