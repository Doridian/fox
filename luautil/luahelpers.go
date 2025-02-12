package luautil

import (
	lua "github.com/yuin/gopher-lua"
)

func MergeFuncMaps(maps ...map[string]lua.LGFunction) map[string]lua.LGFunction {
	ret := make(map[string]lua.LGFunction)

	for _, m := range maps {
		for k, v := range m {
			ret[k] = v
		}
	}

	return ret
}
