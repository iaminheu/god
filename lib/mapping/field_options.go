package mapping

import "fmt"

const notSymbol = '!'

type (
	// 使用context和OptionalDep选项来确定与context.Context无关的可选值
	fieldOptionsWithContext struct {
		FromString bool
		Optional   bool
		Options    []string
		Default    string
		Range      *numberRange
	}

	fieldOptions struct {
		fieldOptionsWithContext
		OptionalDep string
	}

	numberRange struct {
		left         float64
		leftInclude  bool
		right        float64
		rightInclude bool
	}
)

func (o *fieldOptionsWithContext) fromString() bool {
	return o != nil && o.FromString
}

func (o *fieldOptionsWithContext) getDefault() (string, bool) {
	if o == nil {
		return "", false
	} else {
		return o.Default, len(o.Default) > 0
	}
}

func (o *fieldOptionsWithContext) optional() bool {
	return o != nil && o.Optional
}

func (o *fieldOptionsWithContext) options() []string {
	if o == nil {
		return nil
	}

	return o.Options
}

func (o *fieldOptions) optionalDep() string {
	if o == nil {
		return ""
	} else {
		return o.OptionalDep
	}
}

func (o *fieldOptions) toOptionsWithContext(key string, m Valuer, fullName string) (
	*fieldOptionsWithContext, error) {
	var optional bool
	if o.optional() {
		dep := o.optionalDep()
		if len(dep) == 0 {
			optional = true
		} else if dep[0] == notSymbol {
			dep = dep[1:]
			if len(dep) == 0 {
				return nil, fmt.Errorf("wrong optional value for %q in %q", key, fullName)
			}

			_, baseOn := m.Value(dep)
			_, selfOn := m.Value(key)
			if baseOn == selfOn {
				return nil, fmt.Errorf("set value for either %q or %q in %q", dep, key, fullName)
			} else {
				optional = baseOn
			}
		} else {
			_, baseOn := m.Value(dep)
			_, selfOn := m.Value(key)
			if baseOn != selfOn {
				return nil, fmt.Errorf("values for %q and %q should be both provided or both not in %q",
					dep, key, fullName)
			} else {
				optional = !baseOn
			}
		}
	}

	if o.fieldOptionsWithContext.Optional == optional {
		return &o.fieldOptionsWithContext, nil
	} else {
		return &fieldOptionsWithContext{
			FromString: o.FromString,
			Optional:   optional,
			Options:    o.Options,
			Default:    o.Default,
		}, nil
	}
}
