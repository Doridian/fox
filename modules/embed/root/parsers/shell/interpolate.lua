local fs = require("go:fs")
local vars = require("embed:parsers.shell.vars")

local M = {}

-- TODO?: Implement $(cmd) interpolation (LAZY EVAL CRITICAL)
--        This basically requires () subshell parsing
--        such that we can make a stoppable parser that ends when
--        $() actually ends and not on a random ) char inside some string or whatever

-- TODO: Implement ${..#replace}
-- TODO: Implement ${..-default}
-- TODO?: Implement $0 (just $0, all others work)
-- TODO: Likely have to go per-char for above features

M.InterpTypeString = 1
M.InterpTypeFunc = 2

function M.generate(str, toks, escapeAllGlobs)
    local i = 1
    local varStart, varEnd, varName

    toks = toks or {}

    while i <= #str do
        varStart = str:find("$", i, true)
        if not varStart then
            break
        end

        if i < varStart then
            table.insert(toks, {
                type = M.InterpTypeString,
                escapeGlobs = escapeAllGlobs,
                value = str:sub(i, varStart - 1)
            })
        end

        varName = str:sub(varStart + 1, varStart + 1)
        if varName == "{"  then
            varStart = varStart + 1
            varEnd = str:find("}", varStart + 1, true)
            if not varEnd then
                return nil, "Unclosed variable ${}"
            end
        else
            varEnd = str:find("[^%%%w_]", varStart + 1)
        end
        if not varEnd then
            varEnd = #str + 1
        end
        varName = str:sub(varStart + 1, varEnd - 1)
    
        table.insert(toks, {
            type = M.InterpTypeFunc,
            escapeGlobs = true,
            value = function()
                return vars.get(varName) or ""
            end
        })

        i = varEnd + 1
    end

    if i > #str then
        return toks
    end

    table.insert(toks, {
        type = M.InterpTypeString,
        escapeGlobs = escapeAllGlobs,
        value = str:sub(i)
    })

    return toks
end

function M.singleToken(str, escapeGlobs)
    return {
        type = M.InterpTypeString,
        value = str,
        escapeGlobs = escapeGlobs
    }
end

function M.eval(toks, noFuncs, noGlobs)
    local ret = {}
    local retEsc = {}

    local hasGlobs = false
    for _, tok in ipairs(toks) do
        local v, vEsc
        if tok.type == M.InterpTypeString then
            v = tok.value
        elseif tok.type == M.InterpTypeFunc then
            v = tok:value()
            if noFuncs then
                return nil, nil
            end
        end
        vEsc = v
        if (not noGlobs) and fs.hasGlob(v) then
            if tok.escapeGlobs then
                vEsc = fs.globEscape(v)
            else
                hasGlobs = true
            end
        end
        table.insert(ret, v)
        table.insert(retEsc, vEsc)
    end

    return table.concat(ret, ""), table.concat(retEsc, ""), hasGlobs
end

return M
