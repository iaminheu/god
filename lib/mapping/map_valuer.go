package mapping

type (
	Valuer interface {
		Value(key string) (interface{}, bool)
	}

	MapStrAny map[string]interface{}
)

func (m MapStrAny) Value(key string) (interface{}, bool) {
	v, ok := m[key]
	return v, ok
}
