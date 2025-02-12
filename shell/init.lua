function exit(code)
    _LAST_DO_EXIT = true
    _LAST_EXIT_CODE = code
end

local c = cmd.new()
c:cmd({"/bin/cat", "-"})

local p = c:stdinPipe()
print("S", c:start())

p:write("meow\nim a fomx :3\n")
p:close()

print("GO")
print("W", c:wait())
