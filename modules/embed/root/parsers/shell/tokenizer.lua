local fs = require("go:fs")
local interpolate = require("embed:parsers.shell.interpolate")

local M = {}

M.ArgTypeString = 1
M.ArgTypeStringFunc = 2
M.ArgTypeOp = 3

-- TODO: Variable sub-parser should decide end of variables
--       This makes ${} variables fully work

local function resolveStringFunc(token, args)
    local str, strEscaped, hasGlobs = interpolate.eval(token.sub)

    args = args or {}
    if hasGlobs then
        local matches = fs.glob(strEscaped)
        if #matches > 0 then
            for _, match in pairs(matches) do
                table.insert(args, match)
            end
            return
        end
    end
    table.insert(args, str)
    return args
end

function M.oneStringVal(token)
    local v = M.stringVals(token)
    return v and v[1]
end

function M.stringVals(token, args)
    if token.type == M.ArgTypeString then
        if args then
            table.insert(args, token.value)
            return args
        end
        return { token.value }
    elseif token.type == M.ArgTypeStringFunc then
        return token:value(args)
    end
    error("invalid token type for stringVals")
end

function M.run(parsed)
    local i = 1
    local tokens = {}
    local curToken, nextControlIdx, nextControl, controlEndIdx
    local function bufToken(container)
        if not curToken then
            curToken = {
                type = M.ArgTypeStringFunc,
                sub = {},
                value = resolveStringFunc,
            }
        end

        local sub, err
        if nextControlIdx then
            sub = parsed:sub(i, nextControlIdx - 1)
            i = nextControlIdx + 1
        else
            sub = parsed:sub(i)
            i = #parsed + 1
        end

        if sub ~= "" then
            local escapeGlobs = false
            if container then
                escapeGlobs = true
            end

            if container ~= "'" then
                _, err = interpolate.run(sub, curToken.sub, escapeGlobs)
                if err then
                    return nil, "shell.interpolate error: " .. tostring(err)
                end
            else
                table.insert(curToken.sub, {
                    type = "str",
                    value = sub,
                    escapeGlobs = escapeGlobs,
                })
            end
        end
    end
    local function manualPushToken()
        if #curToken.sub > 0 then
            local arg = curToken
            curToken = nil
            table.insert(tokens, arg)
        end
    end
    local function pushToken(container)
        bufToken(container)
        manualPushToken()
    end
    while i <= #parsed do
        nextControlIdx = parsed:find("[ \n\t\"';&|><!]", i)
        if not nextControlIdx then
            pushToken()
            break
        end

        nextControl = parsed:sub(nextControlIdx, nextControlIdx)
        if nextControl == "'" or nextControl == '"' then
            controlEndIdx = parsed:find(nextControl, nextControlIdx + 1)
            if not controlEndIdx then
                return nil, "unclosed " .. nextControl
            end

            bufToken()
            nextControlIdx = controlEndIdx
            bufToken(nextControl)
        elseif nextControl == "\n" or nextControl == "\r" or nextControl == "\t" or nextControl == " " then
            pushToken()
        else
            bufToken()
            local curTokenSub = curToken and curToken.sub
            local curTokenSubOne = curTokenSub and curTokenSub[1]
            if (nextControl ~= "<" and nextControl ~= ">") or not (curTokenSubOne and #curTokenSub == 1 and tonumber(curTokenSubOne.value)) then
                manualPushToken()
            end

            controlEndIdx = nextControlIdx
            local hasAmpersand = false
            while true do
                local nextC = parsed:sub(controlEndIdx + 1, controlEndIdx + 1)
                if nextC ~= nextControl then
                    if nextC == "&" then
                        if hasAmpersand then
                            return nil, "cannot have && right after >"
                        end
                        hasAmpersand = true
                    else
                        break
                    end
                end
                controlEndIdx = controlEndIdx + 1
            end

            table.insert(tokens, {
                pre = curTokenSubOne and curTokenSubOne.value,
                value = nextControl,
                len = controlEndIdx - nextControlIdx + 1,
                raw = parsed:sub(nextControlIdx, controlEndIdx),
                hasAmpersand = hasAmpersand,
                type = M.ArgTypeOp,
            })
            curToken = nil
            i = controlEndIdx + 1
        end
    end

    return tokens
end

return M
