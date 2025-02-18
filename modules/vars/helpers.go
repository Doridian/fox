package vars

import lua "github.com/yuin/gopher-lua"

func Set(key string, value lua.LString) {
	varLock.Lock()
	varTable[key] = value
	varLock.Unlock()
}
