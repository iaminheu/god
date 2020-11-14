package mapping

type (
	Valuer interface {
		Value(key string) (interface{}, bool)
	}

	MapValuer map[string]interface{}
)

func (m MapValuer) Value(key string) (interface{}, bool) {
	v, ok := m[key]
	return v, ok
}
