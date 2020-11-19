package discovery

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEtcdConf_Validate(t *testing.T) {
	tests := []struct {
		EtcdConf
		pass bool
	}{
		{
			EtcdConf{},
			false,
		},
		{
			EtcdConf{
				Key: "any",
			},
			false,
		},
		{
			EtcdConf{
				Hosts: []string{"any"},
			},
			false,
		},
		{
			EtcdConf{
				Hosts: []string{"any"},
				Key:   "key",
			},
			true,
		},
	}

	for _, test := range tests {
		if test.pass {
			assert.Nil(t, test.EtcdConf.Validate())
		} else {
			assert.NotNil(t, test.EtcdConf.Validate())
		}
	}
}
