local shell = require("go:shell")

local M = {}
local parserCache = {}

M.defaultParser = nil

M.search = {
    "embed:parsers",
    "go:parsers",
    "parsers",
}
local function getParser(name)
    if parserCache[name] then
        return parserCache[name]
    end

    local errs = {}
    for _, prefix in pairs(M.search) do
        local ok, mod = pcall(require, prefix .. "." .. name)
        if ok then
            parserCache[name] = mod
            return mod
        else
            table.insert(errs, mod)
        end
    end

    if #errs > 0 then
        error(table.concat(errs, "\n"))
    end

    return nil
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
            prev.parser = getParser(cmdAdd)

            if not prev.parser then
                print("Unknown parser " .. cmdAdd)
                return ""
            end

            return prev, true
        elseif cmdPrefix == "=" or cmdAdd:sub(1, 2) == "--" then
            -- noop
        elseif cmdPrefix == "/" then
            cmdAdd = cmdAdd:sub(2)
        elseif M.defaultParser then
            local parser = getParser(M.defaultParser)
            if parser then
                prev.parser = parser
            end
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
