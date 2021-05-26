package sqlx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name  string
		query string
		args  []interface{}
		want  string
	}{
		{
			name:  "mysql 常规语句",
			query: "select name, age from users where bool=? and phone=?",
			args:  []interface{}{true, "133"},
			want:  "select name, age from users where bool=1 and phone='133'",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := format(test.query, test.args...)
			assert.Nil(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestDesensitize(t *testing.T) {
	datasource := "user:pass@tcp(111.222.333.44:3306)/any_table?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
	datasource = desensitize(datasource)
	assert.False(t, strings.Contains(datasource, "user"))
	assert.False(t, strings.Contains(datasource, "pass"))
	assert.True(t, strings.Contains(datasource, "tcp(111.222.333.44:3306)"))
}

func TestDesensitize_WithoutAccount(t *testing.T) {
	datasource := "tcp(111.222.333.44:3306)/any_table?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
	datasource = desensitize(datasource)
	assert.True(t, strings.Contains(datasource, "tcp(111.222.333.44:3306)"))
}

func TestEscape(t *testing.T) {
	s := "a\x00\n\r\\'\"\x1ab"

	out := escape(s)

	assert.Equal(t, `a\x00\n\r\\\'\"\x1ab`, out)
}
