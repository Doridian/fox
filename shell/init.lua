function exit(code)
    _LAST_DO_EXIT = true
    _LAST_EXIT_CODE = code
end

local c = cmd.new()
c:cmd({"/bin/cat", "-"})

local c2 = cmd.new()
c2:cmd({"/bin/echo", "meow", "test"})

c:stdin(c2:stdoutPipe())

print("GO")

print(c:run())
