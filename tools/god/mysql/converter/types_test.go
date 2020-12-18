package converter_test

import (
	"git.zc0901.com/go/god/tools/god/mysql/converter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertDataType(t *testing.T) {
	v, err := converter.ConvertDataType("tinyint", false)
	assert.Nil(t, err)
	assert.Equal(t, "int64", v)

	v, err = converter.ConvertDataType("tinyint", true)
	assert.Nil(t, err)
	assert.Equal(t, "sql.NullInt64", v)

	v, err = converter.ConvertDataType("timestamp", false)
	assert.Nil(t, err)
	assert.Equal(t, "time.Time", v)

	v, err = converter.ConvertDataType("timestamp", true)
	assert.Nil(t, err)
	assert.Equal(t, "sql.NullTime", v)
}
