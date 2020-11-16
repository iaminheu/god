package netx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternalIp(t *testing.T) {
	ip := InternalIp()
	assert.True(t, len(InternalIp()) > 0)
	fmt.Println(ip)
}
