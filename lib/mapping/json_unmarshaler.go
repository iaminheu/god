package mapping

import (
	"git.zc0901.com/go/god/lib/gconv"
	"git.zc0901.com/go/god/lib/gvalid"
	"git.zc0901.com/go/god/lib/jsonx"
	"io"
)

const jsonTagKey = "json"

var jsonUnmarshaler = NewUnmarshaler(jsonTagKey)

func UnmarshalJsonBytes(content []byte, v interface{}) error {
	return unmarshalJsonBytes(content, v, jsonUnmarshaler)
}

func UnmarshalJsonReader(reader io.Reader, v interface{}) error {
	return unmarshalJsonReader(reader, v, jsonUnmarshaler)
}

func unmarshalJsonBytes(content []byte, v interface{}, unmarshaler *Unmarshaler) error {
	var m map[string]interface{}
	if err := jsonx.Unmarshal(content, &m); err != nil {
		return err
	}

	return unmarshaler.Unmarshal(m, v)
}

func unmarshalJsonReader(reader io.Reader, v interface{}, unmarshaler *Unmarshaler) error {
	var m map[string]interface{}
	if err := jsonx.UnmarshalFromReader(reader, &m); err != nil {
		return err
	}

	// 转换
	if err := gconv.Struct(m, v); err != nil {
		return err
	}
	// 验证
	if err := gvalid.CheckStruct(v, nil); err != nil {
		return err.Current()
	}

	return nil

	//return unmarshaler.Unmarshal(m, v)
}
