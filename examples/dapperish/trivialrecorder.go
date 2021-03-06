package dapperish

import (
	"fmt"
	"reflect"

	"github.com/life360/basictracer-go"
)

// TrivialRecorder implements the basictracer.Recorder interface.
type TrivialRecorder struct {
	processName string
	tags        map[string]string
}

// NewTrivialRecorder returns a TrivialRecorder for the given `processName`.
func NewTrivialRecorder(processName string) *TrivialRecorder {
	return &TrivialRecorder{
		processName: processName,
		tags:        make(map[string]string),
	}
}

// ProcessName returns the process name.
func (t *TrivialRecorder) ProcessName() string { return t.processName }

// SetTag sets a tag.
func (t *TrivialRecorder) SetTag(key string, val interface{}) *TrivialRecorder {
	t.tags[key] = fmt.Sprint(val)
	return t
}

// RecordSpan complies with the basictracer.Recorder interface.
func (t *TrivialRecorder) RecordSpan(span basictracer.RawSpan) {
	fmt.Printf(
		"RecordSpan: %v[%v, %v us] --> %v logs. std context: %v; baggage: %v\n",
		span.Operation, span.Start, span.Duration, len(span.Logs),
		span.SpanContext, span.Baggage)
	for i, l := range span.Logs {
		fmt.Printf(
			"    log %v @ %v: %v --> %v\n", i, l.Timestamp, l.Event, reflect.TypeOf(l.Payload))
	}
}
