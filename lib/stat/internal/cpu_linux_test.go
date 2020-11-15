package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRefreshCpu(t *testing.T) {
	assert.True(t, RefreshCpuUsage() > 0)
}

func BenchmarkRefreshCpu(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RefreshCpuUsage()
	}
}
