package parser

import "git.zc0901.com/go/god/tools/god/api/spec"

type state interface {
	process(api *spec.ApiSpec) (state, error)
}
