package fs

import (
	"fmt"
	"os"

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

	return file.Push(L, f)
}
