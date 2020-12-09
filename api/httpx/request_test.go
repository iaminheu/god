package httpx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseForm(t *testing.T) {
	var v struct {
		Key      string `form:"key"`
		Location string `form:"location"`
		Radius   int64  `form:"radius" v:"required|min:1"`
		Offset   int64  `form:"offset"`
		Page     int64  `form:"page"`
		Phone    string `v:"required|phone#手机号必填|手机号格式不正确"`
	}

	r, e := http.NewRequest(http.MethodGet, "http://localhost:8888/place/around?key=6e10597c6b5f745d2ff915a4a721edfb&location=116.473168,39.993015&radius2=3000&extensions=base&output=json&offset=20&page=1", nil)
	if e = Parse(r, &v); e != nil {
		fmt.Println(e)
	}
	fmt.Println("key", v.Key)
	fmt.Println("location", v.Location)
	fmt.Println("radius", v.Radius)
	fmt.Println("offset", v.Offset)
	fmt.Println("page", v.Page)
}

func TestParseHeader(t *testing.T) {
	m := ParseHeader("key=value;")
	assert.EqualValues(t, map[string]string{
		"key": "value",
	}, m)
}

func TestParseJsonBody(t *testing.T) {
	var v struct {
		Name string `json:"name" v:"required"`
		Age  int    `json:"age"`
	}

	body := `{"age": 18}`
	r := httptest.NewRequest(http.MethodPost, "http://localhost:3333/", strings.NewReader(body))
	r.Header.Set(ContentType, ApplicationJson)

	if e := Parse(r, &v); e != nil {
		fmt.Println(e)
		fmt.Println("name", v.Name)
		fmt.Println("age", v.Age)
	}

	//assert.Nil(t, Parse(r, &v))
	//assert.Equal(t, "kevin", v.Name)
	//assert.Equal(t, 18, v.Age)
}
