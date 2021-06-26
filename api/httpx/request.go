package httpx

import (
	"io"
	"net/http"
	"strings"

	"git.zc0901.com/go/god/api/internal/context"
	"git.zc0901.com/go/god/lib/container/gmap"
	"git.zc0901.com/go/god/lib/gconv"
	"git.zc0901.com/go/god/lib/gvalid"
	"git.zc0901.com/go/god/lib/mapping"
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

// Parse 依次将请求路径、表单和JSON中的参数，解析值目标 v
func Parse(r *http.Request, pointer interface{}) error {
	pathParams, err := ParsePath(r, pointer)
	if err != nil {
		return err
	}

	formParams, err := ParseForm(r, pointer)
	if err != nil {
		return err
	}

	bodyParams, err := ParseJsonBody(r, pointer)
	if err != nil {
		return err
	}

	params := pathParams.Clone()
	params.Merge(formParams)
	params.Merge(bodyParams)

	// 转换
	if err := gconv.Struct(params, pointer); err != nil {
		return err
	}
	// 验证
	if err := gvalid.CheckStruct(pointer, nil); err != nil {
		return err.Current()
	}

	return nil
}

// ParseJsonBody 解析请求体为JSON的参数
func ParseJsonBody(r *http.Request, pointer interface{}) (*gmap.StrAnyMap, error) {
	var reader io.Reader
	if withJsonBody(r) {
		reader = io.LimitReader(r.Body, maxBodyLen)
	} else {
		reader = strings.NewReader(emptyJson)
	}

	return mapping.UnmarshalJsonReader(reader, pointer)
}

// ParseForm 解析表单请求参数（即Query参数）
func ParseForm(r *http.Request, pointer interface{}) (*gmap.StrAnyMap, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		if err != http.ErrNotMultipart {
			return nil, err
		}
	}

	params := make(map[string]interface{}, len(r.Form))
	for key := range r.Form {
		value := r.Form.Get(key)
		if len(value) > 0 {
			params[key] = value
		}
	}

	//// 转换
	//if err := gconv.Struct(params, pointer); err != nil {
	//	return err
	//}
	//// 验证
	//if err := gvalid.CheckStruct(pointer, nil); err != nil {
	//	return err.Current()
	//}

	return gmap.NewStrAnyMapFrom(params), nil
	// return formUnmarshaler.Unmarshal(params, pointer)
}

// ParsePath 解析URL中的路径参数。
// 如：http://localhost/users/:name
func ParsePath(r *http.Request, pointer interface{}) (*gmap.StrAnyMap, error) {
	vars := context.Vars(r)
	params := make(map[string]interface{}, len(vars))
	for k, v := range vars {
		params[k] = v
	}

	//// 转换
	//if err := gconv.Struct(params, pointer); err != nil {
	//	return err
	//}
	//// 验证
	//if err := gvalid.CheckStruct(pointer, nil); err != nil {
	//	return err.Current()
	//}

	return gmap.NewStrAnyMapFrom(params), nil
	// return pathUnmarshaler.Unmarshal(params, pointer)
}

func ParseHeader(headerValue string) map[string]string {
	params := make(map[string]string)
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

		params[kv[0]] = kv[1]
	}

	return params
}

// 判断是否带有JSON请求体
func withJsonBody(r *http.Request) bool {
	return r.ContentLength > 0 && strings.Contains(r.Header.Get(ContentType), ApplicationJson)
}
