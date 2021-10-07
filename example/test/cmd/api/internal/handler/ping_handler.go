package handler

import (
	"net/http"

	"git.zc0901.com/go/god/example/test/cmd/api/internal/logic"
	"git.zc0901.com/go/god/example/test/cmd/api/internal/svc"
	"git.zc0901.com/go/god/example/test/cmd/api/internal/types"

	"git.zc0901.com/go/god/api/httpx"
)

func PingHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PingReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewPingLogic(r.Context(), ctx)
		resp, err := l.Ping(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
