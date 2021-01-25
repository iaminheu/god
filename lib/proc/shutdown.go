// +build linux darwin

package proc

import (
	"git.zc0901.com/go/god/lib/logx"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	wrapUpTime = time.Second
	// why we use 5500 milliseconds is because most of our queue are blocking mode with 5 seconds
	waitTime = 5500 * time.Millisecond
)

var (
	shutdownListeners        = new(ListenerManager)
	wrapUpListeners          = new(ListenerManager)
	delayTimeBeforeForceQuit = waitTime
)

// AddShutdownListener 添加一个程序关闭监听器
func AddShutdownListener(listener func()) (waitForCalled func()) {
	return shutdownListeners.add(listener)
}

// AddWrapUpListener 添加一个包装监听器
func AddWrapUpListener(listener func()) (waitForCalled func()) {
	return wrapUpListeners.add(listener)
}

// gracefulStop 平滑停止程序（为关闭类监听器的执行留有时间）
func gracefulStop(signals chan os.Signal) {
	signal.Stop(signals)

	logx.Info("捕获终止信号 SIGTERM，程序关闭中...")
	wrapUpListeners.notify()

	// 通知关闭类监听器，如定时执行器在关闭前进行任务 Flush
	time.Sleep(wrapUpTime)
	shutdownListeners.notify()

	// 静候5秒，等待监听器处理完毕，然后关闭程序
	time.Sleep(delayTimeBeforeForceQuit - wrapUpTime)
	logx.Info("已等待 %v 秒，即将强制杀死该进程", delayTimeBeforeForceQuit)
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
}
