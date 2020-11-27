package handler

import (
	"net/http"

	"git.zc0901.com/go/god/api/httpx"
	"git.zc0901.com/go/god/example/graceful/dns/api/internal/logic"
	"git.zc0901.com/go/god/example/graceful/dns/api/internal/svc"
)

func GracefulHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewGracefulLogic(r.Context(), ctx)
		resp, err := l.Graceful()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
