function exit(code)
    _LAST_DO_EXIT = true
    _LAST_EXIT_CODE = code
end

local c = shellcmd.new()
c:path("/bin/echo")
c:args({"hello", "world"})
c:errorPropagation(false)
local c2, code, err = c:run()

print("INIT OK!", c2, code, err)
