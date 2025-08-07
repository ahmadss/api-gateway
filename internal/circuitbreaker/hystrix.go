package circuitbreaker

import (
	"github.com/afex/hystrix-go/hystrix"
)

type CircuitBreaker struct {
	Timeout               int
	MaxConcurrentRequests int
	ErrorPercentThreshold int
}

func (cb *CircuitBreaker) Configure(name string) {
	hystrix.ConfigureCommand(name, hystrix.CommandConfig{
		Timeout:               cb.Timeout,
		MaxConcurrentRequests: cb.MaxConcurrentRequests,
		ErrorPercentThreshold: cb.ErrorPercentThreshold,
	})
}

func Do(name string, run func() error, fallback func(error) error) error {
	return hystrix.Do(name, run, fallback)
}
