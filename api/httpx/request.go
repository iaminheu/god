package httpx

import (
	"git.zc0901.com/go/god/api/internal/context"
	"git.zc0901.com/go/god/lib/mapping"
	"io"
	"net/http"
	"strings"
)

const (
	pathKey           = "path"
	formKey           = "form"
	maxMemory         = 32 << 20 // 32MB
	maxBodyLen        = 8 << 20  // 8MB
	emptyJson         = "{}"
	separator         = ";"
	tokensInAttribute = 2
)

var (
	// 路径参数解编排
	pathUnmarshaler = mapping.NewUnmarshaler(pathKey, mapping.WithStringValues())

	// 表单参数解编排
	formUnmarshaler = mapping.NewUnmarshaler(formKey, mapping.WithStringValues())
)

// 依次将请求路径、表单和JSON中的参数，解析值目标 v
func Parse(r *http.Request, v interface{}) error {
	if err := ParsePath(r, v); err != nil {
		return err
	}

	if err := ParseForm(r, v); err != nil {
		return err
	}

	return ParseJsonBody(r, v)
}

// 解析请求体为JSON的参数
func ParseJsonBody(r *http.Request, v interface{}) error {
	var reader io.Reader
	if withJsonBody(r) {
		reader = io.LimitReader(r.Body, maxBodyLen)
	} else {
		reader = strings.NewReader(emptyJson)
	}

	return mapping.UnmarshalJsonReader(reader, v)
}

// 解析表单请求参数（即Query参数）
func ParseForm(r *http.Request, v interface{}) error {
	if strings.Contains(r.Header.Get(ContentType), MultipartFormData) {
		if err := r.ParseMultipartForm(maxMemory); err != nil {
			return err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return err
		}
	}

	params := make(map[string]interface{}, len(r.Form))
	for key := range r.Form {
		value := r.Form.Get(key)
		if len(value) > 0 {
			params[key] = value
		}
	}

	return formUnmarshaler.Unmarshal(params, v)
}

// 解析URL中的路径参数
// 如：http://localhost/users/:name
func ParsePath(r *http.Request, v interface{}) error {
	vars := context.Vars(r)
	m := make(map[string]interface{}, len(vars))
	for k, v := range vars {
		m[k] = v
	}

	return pathUnmarshaler.Unmarshal(m, v)
}

func ParseHeader(headerValue string) map[string]string {
	m := make(map[string]string)
	fields := strings.Split(headerValue, separator)
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if len(field) == 0 {
			continue
		}

		kv := strings.SplitN(field, "=", tokensInAttribute)
		if len(kv) != tokensInAttribute {
			continue
		}

		m[kv[0]] = kv[1]
	}

	return m
}

// 判断是否带有JSON请求体
func withJsonBody(r *http.Request) bool {
	return r.ContentLength > 0 && strings.Contains(r.Header.Get(ContentType), ApplicationJson)
}
