package sysx

import (
	"git.zc0901.com/go/god/lib/stringx"
	"os"
)

var hostname string

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		hostname = stringx.RandId()
	}
}

func Hostname() string {
	return hostname
}
