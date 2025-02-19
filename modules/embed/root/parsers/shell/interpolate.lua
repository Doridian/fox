local fs = require("go:fs")
local vars = require("embed:parsers.shell.vars")

local M = {}

-- TODO: Lazy interpolation, only once the string is actually grabbed/neeeded
--       PreReq for $(cmd) interpolation
--       And also correct behaviour

-- TODO?: Implement $(cmd) interpolation (LAZY EVAL CRITICAL)
--        This basically requires () subshell parsing
--        such that we can make a stoppable parser that ends when
--        $() actually ends and not on a random ) char inside some string or whatever

-- TODO: Implement ${..#replace}
-- TODO: Implement ${..-default}

-- return true to indicate that glob processing mode should be enabled
function M.run(str, toks)
    local i = 1
    local varStart, varEnd, varTmp, varType

    toks = toks or {}

    while i <= #str do
        varStart = str:find("[$%%]", i)
        if not varStart then
            break
        end

        table.insert(toks, {
            type = "str",
            value = str:sub(i, varStart - 1)
        })

        varType = str:sub(varStart, varStart)

        varTmp = str:sub(varStart + 1, varStart + 1)
        if varTmp == "{"  then
            varStart = varStart + 1
            varEnd = str:find("}", varStart + 1, true)
            if not varEnd then
                return nil, "Unclosed variable ${}"
            end
        else
            varEnd = str:find("[^%w_]", varStart + 1)
        end
        if not varEnd then
            varEnd = #str + 1
        end
        varTmp = str:sub(varStart + 1, varEnd - 1)

        table.insert(toks, {
            type = "func",
            escapeGlobs = true,
            value = function()
                return vars.get(varType, varTmp)
            end
        })

        i = varEnd + 1
    end

    table.insert(toks, {
        type = "str",
        value = str:sub(i)
    })

    return toks
end

function M.eval(toks, escapeGlobs)
    local ret = ""
    for _, tok in ipairs(toks) do
        local v
        if tok.type == "str" then
            v = tok.value
        elseif tok.type == "func" then
            v = tok.value()
        end
        if tok.escapeGlobs and escapeGlobs then
            v = fs.globEscape(v)
        end
        ret = ret .. v
    end
    return ret
end

return M
