function exit(code)
    _LAST_DO_EXIT = true
    _LAST_EXIT_CODE = code
end

local c = shellcmd.new()
print(#c:args())
c:path("/bin/echo")
print(#c:args())
c:args({"echo", "hello", "world"})
print(#c:args())
c:errorPropagation(false)
local c2, code, err = c:run()

print("INIT OK!", c2, code, err)
