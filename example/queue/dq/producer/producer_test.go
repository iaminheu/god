package main

import (
	"fmt"
	"testing"
	"time"
)

func TestNano(t *testing.T) {
	fmt.Println(time.Now().UnixNano())
}
