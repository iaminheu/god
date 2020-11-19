package discovery

import "errors"

type EtcdConf struct {
	Hosts []string
	Key   string
}

func (c EtcdConf) Validate() error {
	if len(c.Hosts) == 0 {
		return errors.New("未配置用于服务发现的 Etcd Hosts")
	} else if len(c.Key) == 0 {
		return errors.New("未配置用于服务发现的 Etcd Key")
	} else {
		return nil
	}
}
