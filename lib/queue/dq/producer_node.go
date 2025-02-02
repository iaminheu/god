package dq

import (
	"errors"
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
	"strconv"
	"strings"
	"time"
)

var ErrTimeBeforeNow = errors.New("不能把任务安排到过去的时间")

type producerNode struct {
	endpoint string
	tube     string
	conn     *connection
}

// 新建生产者节点
func NewProducerNode(endpoint, tube string) Producer {
	return &producerNode{
		endpoint: endpoint,
		tube:     tube,
		conn:     newConnection(endpoint, tube),
	}
}

// 定时执行
func (p producerNode) At(body []byte, at time.Time) (string, error) {
	now := time.Now()
	if at.Before(now) {
		return "", ErrTimeBeforeNow
	}

	delay := at.Sub(now)
	return p.Delay(body, delay)
}

// 延迟执行
func (p producerNode) Delay(body []byte, delay time.Duration) (string, error) {
	conn, err := p.conn.get()
	if err != nil {
		return "", err
	}

	id, err := conn.Put(body, PriorityNormal, delay, defaultTimeToRun)

	// 推送成功
	if err == nil {
		return fmt.Sprintf("%s/%s/%d", p.endpoint, p.tube, id), nil
	}

	// 推送失败
	switch e := err.(type) {
	case beanstalk.ConnError:
		switch e.Err {
		case beanstalk.ErrBadChar, beanstalk.ErrBadFormat, beanstalk.ErrBuried, beanstalk.ErrDeadline,
			beanstalk.ErrDraining, beanstalk.ErrEmpty, beanstalk.ErrInternal, beanstalk.ErrJobTooBig,
			beanstalk.ErrNoCRLF, beanstalk.ErrNotFound, beanstalk.ErrNotIgnored, beanstalk.ErrTooLong:
		// 不重置连接
		default:
			// 重置连接的错误类型：
			// beanstalk.ErrOOM, beanstalk.ErrTimeout, beanstalk.ErrUnknown 和其他错误。
			p.conn.reset()
		}
	}

	return "", err
}

//  撤回一批任务
//
// ids: endpoint/tube/id,endpoint/tube/id,endpoint/tube/id
func (p producerNode) Revoke(jointId string) error {
	ids := strings.Split(jointId, idSep)
	for _, id := range ids {
		fields := strings.Split(id, "/")
		if len(fields) < 3 {
			continue
		}
		if fields[0] != p.endpoint || fields[1] != p.tube {
			continue
		}

		conn, err := p.conn.get()
		if err != nil {
			return err
		}

		n, err := strconv.ParseUint(fields[2], 10, 64)
		if err != nil {
			return err
		}

		return conn.Delete(n)
	}

	// 如果任务不再该连接中，则忽略。
	return nil
}

// 关闭该集群节点
func (p producerNode) Close() error {
	return p.conn.Close()
}
