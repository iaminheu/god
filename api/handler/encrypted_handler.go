package handler

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"git.zc0901.com/go/god/lib/codec"
	"git.zc0901.com/go/god/lib/logx"
	"io"
	"io/ioutil"
	"net"
	"net/http"
)

const maxBytes = 1 << 20 // 1MB

var errContentLengthExceeded = errors.New("内容超出长度限制")

// API 加密响应中间件
// key: 解密秘钥
func EncrpytedHandler(key []byte) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ew := newEncryptedResponseWriter(w)
			defer ew.flush(key)

			if r.ContentLength <= 0 {
				next.ServeHTTP(ew, r)
				return
			}

			if err := decryptBody(key, r); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			next.ServeHTTP(ew, r)
		})
	}
}

// 加密响应输出器
type encryptedResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (w *encryptedResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *encryptedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacked, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacked.Hijack()
	}

	return nil, nil, errors.New("server doesn't support hijacking")
}

func (w *encryptedResponseWriter) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

func (w *encryptedResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *encryptedResponseWriter) flush(key []byte) {
	if w.buf.Len() == 0 {
		return
	}

	content, err := codec.EcbEncrypt(key, w.buf.Bytes())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body := base64.StdEncoding.EncodeToString(content)
	if n, err := io.WriteString(w.ResponseWriter, body); err != nil {
		logx.Errorf("写响应失败，错误：%s", err)
	} else if n < len(content) {
		logx.Errorf("实际字节数: %d，写字节数: %d", len(content), n)
	}
}

func newEncryptedResponseWriter(w http.ResponseWriter) *encryptedResponseWriter {
	return &encryptedResponseWriter{
		ResponseWriter: w,
		buf:            new(bytes.Buffer),
	}
}

func decryptBody(key []byte, r *http.Request) error {
	if r.ContentLength > maxBytes {
		return errContentLengthExceeded
	}

	var content []byte
	var err error
	if r.ContentLength > 0 {
		content = make([]byte, r.ContentLength)
		_, err = io.ReadFull(r.Body, content)
	} else {
		content, err = ioutil.ReadAll(io.LimitReader(r.Body, maxBytes))
	}
	if err != nil {
		return err
	}

	content, err = base64.StdEncoding.DecodeString(string(content))
	if err != nil {
		return err
	}

	output, err := codec.EcbDecrypt(key, content)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.Write(output)
	r.Body = ioutil.NopCloser(&buf)

	return nil
}
