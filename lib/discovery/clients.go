package discovery

import (
	"fmt"
	"git.zc0901.com/go/god/lib/discovery/internal"
	"strings"
)

const timeToLive int64 = 10

var TimeToLive = timeToLive

func makeEtcdKey(key string, id int64) string {
	return fmt.Sprintf("%s%c%d", key, internal.Delimiter, id)
}

func extract(etcdKey string, index int) (string, bool) {
	if index < 0 {
		return "", false
	}

	fields := strings.FieldsFunc(etcdKey, func(ch rune) bool {
		return ch == internal.Delimiter
	})
	if index >= len(fields) {
		return "", false
	}

	return fields[index], true
}
