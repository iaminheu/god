package mathx

import (
	"fmt"
	"testing"
)

func TestProb_TrueOnProb(t *testing.T) {
	prob := NewProba()
	fmt.Println(prob.TrueOnProba(0.3))
}
