package context

import (
	"context"
	"net/http"
)

// 路径变量键
var pathVars = contextKey("pathVars")

type contextKey string

func (c contextKey) String() string {
	return "api/internal/context key: " + string(c)
}

// 提取路径中的绑定变量
func Vars(r *http.Request) map[string]string {
	vars, ok := r.Context().Value(pathVars).(map[string]string)
	if ok {
		return vars
	}

	return nil
}

// 包装带有pathVars上下文
func WithPathVars(r *http.Request, params map[string]string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), pathVars, params))
}
