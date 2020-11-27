package httpx

import (
	"encoding/json"
	"git.zc0901.com/go/god/lib/logx"
	"net/http"
	"sync"
)

var (
	errorHandler func(error) (int, interface{})
	lock         sync.RWMutex
)

// 错误响应，支持自定义错误处理器
func Error(w http.ResponseWriter, err error) {
	lock.Lock()
	handler := errorHandler
	lock.RUnlock()

	if handler == nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	code, body := errorHandler(err)
	e, ok := body.(error)
	if ok {
		http.Error(w, e.Error(), code)
	} else {
		WriteJson(w, code, body)
	}
}

// 正常响应
func Ok(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

// 正常JSON响应
func OkJson(w http.ResponseWriter, body interface{}) {
	WriteJson(w, http.StatusOK, body)
}

// 设置自定义错误处理器
func SetErrorHandler(handler func(error) (int, interface{})) {
	lock.Lock()
	defer lock.Unlock()
	errorHandler = handler
}

// 写JSON响应
func WriteJson(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set(ContentType, ApplicationJson)
	w.WriteHeader(code)

	if bytes, err := json.Marshal(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if n, err := w.Write(bytes); err != nil {
		// http.ErrHandlerTimeout 已经被 http.TimeoutHandler 处理了
		// 所以此处忽略。
		if err != http.ErrHandlerTimeout {
			logx.Errorf("写响应失败，错误：%s", err)
		}
	} else if n < len(bytes) {
		logx.Errorf("实际字节数：%d，写字节数：%d", len(bytes), n)
	}
}