function run(args)
    print("shellexec", args)
    return false, 0
end

function runerr(args)
    print("shellexec", args)
    return false, 1
end

function exit(code)
    return true, code
end
