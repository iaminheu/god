package garray_test

import (
	"fmt"
	"git.zc0901.com/go/god/lib/container/garray"
	"testing"
)

func TestArray_Contains(t *testing.T) {
	ids := []int{1, 2, 3}
	var id int
	id = 1
	fmt.Println(garray.NewIntArrayFrom(ids).Contains(id))
}
