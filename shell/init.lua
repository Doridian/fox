function exit(code)
    _LAST_DO_EXIT = true
    _LAST_EXIT_CODE = code
end

local c = cmd.new()
c:cmd({"/bin/cat", "-"})

local c2 = cmd.new()
c2:cmd({"/bin/echo", "meow", "fomx"})

local p = c:stdinPipe()
c:start()
c2:stdout(p, false)
c2:run()

p:write("\nhi\n")
p:close()

print("stdout", pipe.stdout)
pipe.stdout:write("stdout direct\n")

print("W", c:wait())
