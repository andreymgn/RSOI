package tracer

import (
	"fmt"
	"io"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

// NewTracer returns new tracer
func NewTracer(serviceName, host string) (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  host,
		},
	}

	tracer, closer, err := cfg.New(serviceName)
	if err != nil {
		return nil, nil, fmt.Errorf("new tracer error: %v", err)
	}
	return tracer, closer, nil
}
