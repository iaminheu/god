package search

import "errors"

const (
	colon = ':'
	slash = '/'
)

var (
	ErrDuplicateItem  = errors.New("重复的项")
	ErrDuplicateSlash = errors.New("重复的 /")
	ErrEmptyItem      = errors.New("不能为空项")
	ErrInvalidState   = errors.New("搜索树处在无效状态")
	ErrNotFromRoot    = errors.New("路径必须以 / 开头")
	NotFound          Result // 未找到
)

type (
	Tree struct {
		root *node
	}

	node struct {
		item     interface{}
		children [2]map[string]*node
	}

	Result struct {
		Item   interface{}
		Params map[string]string
	}

	innerResult struct {
		key   string
		value string
		named bool
		found bool
	}
)

func NewTree() *Tree {
	return &Tree{root: newNode(nil)}
}

func (t *Tree) Add(route string, item interface{}) error {
	if len(route) == 0 || route[0] != slash {
		return ErrNotFromRoot
	}

	if item == nil {
		return ErrEmptyItem
	}

	return add(t.root, route[1:], item)
}

func (t *Tree) Search(route string) (Result, bool) {
	if len(route) == 0 || route[0] != slash {
		return NotFound, false
	}

	var result Result
	ok := t.next(t.root, route[1:], &result)
	return result, ok
}

func (t *Tree) next(n *node, route string, result *Result) bool {
	if len(route) == 0 && n.item != nil {
		result.Item = n.item
		return true
	}

	for i := range route {
		if route[i] == slash {
			token := route[:i]
			for _, children := range n.children {
				for k, v := range children {
					if r := match(k, token); r.found {
						if t.next(v, route[i+1:], result) {
							if r.named {
								addParam(result, r.key, r.value)
							}

							return true
						}
					}
				}
			}

			return false
		}
	}

	for _, children := range n.children {
		for k, v := range children {
			if r := match(k, route); r.found && v.item != nil {
				result.Item = v.item
				if r.named {
					addParam(result, r.key, r.value)
				}

				return true
			}
		}
	}

	return false
}

// 添加结果参数
func addParam(result *Result, key, value string) {
	if result.Params == nil {
		result.Params = make(map[string]string)
	}

	result.Params[key] = value
}

func match(pat string, token string) innerResult {
	// 冒号开头
	if pat[0] == colon {
		return innerResult{
			key:   pat[1:],
			value: token,
			named: true,
			found: true,
		}
	}

	// 非冒号开头
	return innerResult{found: pat == token}
}

func add(n *node, route string, item interface{}) error {
	// TODO 添加根目录？
	if len(route) == 0 {
		if n.item != nil {
			return ErrDuplicateItem
		}

		n.item = item
		return nil
	}

	if route[0] == slash {
		return ErrDuplicateSlash
	}

	for i := range route {
		if route[i] == slash {
			token := route[:i]
			children := n.getChildren(token)
			if child, ok := children[token]; ok {
				if child != nil {
					return add(child, route[i+1:], item)
				} else {
					return ErrInvalidState
				}
			} else {
				child := newNode(nil)
				children[token] = child
				return add(child, route[i+1:], item)
			}
		}
	}

	children := n.getChildren(route)
	if child, ok := children[route]; ok {
		if child.item != nil {
			return ErrDuplicateItem
		}

		child.item = item
	} else {
		children[route] = newNode(item)
	}

	return nil
}

func newNode(item interface{}) *node {
	return &node{
		item: item,
		children: [2]map[string]*node{
			make(map[string]*node),
			make(map[string]*node),
		},
	}
}

func (n *node) getChildren(route string) map[string]*node {
	if len(route) > 0 && route[0] == colon {
		return n.children[1]
	} else {
		return n.children[0]
	}
}
