package handler

import (
	"bufio"
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"

	"git.zc0901.com/go/god/api/token"
	"git.zc0901.com/go/god/lib/logx"
	"github.com/dgrijalva/jwt-go"
)

const (
	jwtAudience    = "aud"
	jwtExpire      = "exp" // 过期时间
	jwtId          = "jti"
	jwtIssueAt     = "iat" // 签发时间
	jwtIssuer      = "iss" // 发行人
	jwtNotBefore   = "nbf"
	jwtSubject     = "sub"
	noDetailReason = "no detail reason"
)

var (
	errInvalidToken = errors.New("无效的鉴权令牌")
	errNoClaims     = errors.New("鉴权参数未提供")
)

type (
	AuthorizeOptions struct {
		PrevSecret string
		Callback   UnauthorizedCallback
	}

	UnauthorizedCallback func(w http.ResponseWriter, r *http.Request, err error)
	AuthorizeOption      func(opts *AuthorizeOptions)
)

// Authorize API JWT鉴权中间件
func Authorize(secret string, opts ...AuthorizeOption) func(http.Handler) http.Handler {
	var authOpts AuthorizeOptions
	for _, opt := range opts {
		opt(&authOpts)
	}

	parser := token.NewTokenParser()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwtToken, err := parser.ParseToken(r, secret, authOpts.PrevSecret)
			if err != nil {
				unauthorized(w, r, err, authOpts.Callback)
				return
			}

			if !jwtToken.Valid {
				unauthorized(w, r, errInvalidToken, authOpts.Callback)
				return
			}

			claims, ok := jwtToken.Claims.(jwt.MapClaims)
			if !ok {
				unauthorized(w, r, errNoClaims, authOpts.Callback)
				return
			}

			ctx := r.Context()
			for k, v := range claims {
				switch k {
				case jwtAudience, jwtExpire, jwtId, jwtIssueAt, jwtIssuer, jwtNotBefore, jwtSubject:
					// ignore the standard claims
				default:
					ctx = context.WithValue(ctx, k, v)
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WithPrevSecret(secret string) AuthorizeOption {
	return func(opts *AuthorizeOptions) {
		opts.PrevSecret = secret
	}
}

func WithUnauthorizedCallback(callback UnauthorizedCallback) AuthorizeOption {
	return func(opts *AuthorizeOptions) {
		opts.Callback = callback
	}
}

func detailAuthLog(r *http.Request, reason string) {
	// discard dump error, only for debug purpose
	details, _ := httputil.DumpRequest(r, true)
	logx.Errorf("鉴权失败：%s\n=> %+v", reason, string(details))
}

func unauthorized(w http.ResponseWriter, r *http.Request, err error, callback UnauthorizedCallback) {
	writer := newGuardedResponseWriter(w)

	if err != nil {
		detailAuthLog(r, err.Error())
	} else {
		detailAuthLog(r, noDetailReason)
	}
	if callback != nil {
		callback(writer, r, err)
	}

	writer.WriteHeader(http.StatusUnauthorized)
}

type guardedResponseWriter struct {
	writer      http.ResponseWriter
	wroteHeader bool
}

func newGuardedResponseWriter(w http.ResponseWriter) *guardedResponseWriter {
	return &guardedResponseWriter{
		writer: w,
	}
}

func (grw *guardedResponseWriter) Flush() {
	if flusher, ok := grw.writer.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (grw *guardedResponseWriter) Header() http.Header {
	return grw.writer.Header()
}

// Hijack 实现了 http.Hijacker 接口。
// 如果底层 http.ResponseWriter 支持的话此举将 Response 填充到 http.Hijacker。
func (grw *guardedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacked, ok := grw.writer.(http.Hijacker); ok {
		return hijacked.Hijack()
	}

	return nil, nil, errors.New("服务器不支持 hijack")
}

func (grw *guardedResponseWriter) Write(body []byte) (int, error) {
	return grw.writer.Write(body)
}

func (grw *guardedResponseWriter) WriteHeader(statusCode int) {
	if grw.wroteHeader {
		return
	}

	grw.wroteHeader = true
	grw.writer.WriteHeader(statusCode)
}
