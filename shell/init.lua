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

c2:stdout(c)

c2:start()
print("Ran c2")
c:start()
print("Ran c")

print(c:wait())
