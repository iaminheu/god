package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUUID(t *testing.T) {
	assert.Equal(t, 36, len(NewUUID()))
	fmt.Println(NewUUID())
}
