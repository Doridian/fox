local shell = require("go:shell")

local M = {}
local parserCache = {}

M.defaultParser = nil

M.search = {
    "parsers",
    "embed:parsers",
}
local function getParser(name)
    local c = parserCache[name]
    if c then
        if c.err then
            return nil, c.err
        end
        return c
    end

    local errs = {}
    for _, prefix in pairs(M.search) do
        local modName = prefix .. "." .. name
        local ok, mod = pcall(require, modName)
        if ok then
            parserCache[name] = mod
            return mod
        elseif not mod then
            table.insert(errs, "require() did not return table for ".. modName)
        else
            table.insert(errs, mod)
        end
    end

    local err = table.concat(errs, "\n")
    parserCache[name] = {
        err = err,
    }
    return nil, err
end

local function mustGetParser(name)
    local p, err = getParser(name)
    if not p then
        error("Error loading parser " .. name .. ": " .. err)
    end
    return p
end

local function defaultShellParserModule()
    return {
        run = shell.defaultShellParser,
    }
end

function M.run(cmdAdd, lineNo, prev)
    if not prev then
        prev = {
            prev = nil,
            parser = defaultShellParserModule(),
        }
        local cmdPrefix = cmdAdd:sub(1, 1)
        if cmdPrefix == "!" then
            cmdAdd = cmdAdd:sub(2)
            prev.parser = mustGetParser(cmdAdd)
            return prev, true
        elseif cmdPrefix == "=" or cmdAdd:sub(1, 2) == "--" then
            -- noop
        elseif cmdPrefix == "\\" then
            cmdAdd = cmdAdd:sub(2)
        elseif M.defaultParser then
            prev.parser = mustGetParser(M.defaultParser)
        end
    end

    local state, needMore, promptOverride = prev.parser.run(cmdAdd, lineNo, prev.prev)
    if needMore then
        prev.prev = state
        return prev, true, promptOverride
    end
    return state, false, promptOverride
end

return M
