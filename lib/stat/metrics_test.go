package stat

import "testing"

func TestNewMetrics(t *testing.T) {
	counts := []int{1, 5, 10, 100, 1000, 1000}
	for _, count := range counts {
		m := NewMetrics("foo")
		m.SetName("bar")
		fori
	}
}
