package gi18n

import "git.zc0901.com/go/god/lib/container/gmap"

const (
	// Default group name for instance usage.
	DefaultName = "default"
)

var (
	// instances is the instances map for management
	// for multiple i18n instance by name.
	instances = gmap.NewStrAnyMap(true)
)

// Instance returns an instance of Resource.
// The parameter <name> is the name for the instance.
func Instance(name ...string) *Manager {
	key := DefaultName
	if len(name) > 0 && name[0] != "" {
		key = name[0]
	}
	return instances.GetOrSetFuncLock(key, func() interface{} {
		return New()
	}).(*Manager)
}
