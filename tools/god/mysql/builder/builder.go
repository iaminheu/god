package builder

import (
	"fmt"
	"reflect"

	"github.com/go-xorm/builder"
)

const dbTag = "db"

func NewEq(in interface{}) builder.Eq {
	return builder.Eq(ToMap(in))
}

func NewGt(in interface{}) builder.Gt {
	return builder.Gt(ToMap(in))
}

func ToMap(in interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// we only accept structs
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
	}
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(dbTag); tagv != "" {
			// set key of map to value in struct field
			val := v.Field(i)
			zero := reflect.Zero(val.Type()).Interface()
			current := val.Interface()

			if reflect.DeepEqual(current, zero) {
				continue
			}
			out[tagv] = current
		}
	}
	return out
}

// FieldList 返回表的字段列表
func FieldList(table interface{}) []string {
	out := make([]string, 0)
	v := reflect.ValueOf(table)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// we only accept structs
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
	}
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(dbTag); tagv != "" {
			out = append(out, fmt.Sprintf("`%v`", tagv))
		} else {
			out = append(out, fmt.Sprintf("`%v`", fi.Name))
		}
	}
	return out
}

func FieldListAlias(in interface{}, alias string) []string {
	out := make([]string, 0)
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// we only accept structs
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
	}
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		tagName := ""
		if tagv := fi.Tag.Get(dbTag); tagv != "" {
			tagName = tagv
		} else {
			tagName = fi.Name
		}
		if len(alias) > 0 {
			tagName = alias + "." + tagName
		}
		out = append(out, tagName)
	}
	return out
}
