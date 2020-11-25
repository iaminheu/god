package breaker

import (
	"fmt"
	"testing"
)

func TestNewErrorWindow(t *testing.T) {
	errWin := new(errorWindow)

	errWin.add("因为1....")
	errWin.add("因为2....")
	errWin.add("因为3....")

	fmt.Println(errWin.count)

	for _, reason := range errWin.reasons {
		fmt.Println(reason)
	}

	fmt.Println(errWin.String())
}
