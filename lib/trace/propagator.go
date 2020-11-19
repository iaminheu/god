package trace

import (
	"google.golang.org/grpc/metadata"
	"net/http"
)

const (
	HttpFormat = iota
	GrpcFormat
)

var (
	emptyHttpPropagator httpPropagator
	emptyGrpcPropagator grpcPropagator
)

type (
	Propagator interface {
		Extract(payload interface{}) (Payload, error)
		Inject(payload interface{}) (Payload, error)
	}

	httpPropagator struct{}
	grpcPropagator struct{}
)

func (h httpPropagator) Extract(payload interface{}) (Payload, error) {
	if p, ok := payload.(http.Header); ok {
		return httpPayload(p), nil
	} else {
		return nil, ErrInvalidPayload
	}
}

func (h httpPropagator) Inject(payload interface{}) (Payload, error) {
	if p, ok := payload.(http.Header); ok {
		return httpPayload(p), nil
	} else {
		return nil, ErrInvalidPayload
	}
}

func (h grpcPropagator) Extract(payload interface{}) (Payload, error) {
	if p, ok := payload.(metadata.MD); ok {
		return grpcPayload(p), nil
	} else {
		return nil, ErrInvalidPayload
	}
}

func (h grpcPropagator) Inject(payload interface{}) (Payload, error) {
	if p, ok := payload.(metadata.MD); ok {
		return grpcPayload(p), nil
	} else {
		return nil, ErrInvalidPayload
	}
}

// Extract 提取指定格式的负载参数
func Extract(format, payload interface{}) (Payload, error) {
	switch v := format.(type) {
	case int:
		if v == HttpFormat {
			return emptyHttpPropagator.Extract(payload)
		} else if v == GrpcFormat {
			return emptyGrpcPropagator.Extract(payload)
		}
	}

	return nil, ErrInvalidPayload
}

// Inject 注入指定格式的负载参数
func Inject(format, payload interface{}) (Payload, error) {
	switch v := format.(type) {
	case int:
		if v == HttpFormat {
			return emptyHttpPropagator.Inject(payload)
		} else if v == GrpcFormat {
			return emptyGrpcPropagator.Inject(payload)
		}
	}

	return nil, ErrInvalidPayload
}
