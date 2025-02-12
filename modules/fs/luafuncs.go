package fs

import (
	"os"

	"github.com/Doridian/fox/modules/fs/direntry"
	"github.com/Doridian/fox/modules/fs/fileinfo"
	lua "github.com/yuin/gopher-lua"
)

func doStat(L *lua.LState) int {
	path := L.CheckString(1)
	if path == "" {
		L.ArgError(1, "non-empty path expected")
		return 0
	}

	fi, err := os.Stat(path)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	return fileinfo.Push(L, fi)
}

func doLStat(L *lua.LState) int {
	path := L.CheckString(1)
	if path == "" {
		L.ArgError(1, "non-empty path expected")
		return 0
	}

	fi, err := os.Lstat(path)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	return fileinfo.Push(L, fi)
}

func doReadDir(L *lua.LState) int {
	path := L.CheckString(1)
	if path == "" {
		L.ArgError(1, "non-empty path expected")
		return 0
	}

	dirents, err := os.ReadDir(path)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	return direntry.PushArray(L, dirents)
}
