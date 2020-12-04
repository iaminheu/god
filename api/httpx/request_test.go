package httpx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestParseForm(t *testing.T) {
	var v struct {
		Key      string `form:"key"`
		Location string `form:"location"`
		Radius   int64  `form:"radius"`
		Offset   int64  `form:"offset"`
		Page     int64  `form:"page"`
	}

	r, err := http.NewRequest(http.MethodGet, "http://localhost:8888/place/around?key=6e10597c6b5f745d2ff915a4a721edfb&location=116.473168,39.993015&radius=3000&extensions=base&output=json&offset=20&page=1", nil)
	assert.Nil(t, err)
	assert.Nil(t, Parse(r, &v))
	fmt.Println(v)
}

func TestParseHeader(t *testing.T) {
	m := ParseHeader("key=value;")
	assert.EqualValues(t, map[string]string{
		"key": "value",
	}, m)
}
