// Copyright 2019 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package nrcensus

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/newrelic/newrelic-telemetry-sdk-go/telemetry"
	"go.opencensus.io/trace"
)

var (
	testTime     = time.Date(2014, time.November, 28, 1, 1, 0, 0, time.UTC)
	testSpanID   = trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8}
	testParentID = trace.SpanID{9, 10, 11, 12, 13, 14, 15, 16}
	testTraceID  = trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
)

type testHarvester struct {
	spans   []telemetry.Span
	metrics []telemetry.Metric
}

func (h *testHarvester) RecordMetric(m telemetry.Metric) {
	h.metrics = append(h.metrics, m)
}
func (h *testHarvester) RecordSpan(sp telemetry.Span) error {
	h.spans = append(h.spans, sp)
	return nil
}

func TestSpanGeneric(t *testing.T) {
	h := &testHarvester{}
	exp := &Exporter{
		Harvester:   h,
		ServiceName: "serviceName",
	}
	sd := &trace.SpanData{
		SpanContext: trace.SpanContext{
			SpanID:  testSpanID,
			TraceID: testTraceID,
		},
		Name:      "spanName",
		StartTime: testTime,
		EndTime:   testTime.Add(time.Second),
		Attributes: map[string]interface{}{
			"color": "purple",
		},
	}
	exp.ExportSpan(sd)
	if span := h.spans[0]; !reflect.DeepEqual(span, telemetry.Span{
		ID:          "0102030405060708",
		TraceID:     "0102030405060708090a0b0c0d0e0f10",
		Name:        "spanName",
		ParentID:    "",
		ServiceName: "serviceName",
		Timestamp:   testTime,
		Duration:    time.Second,
		Attributes: map[string]interface{}{
			"color":                    "purple",
			"instrumentation.provider": instrumentationProvider,
			"collector.name":           collectorName,
		},
	}) {
		t.Errorf("span fields are incorrect: %#v", span)
	}
}

func TestSpanParentID(t *testing.T) {
	// test that when available, the parent span id is recorded
	h := &testHarvester{}
	exp := &Exporter{
		Harvester:   h,
		ServiceName: "serviceName",
	}
	sd := &trace.SpanData{
		SpanContext: trace.SpanContext{
			SpanID:  testSpanID,
			TraceID: testTraceID,
		},
		ParentSpanID: testParentID,
		Name:         "spanName",
		StartTime:    testTime,
		EndTime:      testTime.Add(time.Second),
		Attributes: map[string]interface{}{
			"color": "purple",
		},
	}
	exp.ExportSpan(sd)
	if span := h.spans[0]; !reflect.DeepEqual(span, telemetry.Span{
		ID:          "0102030405060708",
		TraceID:     "0102030405060708090a0b0c0d0e0f10",
		Name:        "spanName",
		ParentID:    "090a0b0c0d0e0f10",
		ServiceName: "serviceName",
		Timestamp:   testTime,
		Duration:    time.Second,
		Attributes: map[string]interface{}{
			"color":                    "purple",
			"instrumentation.provider": instrumentationProvider,
			"collector.name":           collectorName,
		},
	}) {
		t.Errorf("span fields are incorrect: %#v", span)
	}
}

func TestSpanErrorAttrHonored(t *testing.T) {
	h := &testHarvester{}
	exp := &Exporter{
		Harvester:   h,
		ServiceName: "serviceName",
	}
	sd := &trace.SpanData{
		StartTime: testTime,
		EndTime:   testTime.Add(time.Second),
		Attributes: map[string]interface{}{
			"error": "hello world",
		},
		Status: trace.Status{Code: 1},
	}
	exp.ExportSpan(sd)
	if span := h.spans[0]; !reflect.DeepEqual(span, telemetry.Span{
		ID:          "0000000000000000",
		TraceID:     "00000000000000000000000000000000",
		ServiceName: "serviceName",
		Timestamp:   testTime,
		Duration:    time.Second,
		Attributes: map[string]interface{}{
			"error":                    "hello world",
			"instrumentation.provider": instrumentationProvider,
			"collector.name":           collectorName,
		},
	}) {
		t.Errorf("span fields are incorrect: %#v", span)
	}
}

func TestSpanErrorIgnoredHonored(t *testing.T) {
	h := &testHarvester{}
	exp := &Exporter{
		Harvester:   h,
		ServiceName: "serviceName",
	}
	exp.IgnoreStatusCodes = append(exp.IgnoreStatusCodes, 1)
	sd := &trace.SpanData{
		StartTime: testTime,
		EndTime:   testTime.Add(time.Second),
		Status:    trace.Status{Code: 1},
	}
	exp.ExportSpan(sd)
	if span := h.spans[0]; !reflect.DeepEqual(span, telemetry.Span{
		ID:          "0000000000000000",
		TraceID:     "00000000000000000000000000000000",
		ServiceName: "serviceName",
		Timestamp:   testTime,
		Duration:    time.Second,
		Attributes: map[string]interface{}{
			"instrumentation.provider": instrumentationProvider,
			"collector.name":           collectorName,
		},
	}) {
		t.Errorf("span fields are incorrect: %#v", span)
	}
}

func TestSpanErrorRecorded(t *testing.T) {
	h := &testHarvester{}
	exp := &Exporter{
		Harvester:   h,
		ServiceName: "serviceName",
	}
	sd := &trace.SpanData{
		StartTime: testTime,
		EndTime:   testTime.Add(time.Second),
		Status:    trace.Status{Code: 1},
	}
	exp.ExportSpan(sd)
	if span := h.spans[0]; !reflect.DeepEqual(span, telemetry.Span{
		ID:          "0000000000000000",
		TraceID:     "00000000000000000000000000000000",
		ServiceName: "serviceName",
		Timestamp:   testTime,
		Duration:    time.Second,
		Attributes: map[string]interface{}{
			"error":                    true,
			"instrumentation.provider": instrumentationProvider,
			"collector.name":           collectorName,
		},
	}) {
		t.Errorf("span fields are incorrect: %#v", span)
	}
}

func TestSpanInstAttrsOverwritten(t *testing.T) {
	h := &testHarvester{}
	exp := &Exporter{
		Harvester:   h,
		ServiceName: "serviceName",
	}
	sd := &trace.SpanData{
		StartTime: testTime,
		EndTime:   testTime.Add(time.Second),
		Attributes: map[string]interface{}{
			"instrumentation.provider": "invalid provider",
			"collector.name":           "invalid collector",
		},
	}
	exp.ExportSpan(sd)
	sd.Attributes = map[string]interface{}{"instrumentation.provider": "totally invalid provider"}
	exp.ExportSpan(sd)
	sd.Attributes = map[string]interface{}{"collector.name": "totally invalid collector"}
	exp.ExportSpan(sd)
	want := map[string]interface{}{
		"instrumentation.provider": instrumentationProvider,
		"collector.name":           collectorName,
	}
	for _, s := range h.spans {
		if !reflect.DeepEqual(s.Attributes, want) {
			t.Errorf("invalid attributes: got %#v, want %#v", s.Attributes, want)
		}
	}
}

type testIDGen struct {
	cnt int
}

func (*testIDGen) NewTraceID() [16]byte { return testTraceID }
func (g *testIDGen) NewSpanID() [8]byte {
	g.cnt++
	if g.cnt == 1 {
		return testParentID
	}
	return testSpanID
}

func TestSpanUsingOpenCensusAPI(t *testing.T) {
	h := &testHarvester{}
	exp := &Exporter{
		Harvester:   h,
		ServiceName: "serviceName",
	}
	trace.RegisterExporter(exp)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
		IDGenerator:    &testIDGen{},
	})

	ctx, parent := trace.StartSpan(context.Background(), "parent")
	_, child := trace.StartSpan(ctx, "child")
	child.SetStatus(trace.Status{
		Code:    trace.StatusCodePermissionDenied,
		Message: "oops permission denied",
	})
	child.End()
	parent.End()

	// first span is the child
	if childSpan := h.spans[0]; !reflect.DeepEqual(childSpan, telemetry.Span{
		ID:          "0102030405060708",
		TraceID:     "0102030405060708090a0b0c0d0e0f10",
		Name:        "child",
		ParentID:    "090a0b0c0d0e0f10",
		ServiceName: "serviceName",
		Timestamp:   childSpan.Timestamp,
		Duration:    childSpan.Duration,
		Attributes: map[string]interface{}{
			"error":                    true,
			"instrumentation.provider": instrumentationProvider,
			"collector.name":           collectorName,
		},
	}) {
		t.Errorf("child span fields are incorrect: %#v", childSpan)
	}

	// second span is the parent
	if parentSpan := h.spans[1]; !reflect.DeepEqual(parentSpan, telemetry.Span{
		ID:          "090a0b0c0d0e0f10",
		TraceID:     "0102030405060708090a0b0c0d0e0f10",
		Name:        "parent",
		ParentID:    "",
		ServiceName: "serviceName",
		Timestamp:   parentSpan.Timestamp,
		Duration:    parentSpan.Duration,
		Attributes: map[string]interface{}{
			"instrumentation.provider": instrumentationProvider,
			"collector.name":           collectorName,
		},
	}) {
		t.Errorf("parent span fields are incorrect: %#v", parentSpan)
	}
}

func TestSpanNilHarvester(t *testing.T) {
	exp := &Exporter{
		ServiceName: "serviceName",
	}
	sd := &trace.SpanData{
		SpanContext: trace.SpanContext{
			SpanID:  testSpanID,
			TraceID: testTraceID,
		},
		Name:      "spanName",
		StartTime: testTime,
		EndTime:   testTime.Add(time.Second),
		Attributes: map[string]interface{}{
			"color": "purple",
		},
	}
	exp.ExportSpan(sd)
}

func TestSpanNilExporter(t *testing.T) {
	var exp *Exporter
	sd := &trace.SpanData{
		SpanContext: trace.SpanContext{
			SpanID:  testSpanID,
			TraceID: testTraceID,
		},
		Name:      "spanName",
		StartTime: testTime,
		EndTime:   testTime.Add(time.Second),
		Attributes: map[string]interface{}{
			"color": "purple",
		},
	}
	exp.ExportSpan(sd)
}

func TestSpanAttrLen(t *testing.T) {
	exp := &Exporter{IgnoreStatusCodes: []int32{}}
	sd := &trace.SpanData{
		SpanContext: trace.SpanContext{
			SpanID:  testSpanID,
			TraceID: testTraceID,
		},
		Name:      "spanName",
		StartTime: testTime,
		EndTime:   testTime.Add(time.Second),
		Attributes: map[string]interface{}{
			"instrumentation.provider": instrumentationProvider,
			"collector.name":           collectorName,
		},
	}

	// If desired attributes already exist and no errors, the length should
	// just be the existing number of attributes.
	want := len(sd.Attributes)
	if got := exp.spanAttrLen(sd); got != want {
		t.Errorf("unexpected number of attributes: want %d, got %d", want, got)
	}

	// Removing the expected attributes means the length returned should be
	// that many more.
	want = 2
	sd.Attributes = map[string]interface{}{}
	if got := exp.spanAttrLen(sd); got != want {
		t.Errorf("unexpected number of attributes: want %d, got %d", want, got)
	}

	// Marking as errored should add an additional attribute.
	want++
	sd.Status = trace.Status{Code: trace.StatusCodePermissionDenied}
	if got := exp.spanAttrLen(sd); got != want {
		t.Errorf("unexpected number of attributes: want %d, got %d", want, got)
	}
}
