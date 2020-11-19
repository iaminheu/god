package p2c

import (
	"context"
	"fmt"
	"git.zc0901.com/go/god/lib/logx"
	"git.zc0901.com/go/god/lib/mathx"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
	"runtime"
	"strconv"
	"sync"
	"testing"
)

func init() {
	logx.Disable()
}

func TestPicker_PickNil(t *testing.T) {
	builder := new(pickerBuilder)
	picker := builder.Build(nil)
	_, _, err := picker.Pick(context.Background(),
		balancer.PickInfo{
			FullMethodName: "/",
			Ctx:            context.Background(),
		})
	assert.NotNil(t, err)
	fmt.Println(err)
}

func TestPicker_Pick(t *testing.T) {
	tests := []struct {
		name       string
		candidates int
		threshold  float64
	}{
		{
			name:       "单个",
			candidates: 1,
			threshold:  0.9,
		},
		{
			name:       "两个",
			candidates: 2,
			threshold:  0.5,
		},
		{
			name:       "多个",
			candidates: 100,
			threshold:  0.95,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			const total = 10000
			builder := new(pickerBuilder)
			ready := make(map[resolver.Address]balancer.SubConn)
			for i := 0; i < test.candidates; i++ {
				ready[resolver.Address{Addr: strconv.Itoa(i)}] = new(mockClientConn)
			}

			picker := builder.Build(ready)
			var wg sync.WaitGroup
			wg.Add(total)
			for i := 0; i < total; i++ {
				_, done, err := picker.Pick(context.Background(), balancer.PickInfo{
					FullMethodName: "/",
					Ctx:            context.Background(),
				})
				assert.Nil(t, err)
				if i%100 == 0 {
					err = status.Error(codes.DeadlineExceeded, "超时啦")
				}
				go func() {
					runtime.Gosched()
					done(balancer.DoneInfo{Err: err})
					wg.Done()
				}()
			}
			wg.Wait()

			dist := make(map[interface{}]int)
			conns := picker.(*p2cPicker).conns
			for _, conn := range conns {
				dist[conn.addr.Addr] = int(conn.requests)
			}

			// 求熵
			entropy := mathx.CalcEntropy(dist)
			assert.True(t, entropy > test.threshold, fmt.Sprintf("熵：%f，小于：%f",
				entropy, test.threshold))
		})
	}
}

type mockClientConn struct {
}

func (m mockClientConn) UpdateAddresses(addresses []resolver.Address) {
}

func (m mockClientConn) Connect() {
}
