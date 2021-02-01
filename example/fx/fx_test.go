package main

import (
	"fmt"
	"git.zc0901.com/go/god/lib/fx"
	"testing"
)

func TestFxSplit(t *testing.T) {
	fx.Just(1, 2, 3, 4, 5).Split(2).ForEach(func(item interface{}) {
		vals := item.([]interface{})
		fmt.Println(vals)
	})
}
