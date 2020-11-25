package stat

import (
	"bytes"
	"encoding/json"
	"errors"
	"git.zc0901.com/go/god/lib/logx"
	"net/http"
	"time"
)

const httpTimeout = time.Second * 5

var ErrWriteFailed = errors.New("远程统计提交错误")

type RemoteWriter struct {
	endpoint string
}

// 新建远程上报器
func NewRemoteWriter(endpoint string) Writer {
	return &RemoteWriter{endpoint}
}

func (rw RemoteWriter) Write(report *StatReport) error {
	bs, err := json.Marshal(report)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: httpTimeout}
	// endpoint 就是推送到prometheus server的地址
	resp, err := client.Post(rw.endpoint, "application/json", bytes.NewReader(bs))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logx.Errorf("提交统计报告失败，错误码：%d, 原因：%s", resp.StatusCode, resp.Status)
		return ErrWriteFailed
	}

	return nil
}
