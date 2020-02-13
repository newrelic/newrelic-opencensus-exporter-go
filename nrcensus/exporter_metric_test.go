// Copyright 2019 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package nrcensus

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/newrelic/newrelic-telemetry-sdk-go/cumulative"
	"github.com/newrelic/newrelic-telemetry-sdk-go/telemetry"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	testKeyFirst, _  = tag.NewKey("first")
	testKeySecond, _ = tag.NewKey("second")
	testMeasure      = stats.Int64("tests", "a test measure", "t")
	testCountView    = &view.View{
		Measure:     testMeasure,
		Name:        "MyTestCount",
		Description: "a count of the test",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{testKeyFirst, testKeySecond},
	}
	testSumView = &view.View{
		Measure:     testMeasure,
		Name:        "MyTestSum",
		Description: "a sum of the test",
		Aggregation: view.Sum(),
		TagKeys:     []tag.Key{testKeyFirst, testKeySecond},
	}
	testLastValueView = &view.View{
		Measure:     testMeasure,
		Name:        "MyTestLastValue",
		Description: "a sum of the test",
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{testKeyFirst, testKeySecond},
	}
	testDistributionView = &view.View{
		Measure:     testMeasure,
		Name:        "MyTestLastValue",
		Description: "a sum of the test",
		Aggregation: view.Distribution(25, 100, 200, 400, 800, 10000),
		TagKeys:     []tag.Key{testKeyFirst, testKeySecond},
	}
)

func TestMetricNilExporter(t *testing.T) {
	var exp *Exporter
	vd := &view.Data{
		View:  testCountView,
		Start: testTime,
		End:   testTime.Add(10 * time.Second),
		Rows: []*view.Row{
			&view.Row{
				Tags: []tag.Tag{
					tag.Tag{Key: testKeyFirst, Value: "firstValue"},
					tag.Tag{Key: testKeySecond, Value: "secondValue"},
				},
				Data: &view.CountData{
					Value: 10,
				},
			},
		},
	}
	exp.ExportView(vd)
}

func TestMetricNilHarvester(t *testing.T) {
	exp := &Exporter{
		ServiceName:     "serviceName",
		DeltaCalculator: cumulative.NewDeltaCalculator(),
	}
	vd := &view.Data{
		View:  testCountView,
		Start: testTime,
		End:   testTime.Add(10 * time.Second),
		Rows: []*view.Row{
			&view.Row{
				Tags: []tag.Tag{
					tag.Tag{Key: testKeyFirst, Value: "firstValue"},
					tag.Tag{Key: testKeySecond, Value: "secondValue"},
				},
				Data: &view.CountData{
					Value: 10,
				},
			},
		},
	}
	exp.ExportView(vd)
}

func TestMetricNilDeltaCalculator(t *testing.T) {
	exp := &Exporter{
		Harvester:   &testHarvester{},
		ServiceName: "serviceName",
	}
	vd := &view.Data{
		View:  testCountView,
		Start: testTime,
		End:   testTime.Add(10 * time.Second),
		Rows: []*view.Row{
			&view.Row{
				Tags: []tag.Tag{
					tag.Tag{Key: testKeyFirst, Value: "firstValue"},
					tag.Tag{Key: testKeySecond, Value: "secondValue"},
				},
				Data: &view.CountData{
					Value: 10,
				},
			},
		},
	}
	exp.ExportView(vd)
}

func TestCountMetrics(t *testing.T) {
	h := &testHarvester{}
	exp := &Exporter{
		Harvester:       h,
		ServiceName:     "serviceName",
		DeltaCalculator: cumulative.NewDeltaCalculator(),
	}

	// first time metric is seen
	vd := &view.Data{
		View:  testCountView,
		Start: testTime,
		End:   testTime.Add(10 * time.Second),
		Rows: []*view.Row{
			&view.Row{
				Tags: []tag.Tag{
					tag.Tag{Key: testKeyFirst, Value: "firstValue"},
					tag.Tag{Key: testKeySecond, Value: "secondValue"},
				},
				Data: &view.CountData{
					Value: 10,
				},
			},
		},
	}
	exp.ExportView(vd)

	// second time metric is seen value does not change
	vd.End = testTime.Add(20 * time.Second)
	exp.ExportView(vd)

	// third time metric is seen value changes
	vd.End = testTime.Add(30 * time.Second)
	vd.Rows[0].Data.(*view.CountData).Value = 15
	exp.ExportView(vd)

	if metric := h.metrics[0]; !reflect.DeepEqual(metric, telemetry.Count{
		Name:      "MyTestCount",
		Value:     10,
		Timestamp: testTime,
		Interval:  10 * time.Second,
		Attributes: map[string]interface{}{
			"first":                    "firstValue",
			"second":                   "secondValue",
			"instrumentation.provider": instrumentationProvider,
			"collector.name":           collectorName,
			"measure.name":             "tests",
			"measure.unit":             "t",
			"service.name":             "serviceName",
		},
	}) {
		t.Errorf("metric fields are incorrect: %#v", metric)
	}
	if metric := h.metrics[1]; !reflect.DeepEqual(metric, telemetry.Count{
		Name:           "MyTestCount",
		Value:          0,
		Timestamp:      testTime.Add(10 * time.Second),
		Interval:       10 * time.Second,
		AttributesJSON: json.RawMessage(`{"collector.name":"` + collectorName + `","first":"firstValue","instrumentation.provider":"` + instrumentationProvider + `","measure.name":"tests","measure.unit":"t","second":"secondValue","service.name":"serviceName"}`),
	}) {
		t.Errorf("metric fields are incorrect: %#v", metric)
	}
	if metric := h.metrics[2]; !reflect.DeepEqual(metric, telemetry.Count{
		Name:           "MyTestCount",
		Value:          5,
		Timestamp:      testTime.Add(20 * time.Second),
		Interval:       10 * time.Second,
		AttributesJSON: json.RawMessage(`{"collector.name":"` + collectorName + `","first":"firstValue","instrumentation.provider":"` + instrumentationProvider + `","measure.name":"tests","measure.unit":"t","second":"secondValue","service.name":"serviceName"}`),
	}) {
		t.Errorf("metric fields are incorrect: %#v", metric)
	}
}

func TestSumMetrics(t *testing.T) {
	h := &testHarvester{}
	exp := &Exporter{
		Harvester:       h,
		ServiceName:     "serviceName",
		DeltaCalculator: cumulative.NewDeltaCalculator(),
	}

	// first time metric is seen
	vd := &view.Data{
		View:  testSumView,
		Start: testTime,
		End:   testTime.Add(10 * time.Second),
		Rows: []*view.Row{
			&view.Row{
				Tags: []tag.Tag{
					tag.Tag{Key: testKeyFirst, Value: "firstValue"},
					tag.Tag{Key: testKeySecond, Value: "secondValue"},
				},
				Data: &view.SumData{
					Value: 10,
				},
			},
		},
	}
	exp.ExportView(vd)

	// second time metric is seen value does not change
	vd.End = testTime.Add(20 * time.Second)
	exp.ExportView(vd)

	// third time metric is seen value changes
	vd.End = testTime.Add(30 * time.Second)
	vd.Rows[0].Data.(*view.SumData).Value = 15
	exp.ExportView(vd)

	if metric := h.metrics[0]; !reflect.DeepEqual(metric, telemetry.Count{
		Name:      "MyTestSum",
		Value:     10,
		Timestamp: testTime,
		Interval:  10 * time.Second,
		Attributes: map[string]interface{}{
			"first":                    "firstValue",
			"second":                   "secondValue",
			"instrumentation.provider": instrumentationProvider,
			"collector.name":           collectorName,
			"measure.name":             "tests",
			"measure.unit":             "t",
			"service.name":             "serviceName",
		},
	}) {
		t.Errorf("metric fields are incorrect: %#v", metric)
	}
	if metric := h.metrics[1]; !reflect.DeepEqual(metric, telemetry.Count{
		Name:           "MyTestSum",
		Value:          0,
		Timestamp:      testTime.Add(10 * time.Second),
		Interval:       10 * time.Second,
		AttributesJSON: json.RawMessage(`{"collector.name":"` + collectorName + `","first":"firstValue","instrumentation.provider":"` + instrumentationProvider + `","measure.name":"tests","measure.unit":"t","second":"secondValue","service.name":"serviceName"}`),
	}) {
		t.Errorf("metric fields are incorrect: %#v", metric)
	}
	if metric := h.metrics[2]; !reflect.DeepEqual(metric, telemetry.Count{
		Name:           "MyTestSum",
		Value:          5,
		Timestamp:      testTime.Add(20 * time.Second),
		Interval:       10 * time.Second,
		AttributesJSON: json.RawMessage(`{"collector.name":"` + collectorName + `","first":"firstValue","instrumentation.provider":"` + instrumentationProvider + `","measure.name":"tests","measure.unit":"t","second":"secondValue","service.name":"serviceName"}`),
	}) {
		t.Errorf("metric fields are incorrect: %#v", metric)
	}
}

func TestLastValueMetrics(t *testing.T) {
	h := &testHarvester{}
	exp := &Exporter{
		Harvester:       h,
		ServiceName:     "serviceName",
		DeltaCalculator: cumulative.NewDeltaCalculator(),
	}

	// first time metric is seen
	vd := &view.Data{
		View:  testLastValueView,
		Start: testTime,
		End:   testTime.Add(10 * time.Second),
		Rows: []*view.Row{
			&view.Row{
				Tags: []tag.Tag{
					tag.Tag{Key: testKeyFirst, Value: "firstValue"},
					tag.Tag{Key: testKeySecond, Value: "secondValue"},
				},
				Data: &view.LastValueData{
					Value: 10,
				},
			},
		},
	}
	exp.ExportView(vd)

	// second time metric is seen value does not change
	vd.End = testTime.Add(20 * time.Second)
	exp.ExportView(vd)

	// third time metric is seen value changes
	vd.End = testTime.Add(30 * time.Second)
	vd.Rows[0].Data.(*view.LastValueData).Value = 15
	exp.ExportView(vd)

	if metric := h.metrics[0]; !reflect.DeepEqual(metric, telemetry.Gauge{
		Name:      "MyTestLastValue",
		Value:     10,
		Timestamp: testTime.Add(10 * time.Second),
		Attributes: map[string]interface{}{
			"first":                    "firstValue",
			"second":                   "secondValue",
			"instrumentation.provider": instrumentationProvider,
			"collector.name":           collectorName,
			"measure.name":             "tests",
			"measure.unit":             "t",
			"service.name":             "serviceName",
		},
	}) {
		t.Errorf("metric fields are incorrect: %#v", metric)
	}
	if metric := h.metrics[1]; !reflect.DeepEqual(metric, telemetry.Gauge{
		Name:      "MyTestLastValue",
		Value:     10,
		Timestamp: testTime.Add(20 * time.Second),
		Attributes: map[string]interface{}{
			"first":                    "firstValue",
			"second":                   "secondValue",
			"instrumentation.provider": instrumentationProvider,
			"collector.name":           collectorName,
			"measure.name":             "tests",
			"measure.unit":             "t",
			"service.name":             "serviceName",
		},
	}) {
		t.Errorf("metric fields are incorrect: %#v", metric)
	}
	if metric := h.metrics[2]; !reflect.DeepEqual(metric, telemetry.Gauge{
		Name:      "MyTestLastValue",
		Value:     15,
		Timestamp: testTime.Add(30 * time.Second),
		Attributes: map[string]interface{}{
			"first":                    "firstValue",
			"second":                   "secondValue",
			"instrumentation.provider": instrumentationProvider,
			"collector.name":           collectorName,
			"measure.name":             "tests",
			"measure.unit":             "t",
			"service.name":             "serviceName",
		},
	}) {
		t.Errorf("metric fields are incorrect: %#v", metric)
	}
}

func TestDistributionMetrics(t *testing.T) {
	h := &testHarvester{}
	exp := &Exporter{
		Harvester:       h,
		ServiceName:     "serviceName",
		DeltaCalculator: cumulative.NewDeltaCalculator(),
	}

	vd := &view.Data{
		View:  testDistributionView,
		Start: testTime,
		End:   testTime.Add(10 * time.Second),
		Rows: []*view.Row{
			&view.Row{
				Tags: []tag.Tag{
					tag.Tag{Key: testKeyFirst, Value: "firstValue"},
					tag.Tag{Key: testKeySecond, Value: "secondValue"},
				},
				Data: &view.DistributionData{
					Count:           5,
					Min:             1,
					Max:             20000,
					Mean:            1234,
					SumOfSquaredDev: 123123123,
					CountPerBucket:  []int64{1, 2, 0, 0, 0, 0, 2},
				},
			},
		},
	}
	exp.ExportView(vd)

	if len(h.metrics) != 0 {
		t.Errorf("distribution type metrics not yet supported: %#v", h.metrics)
	}
}
