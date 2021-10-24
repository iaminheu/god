package fx

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	const N = 5
	var count int32
	var wait sync.WaitGroup
	wait.Add(1)
	From(func(source chan<- interface{}) {
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()

		for i := 0; i < 2*N; i++ {
			select {
			case source <- i:
				fmt.Println("add", i)
				atomic.AddInt32(&count, 1)
			case <-ticker.C:
				wait.Done()
				return
			}
		}
	}).Buffer(N).ForAll(func(pipe <-chan interface{}) {
		wait.Wait()
		// 要多等一个，才能发送到通道
		assert.Equal(t, int32(N+1), atomic.LoadInt32(&count))
		fmt.Println(N+1, atomic.LoadInt32(&count))
	})
}

func TestBufferNegative(t *testing.T) {
	var result int
	Just(1, 2, 3, 4).Buffer(-1).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
		for item := range pipe {
			result += item.(int)
		}
		return result, nil
	})
	assert.Equal(t, 10, result)
}

func TestJust(t *testing.T) {
	var result int
	result2, err := Just(1, 2, 3, 4).Buffer(-1).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
		for item := range pipe {
			result += item.(int)
		}
		return result, nil
	})
	assert.Nil(t, err)
	fmt.Println(result)
	fmt.Println(result2)
}

func TestParallelJust(t *testing.T) {
	var count int32
	Just(1, 2, 3).Parallel(func(item interface{}) {
		time.Sleep(100 * time.Millisecond)
		atomic.AddInt32(&count, int32(item.(int)))
	}, UnlimitedWorkers())
	assert.Equal(t, int32(6), count)
}

func TestConvertVideo(t *testing.T) {
	cmd := exec.Command("ffmpeg", "-i", "/Users/zs/Desktop/video/guandian/75-如何改造我们的住宅.flv", "/Users/zs/Desktop/video/guandian/75-如何改造我们的住宅.mp4")

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("无法获得标准输出 %+v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("命令错误 %+v", err)
	}

	outputBuf := bufio.NewReader(stdoutPipe)
	for {
		output, _, err := outputBuf.ReadLine()
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Printf("错误: %s\n", err)
			}
			return
		}
		fmt.Printf("%s\n", string(output))

		if err := cmd.Wait(); err != nil {
			fmt.Print("等待：", err.Error())
		}
	}
}
