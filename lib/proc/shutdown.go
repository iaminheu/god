package proc

import (
	"git.zc0901.com/go/god/lib/logx"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	shutdownListeners        = new(ListenerManager)
	wrapUpListener           = new(ListenerManager)
	delayTimeBeforeForceQuit = 5 * time.Second
)

// AddShutdownListener 添加一个程序关闭监听器
func AddShutdownListener(listener func()) (waitForCalled func()) {
	return shutdownListeners.add(listener)
}

// AddWrapUpListener 添加一个包装监听器
func AddWrapUpListener(listener func()) (waitForCalled func()) {
	return wrapUpListener.add(listener)
}

// gracefulStop 平滑停止程序（为关闭类监听器的执行留有时间）
func gracefulStop(signals chan os.Signal) {
	signal.Stop(signals)

	logx.Info("得到终止信号 SIGTERM，程序关闭中...")

	// 通知关闭类监听器，如定时执行器在关闭前进行任务 Flush
	shutdownListeners.notify()

	// 静候5秒，等待监听器处理完毕。
	time.Sleep(delayTimeBeforeForceQuit)

	// 关闭程序
	logx.Info("已等待 %v 秒，即将强制杀死该进程", delayTimeBeforeForceQuit)
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
}
