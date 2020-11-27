package conf

import (
	"fmt"
	"git.zc0901.com/go/god/lib/mapping"
	"io/ioutil"
	"log"
	"path"
)

var loaders = map[string]func([]byte, interface{}) error{
	".json": LoadConfigFromJsonBytes,
	".yaml": LoadConfigFromYamlBytes,
	".yml":  LoadConfigFromYamlBytes,
}

func LoadConfigFromJsonBytes(content []byte, v interface{}) error {
	return mapping.UnmarshalJsonBytes(content, v)
}

func LoadConfigFromYamlBytes(bytes []byte, v interface{}) error {
	return mapping.UnmarshalYamlBytes(bytes, v)
}

func LoadConfig(filename string, v interface{}) error {
	if content, err := ioutil.ReadFile(filename); err != nil {
		return err
	} else if loader, ok := loaders[path.Ext(filename)]; ok {
		return loader(content, v)
	} else {
		return fmt.Errorf("不识别的配置文件类型：%s", filename)
	}
}

func MustLoad(path string, v interface{}) {
	if err := LoadConfig(path, v); err != nil {
		log.Fatalf("错误：配置文件 - %s, %s", path, err.Error())
	}
}
