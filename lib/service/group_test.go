package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

var (
	number = 1
	mutex  sync.Mutex
	done   = make(chan struct{})
)

type mockedService struct {
	quit       chan struct{}
	multiplier int
}

func newMockedService(multiplier int) *mockedService {
	return &mockedService{
		quit:       make(chan struct{}),
		multiplier: multiplier,
	}
}

func (s *mockedService) Start() {
	mutex.Lock()
	number *= s.multiplier
	mutex.Unlock()
	done <- struct{}{}
	<-s.quit
}

func (s *mockedService) Stop() {
	close(s.quit)
}

func TestNewServiceGroup(t *testing.T) {
	multipliers := []int{2, 3, 5, 7}
	want := 1

	group := NewServiceGroup()
	for _, multiplier := range multipliers {
		want *= multiplier
		service := newMockedService(multiplier)
		group.Add(service)
	}

	go group.Start()

	for i := 0; i < len(multipliers); i++ {
		<-done
	}

	group.Stop()

	mutex.Lock()
	defer mutex.Unlock()
	assert.Equal(t, want, number)
	fmt.Println(want, number)
}
