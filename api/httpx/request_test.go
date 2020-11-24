package httpx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestParseForm(t *testing.T) {
	var v struct {
		Name    string  `form:"name"`
		Age     int     `form:"age"`
		Percent float64 `form:"percent,optional"`
	}

	r, err := http.NewRequest(http.MethodGet, "http://hello.com/a?name=hello&age=18&percent=3.4", nil)
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
