package basictracer

import (
	"testing"

	"github.com/life360/opentracing-go/ext"
	"github.com/stretchr/testify/assert"
)

func TestSpan_Baggage(t *testing.T) {
	recorder := NewInMemoryRecorder()
	tracer := NewWithOptions(Options{
		Recorder:     recorder,
		ShouldSample: func(traceID uint64) bool { return true }, // always sample
	})
	span := tracer.StartSpan("x")
	span.Context().SetBaggageItem("x", "y")
	assert.Equal(t, "y", span.Context().BaggageItem("x"))
	span.Finish()
	spans := recorder.GetSpans()
	assert.Equal(t, 1, len(spans))
	assert.Equal(t, map[string]string{"x": "y"}, spans[0].Baggage)

	recorder.Reset()
	span = tracer.StartSpan("x")
	span.Context().SetBaggageItem("x", "y")
	baggage := make(map[string]string)
	span.Context().ForeachBaggageItem(func(k, v string) bool {
		baggage[k] = v
		return true
	})
	assert.Equal(t, map[string]string{"x": "y"}, baggage)

	span.Context().SetBaggageItem("a", "b")
	baggage = make(map[string]string)
	span.Context().ForeachBaggageItem(func(k, v string) bool {
		baggage[k] = v
		return false // exit early
	})
	assert.Equal(t, 1, len(baggage))
	span.Finish()
	spans = recorder.GetSpans()
	assert.Equal(t, 1, len(spans))
	assert.Equal(t, 2, len(spans[0].Baggage))
}

func TestSpan_Sampling(t *testing.T) {
	recorder := NewInMemoryRecorder()
	tracer := NewWithOptions(Options{
		Recorder:     recorder,
		ShouldSample: func(traceID uint64) bool { return true },
	})
	span := tracer.StartSpan("x")
	span.Finish()
	assert.Equal(t, 1, len(recorder.GetSampledSpans()), "by default span should be sampled")

	recorder.Reset()
	span = tracer.StartSpan("x")
	ext.SamplingPriority.Set(span, 0)
	span.Finish()
	assert.Equal(t, 0, len(recorder.GetSampledSpans()), "SamplingPriority=0 should turn off sampling")

	tracer = NewWithOptions(Options{
		Recorder:     recorder,
		ShouldSample: func(traceID uint64) bool { return false },
	})

	recorder.Reset()
	span = tracer.StartSpan("x")
	span.Finish()
	assert.Equal(t, 0, len(recorder.GetSampledSpans()), "by default span should not be sampled")

	recorder.Reset()
	span = tracer.StartSpan("x")
	ext.SamplingPriority.Set(span, 1)
	span.Finish()
	assert.Equal(t, 1, len(recorder.GetSampledSpans()), "SamplingPriority=1 should turn on sampling")
}
