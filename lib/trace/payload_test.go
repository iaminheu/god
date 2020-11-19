package trace

import (
	"git.zc0901.com/go/god/lib/stringx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpPayload(t *testing.T) {
	tests := []map[string]string{
		{},
		{
			"first":  "a",
			"second": "b",
		},
	}

	for _, test := range tests {
		t.Run(stringx.RandId(), func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
			payload := httpPayload(req.Header)
			for k, v := range test {
				payload.Set(k, v)
			}
			for k, v := range test {
				assert.Equal(t, v, payload.Get(k))
			}
			assert.Equal(t, "", payload.Get("none"))
		})
	}
}

func TestGrpcPayload(t *testing.T) {
	tests := []map[string]string{
		{},
		{
			"first":  "a",
			"second": "b",
		},
	}

	for _, test := range tests {
		t.Run(stringx.RandId(), func(t *testing.T) {
			m := make(map[string][]string)
			payload := grpcPayload(m)
			for k, v := range test {
				payload.Set(k, v)
			}
			for k, v := range test {
				assert.Equal(t, v, payload.Get(k))
			}
			assert.Equal(t, "", payload.Get("none"))
		})
	}
}
