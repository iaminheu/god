package prometheus

import (
	"fmt"
	"git.zc0901.com/go/god/lib/logx"
	"git.zc0901.com/go/god/lib/threading"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"
)

var once sync.Once

// StartAgent 启动普罗米修斯Http代理服务
func StartAgent(c PromConf) {
	once.Do(func() {
		if len(c.Host) == 0 {
			return
		}

		threading.GoSafe(func() {
			http.Handle(c.Path, promhttp.Handler())
			addr := fmt.Sprintf("%s:%s", c.Host, c.Port)
			logx.Infof("启动普罗米修斯代理，端口：%s", addr)
			if err := http.ListenAndServe(addr, nil); err != nil {
				logx.Error(err)
			}
		})
	})
}
