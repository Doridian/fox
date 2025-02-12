function exit(code)
    _LAST_DO_EXIT = true
    _LAST_EXIT_CODE = code
end

local c = shellcmd.new()
c:path("/bin/grep")
c:args({"-F", "test"})

local c2 = shellcmd.new()
c2:path("/bin/echo")
c2:args({"meow", "test"})

c:stdin(c2:stdout())

print("GO")

print(c:run())
