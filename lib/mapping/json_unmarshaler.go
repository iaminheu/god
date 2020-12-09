package mapping

import (
	"git.zc0901.com/go/god/lib/container/gmap"
	"git.zc0901.com/go/god/lib/jsonx"
	"io"
)

const jsonTagKey = "json"

var jsonUnmarshaler = NewUnmarshaler(jsonTagKey)

func UnmarshalJsonBytes(content []byte, v interface{}) error {
	return unmarshalJsonBytes(content, v, jsonUnmarshaler)
}

func UnmarshalJsonReader(reader io.Reader, v interface{}) (*gmap.StrAnyMap, error) {
	return unmarshalJsonReader(reader, v, jsonUnmarshaler)
}

func unmarshalJsonBytes(content []byte, v interface{}, unmarshaler *Unmarshaler) error {
	var m map[string]interface{}
	if err := jsonx.Unmarshal(content, &m); err != nil {
		return err
	}

	return unmarshaler.Unmarshal(m, v)
}

func unmarshalJsonReader(reader io.Reader, v interface{}, unmarshaler *Unmarshaler) (*gmap.StrAnyMap, error) {
	var m map[string]interface{}
	if err := jsonx.UnmarshalFromReader(reader, &m); err != nil {
		return nil, err
	}

	//// 转换
	//if err := gconv.Struct(m, v); err != nil {
	//	return err
	//}
	//// 验证
	//if err := gvalid.CheckStruct(v, nil); err != nil {
	//	return err.Current()
	//}

	return gmap.NewStrAnyMapFrom(m), nil

	//return unmarshaler.Unmarshal(m, v)
}
