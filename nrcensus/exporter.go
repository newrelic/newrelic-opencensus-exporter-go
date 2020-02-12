// Copyright 2019 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package nrcensus

import (
	"github.com/newrelic/newrelic-telemetry-sdk-go/cumulative"
	"github.com/newrelic/newrelic-telemetry-sdk-go/telemetry"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

// Exporter implements trace.Exporter
// (https://godoc.org/go.opencensus.io/trace#Exporter) and view.Exporter
// (https://godoc.org/go.opencensus.io/stats/view#Exporter).  It enables sending of
// trace and view data from OpenCensus applications to New Relic.
type Exporter struct {
	// Harvester is expected to be populated by the *telemetry.Harvester
	// (https://godoc.org//github.com/newrelic/newrelic-telemetry-sdk-go/telemetry#Harvester)
	// to use for reporting trace and view data.  It is an interface here to
	// facilitate testing.
	Harvester interface {
		RecordSpan(telemetry.Span) error
		RecordMetric(telemetry.Metric)
	}
	// ServiceName is the name of this service or application.
	ServiceName string
	// IgnoreStatusCodes controls which trace.Status
	// (https://opencensus.io/tracing/span/status/) Codes are turned into
	// errors on Spans.  A Span with a trace.Status greater than 0 that is not
	// in this slice will be marked as an error.  When instantiated with
	// NewExporter this field defaults to only include 5 (NOT_FOUND).
	IgnoreStatusCodes []int32
	// DeltaCalculator translates OpenCensus's cumulative metrics into delta
	// metrics.  This field must be populated to record metrics, as is done by
	// NewExporter.
	//
	// This is a cache which is cleared by default every 20 minutes.  If your
	// metrics are being recorded on an intermittent basis, you may have to
	// modify the cache cleaning interval on this DeltaCalculator in order to
	// avoid missing metrics or spikes in graphs when your data is assimilated.
	DeltaCalculator *cumulative.DeltaCalculator
}

var emptySpanID trace.SpanID

// NewExporter creates a new Exporter.  serviceName is the name of this service
// or application.  apiKey is required and refers to a New Relic Insights Insert API key.
func NewExporter(serviceName, apiKey string, options ...func(*telemetry.Config)) (*Exporter, error) {
	options = append([]func(*telemetry.Config){
		func(cfg *telemetry.Config) {
			cfg.Product = userAgentProduct
			cfg.ProductVersion = version
		},
		telemetry.ConfigAPIKey(apiKey),
	}, options...)
	h, err := telemetry.NewHarvester(options...)
	if nil != err {
		return nil, err
	}
	return &Exporter{
		Harvester:         h,
		ServiceName:       serviceName,
		IgnoreStatusCodes: []int32{5},
		DeltaCalculator:   cumulative.NewDeltaCalculator(),
	}, nil
}

func (e *Exporter) responseCodeIsError(code int32) bool {
	if code <= 0 {
		return false
	}
	for _, ignoreCode := range e.IgnoreStatusCodes {
		if code == ignoreCode {
			return false
		}
	}
	return true
}

// ExportSpan implements trace.Exporter and records spans with the Harvester
// for later sending to New Relic.
func (e *Exporter) ExportSpan(s *trace.SpanData) {
	if nil == e {
		return
	}

	// This is a somewhat expensive call, so be sure to only do this once.
	isErr := e.responseCodeIsError(s.Status.Code)
	// Make a new attribute map instead of updating the original in order to
	// not change the passed attributes.
	attrs := make(map[string]interface{}, e.spanAttrLen(s.Attributes, isErr))
	for k, v := range s.Attributes {
		attrs[k] = v
	}
	// Preserve any passed `error` attribute.
	if _, in := s.Attributes["error"]; !in && isErr {
		attrs["error"] = true
	}
	// This exporter defines these values, overwrite if they exist.
	attrs["instrumentation.provider"] = instrumentationProvider
	attrs["collector.name"] = collectorName

	sp := telemetry.Span{
		ID:          s.SpanContext.SpanID.String(),
		TraceID:     s.SpanContext.TraceID.String(),
		Name:        s.Name,
		Timestamp:   s.StartTime,
		Duration:    s.EndTime.Sub(s.StartTime),
		ServiceName: e.ServiceName,
		Attributes:  attrs,
	}

	if s.ParentSpanID != emptySpanID {
		sp.ParentID = s.ParentSpanID.String()
	}

	if nil == e.Harvester {
		return
	}
	e.Harvester.RecordSpan(sp)
}

// spanAttrLen returns the number of attributes that will be exported based on
// the OpenCensus attrs and if the span isErr.
func (e *Exporter) spanAttrLen(attrs map[string]interface{}, isErr bool) int {
	l := len(attrs)
	if _, in := attrs["error"]; !in && isErr {
		l++
	}
	if _, in := attrs["instrumentation.provider"]; !in {
		l++
	}
	if _, in := attrs["collector.name"]; !in {
		l++
	}
	return l
}

func (e *Exporter) recordCountData(vd *view.Data, data *view.CountData, attrs map[string]interface{}) {
	metric, ok := e.DeltaCalculator.CountMetric(vd.View.Name, attrs, float64(data.Value), vd.End)
	if !ok {
		metric.Name = vd.View.Name
		metric.Attributes = attrs
		metric.Value = float64(data.Value)
		metric.Timestamp = vd.Start
		metric.Interval = vd.End.Sub(vd.Start)
	}
	e.Harvester.RecordMetric(metric)
}

func (e *Exporter) recordLastValueData(vd *view.Data, data *view.LastValueData, attrs map[string]interface{}) {
	e.Harvester.RecordMetric(telemetry.Gauge{
		Name:       vd.View.Name,
		Attributes: attrs,
		Value:      data.Value,
		Timestamp:  vd.End,
	})
}

func (e *Exporter) recordSumData(vd *view.Data, data *view.SumData, attrs map[string]interface{}) {
	metric, ok := e.DeltaCalculator.CountMetric(vd.View.Name, attrs, data.Value, vd.End)
	if !ok {
		metric.Name = vd.View.Name
		metric.Attributes = attrs
		metric.Value = data.Value
		metric.Timestamp = vd.Start
		metric.Interval = vd.End.Sub(vd.Start)
	}
	e.Harvester.RecordMetric(metric)
}

// ExportView implements view.Exporter and records metrics with the Harvester
// for later sending to New Relic.
func (e *Exporter) ExportView(vd *view.Data) {
	if nil == e {
		return
	}
	if nil == e.Harvester {
		return
	}
	if nil == e.DeltaCalculator {
		return
	}
	for _, row := range vd.Rows {
		attrs := make(map[string]interface{}, len(row.Tags)+5)
		for _, tag := range row.Tags {
			attrs[tag.Key.Name()] = tag.Value
		}
		attrs["instrumentation.provider"] = instrumentationProvider
		attrs["collector.name"] = collectorName
		attrs["measure.name"] = vd.View.Measure.Name()
		attrs["measure.unit"] = vd.View.Measure.Unit()
		attrs["service.name"] = e.ServiceName

		switch data := row.Data.(type) {
		case *view.CountData:
			e.recordCountData(vd, data, attrs)
		case *view.SumData:
			e.recordSumData(vd, data, attrs)
		case *view.LastValueData:
			e.recordLastValueData(vd, data, attrs)
		case *view.DistributionData:
		default:
		}
	}
}
