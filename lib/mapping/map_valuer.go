package mapping

type (
	// Valuer 接口定义了从底层对象使用指定键获取值的方式。
	Valuer interface {
		// Value 获取指定键的值。
		Value(key string) (interface{}, bool)
	}

	// MapValuer MapValuer 是一个使用 Value 方法获取指定键的值的方法。
	MapValuer map[string]interface{}
)

// Value 从 mv 中获取指定键的值。
func (mv MapValuer) Value(key string) (interface{}, bool) {
	v, ok := mv[key]
	return v, ok
}
