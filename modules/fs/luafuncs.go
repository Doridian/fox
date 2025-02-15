package fs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Doridian/fox/modules/fs/direntry"
	"github.com/Doridian/fox/modules/fs/file"
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
		if errors.Is(err, os.ErrNotExist) {
			L.Push(lua.LNil)
			return 1
		}
		L.RaiseError("%v", err)
		return 0
	}
	return fileinfo.PushNew(L, fi)
}

func doLStat(L *lua.LState) int {
	path := L.CheckString(1)
	if path == "" {
		L.ArgError(1, "non-empty path expected")
		return 0
	}

	fi, err := os.Lstat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			L.Push(lua.LNil)
			return 1
		}
		L.RaiseError("%v", err)
		return 0
	}
	return fileinfo.PushNew(L, fi)
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

// open(path, modes)
// r = read
// w = write
// a = append
// b = binary (ignored)
// o = overwrite (do not create new files)
// e = exclusive (never open an existing file)
// s = sync
func doOpen(L *lua.LState) int {
	path := L.CheckString(1)
	if path == "" {
		L.ArgError(1, "non-empty path expected")
		return 0
	}

	modes := L.CheckString(2)
	if modes == "" {
		L.ArgError(2, "non-empty mode expected")
		return 0
	}

	chmod := os.FileMode(L.OptInt(3, 0666))

	flag := 0
	allowRead := false
	allowWrite := false
	doAppend := false
	doCreate := true
	for _, mode := range modes {
		switch mode {
		case 'r':
			allowRead = true
		case 'w':
			allowWrite = true
		case 'a':
			allowWrite = true
			doAppend = true
		case 'b':
			// ignore b flag
		case 'o':
			doCreate = false
		case 'e':
			flag |= os.O_EXCL
		case 's':
			flag |= os.O_SYNC
		default:
			L.ArgError(2, fmt.Sprintf("invalid mode %c", mode))
			return 0
		}
	}

	if allowRead && allowWrite {
		flag = os.O_RDWR
	} else if allowRead {
		flag = os.O_RDONLY
	} else if allowWrite {
		flag = os.O_WRONLY
	}

	if allowWrite {
		if doAppend {
			flag |= os.O_APPEND
		} else {
			flag |= os.O_TRUNC
		}
		if doCreate {
			flag |= os.O_CREATE
		}
	}

	f, err := os.OpenFile(path, flag, chmod)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	return file.PushNew(L, f)
}

func doRemove(L *lua.LState) int {
	path := L.CheckString(1)
	if path == "" {
		L.ArgError(1, "non-empty path expected")
		return 0
	}

	err := os.Remove(path)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	return 0
}

func doMkdir(L *lua.LState) int {
	path := L.CheckString(1)
	if path == "" {
		L.ArgError(1, "non-empty path expected")
		return 0
	}

	err := os.Mkdir(path, os.FileMode(L.OptInt(2, 0777)))
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	return 0
}

func doMkdirAll(L *lua.LState) int {
	path := L.CheckString(1)
	if path == "" {
		L.ArgError(1, "non-empty path expected")
		return 0
	}

	err := os.MkdirAll(path, os.FileMode(L.OptInt(2, 0777)))
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	return 0
}

func doGlob(L *lua.LState) int {
	pattern := L.CheckString(1)
	if pattern == "" {
		L.ArgError(1, "non-empty pattern expected")
		return 0
	}

	matches, err := filepath.Glob(pattern)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	ret := L.CreateTable(len(matches), 0)
	for _, match := range matches {
		ret.Append(lua.LString(match))
	}
	L.Push(ret)
	return 1
}

func doGlobEscape(L *lua.LState) int {
	str := L.CheckString(1)
	if str == "" {
		L.ArgError(1, "non-empty str expected")
		return 0
	}

	// TODO: Better alternative
	L.Push(lua.LString(regexp.QuoteMeta(str)))
	return 1
}
