package trace

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpPropagator_Extract(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.Header.Set(traceIdKey, "trace")
	req.Header.Set(spanIdKey, "span")
	payload, err := Extract(HttpFormat, req.Header)
	assert.Nil(t, err)
	assert.Equal(t, "trace", payload.Get(traceIdKey))
	assert.Equal(t, "span", payload.Get(spanIdKey))

	_, err = Extract(HttpFormat, req)
	assert.Equal(t, ErrInvalidPayload, err)
}

func TestHttpPropagator_Inject(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.Header.Set(traceIdKey, "trace")
	req.Header.Set(spanIdKey, "span")
	payload, err := Inject(HttpFormat, req.Header)
	assert.Nil(t, err)
	assert.Equal(t, "trace", payload.Get(traceIdKey))
	assert.Equal(t, "span", payload.Get(spanIdKey))

	_, err = Inject(HttpFormat, req)
	assert.Equal(t, ErrInvalidPayload, err)
}

func TestGrpcPropagator_Extract(t *testing.T) {
	md := metadata.New(map[string]string{
		traceIdKey: "trace",
		spanIdKey:  "span",
	})
	payload, err := Extract(GrpcFormat, md)
	assert.Nil(t, err)
	assert.Equal(t, "trace", payload.Get(traceIdKey))
	assert.Equal(t, "span", payload.Get(spanIdKey))

	_, err = Extract(GrpcFormat, 1)
	assert.Equal(t, ErrInvalidPayload, err)
	_, err = Extract(nil, 1)
	assert.Equal(t, ErrInvalidPayload, err)
}

func TestGrpcPropagator_Inject(t *testing.T) {
	md := metadata.New(map[string]string{
		traceIdKey: "trace",
		spanIdKey:  "span",
	})
	payload, err := Inject(GrpcFormat, md)
	assert.Nil(t, err)
	assert.Equal(t, "trace", payload.Get(traceIdKey))
	assert.Equal(t, "span", payload.Get(spanIdKey))

	_, err = Inject(GrpcFormat, 1)
	assert.Equal(t, ErrInvalidPayload, err)
	_, err = Inject(nil, 1)
	assert.Equal(t, ErrInvalidPayload, err)
}
