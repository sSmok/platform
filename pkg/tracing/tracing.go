package tracing

import (
	"github.com/uber/jaeger-client-go/config"
)

// Init инициализирует трейсер
func Init(serviceName string) error {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: "localhost:6831",
		},
	}

	_, err := cfg.InitGlobalTracer(serviceName)
	if err != nil {
		return err
	}

	return nil
}
