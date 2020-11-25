package service

import (
	"git.zc0901.com/go/god/lib/proc"
	"git.zc0901.com/go/god/lib/syncx"
	"git.zc0901.com/go/god/lib/threading"
	"log"
)

type (
	Starter interface {
		Start()
	}

	Stopper interface {
		Stop()
	}

	Service interface {
		Starter
		Stopper
	}

	Group struct {
		services []Service
		stopOnce func()
	}
)

func NewServiceGroup() *Group {
	g := new(Group)
	g.stopOnce = syncx.Once(g.doStop)
	return g
}

func (g *Group) Add(service Service) {
	g.services = append(g.services, service)
}

// 调用该方法后不应有任何逻辑代码，因为该方法是阻塞的，
// 同时，退出该方法后将关闭 logx 输出。
func (g *Group) Start() {
	proc.AddShutdownListener(func() {
		log.Println("服务关闭中...")
		g.stopOnce()
	})

	g.doStart()
}

func (g *Group) Stop() {
	g.stopOnce()
}

func (g *Group) doStart() {
	routineGroup := threading.NewRoutineGroup()

	for i := range g.services {
		service := g.services[i]
		routineGroup.RunSafe(func() {
			service.Start()
		})
	}

	routineGroup.Wait()
}

func (g *Group) doStop() {
	for _, service := range g.services {
		service.Stop()
	}
}

func WithStart(start func()) Service {
	return startOnlyService{
		start: start,
	}
}

func WithStarter(start Starter) Service {
	return starterOnlyService{
		Starter: start,
	}
}

type (
	stopper struct{}

	startOnlyService struct {
		start func()
		stopper
	}

	starterOnlyService struct {
		Starter
		stopper
	}
)

func (s stopper) Stop() {
}

func (s startOnlyService) Start() {
	s.start()
}
