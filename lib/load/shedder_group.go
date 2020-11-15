package load

import (
	"git.zc0901.com/go/god/lib/syncx"
	"io"
)

type (
	ShedderGroup struct {
		options []ShedderOption
		manager *syncx.ResourceManager
	}

	nopCloser struct {
		Shedder
	}
)

func NewShedderGroup(opts ...ShedderOption) *ShedderGroup {
	return &ShedderGroup{
		options: opts,
		manager: syncx.NewResourceManager(),
	}
}

func (g *ShedderGroup) GetShedder(key string) Shedder {
	shedder, _ := g.manager.Get(key, func() (io.Closer, error) {
		return nopCloser{
			Shedder: NewAdaptiveShedder(g.options...),
		}, nil
	})
	return shedder.(Shedder)
}

func (c nopCloser) Close() error {
	return nil
}
