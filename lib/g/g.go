package g

import "git.zc0901.com/go/god/lib/container/gvar"

// Var是通用变量接口，与泛型类似。
type Var = gvar.Var

// 常用的 Map 类型别名。
type Map = map[string]interface{}
type MapAnyAny = map[interface{}]interface{}
type MapAnyStr = map[interface{}]string
type MapAnyInt = map[interface{}]int
type MapStrAny = map[string]interface{}
type MapStrStr = map[string]string
type MapStrInt = map[string]int
type MapIntAny = map[int]interface{}
type MapIntStr = map[int]string
type MapIntInt = map[int]int
type MapAnyBool = map[interface{}]bool
type MapStrBool = map[string]bool
type MapIntBool = map[int]bool

// 常用的 List 类型别名。
type List = []Map
type ListAnyAny = []Map
type ListAnyStr = []MapAnyStr
type ListAnyInt = []MapAnyInt
type ListStrAny = []MapStrAny
type ListStrStr = []MapStrStr
type ListStrInt = []MapStrInt
type ListIntAny = []MapIntAny
type ListIntStr = []MapIntStr
type ListIntInt = []MapIntInt
type ListAnyBool = []MapAnyBool
type ListStrBool = []MapStrBool
type ListIntBool = []MapIntBool

// 常用的 Slice 类型别名。
type Slice = []interface{}
type SliceAny = []interface{}
type SliceStr = []string
type SliceInt = []int

// Array 是 Slice 别名。
type Array = []interface{}
type ArrayAny = []interface{}
type ArrayStr = []string
type ArrayInt = []int
