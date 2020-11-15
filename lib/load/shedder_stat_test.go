package load

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewShedderStat(t *testing.T) {
	ss := NewShedderStat("any")
	for i := 0; i < 3; i++ {
		ss.IncrTotal()
	}
	for i := 0; i < 5; i++ {
		ss.IncrPass()
	}
	for i := 0; i < 7; i++ {
		ss.IncrDrop()
	}
	result := ss.reset()
	assert.Equal(t, int64(3), result.Total)
	assert.Equal(t, int64(5), result.Pass)
	assert.Equal(t, int64(7), result.Drop)
}
