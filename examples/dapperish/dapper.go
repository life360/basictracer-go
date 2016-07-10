package dapperish

import (
	"github.com/life360/basictracer-go"
	"github.com/life360/opentracing-go"
)

// NewTracer returns a new dapperish Tracer instance.
func NewTracer(processName string) opentracing.Tracer {
	return basictracer.New(NewTrivialRecorder(processName))
}
