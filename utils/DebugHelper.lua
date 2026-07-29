local ngx                   = ngx
local ngx_log               = ngx.log
local ngx_say               = ngx.say
local ngx_err               = ngx.ERR
local string                = string
local string_find           = string.find
local table_insert          = table.insert
local tostring              = tostring
local type                  = type
local pairs                 = pairs

local cjson                 = require "cjson"

local _M = {
    _VERSION = "0.0.1"
}

-- log记录table
function _M.log(node)

    if not node then
        ngx_log(ngx_err, "nil")
        return
    end

    if type(node) ~= 'table' then
        ngx_log(ngx_err, tostring(node))
        return
    end

    -- to make output beautiful
    local function tab(amt)
        local str = ""
        for i = 1, amt do
            str = str .. "\t"
        end
        return str
    end

    local cache, stack, output = {}, {}, {}
    local depth = 1
    local output_str = "{\n"

    while true do
        local size = 0
        for k, v in pairs(node) do
            size = size + 1
        end

        local cur_index = 1
        for k, v in pairs(node) do
            if (cache[node] == nil) or (cur_index >= cache[node]) then

                if (string_find(output_str, "}", output_str:len())) then
                    output_str = output_str .. ",\n"
                elseif not (string_find(output_str, "\n", output_str:len())) then
                    output_str = output_str .. "\n"
                end

                -- This is necessary for working with HUGE tables otherwise we run out of memory using concat on huge strings
                table_insert(output, output_str)
                output_str = ""

                local key
                if (type(k) == "number" or type(k) == "boolean") then
                    key = "[" .. tostring(k) .. "]"
                else
                    key = "['" .. tostring(k) .. "']"
                end

                if (type(v) == "number" or type(v) == "boolean") then
                    output_str = output_str .. tab(depth) .. key .. " = " .. tostring(v)
                elseif (type(v) == "table") then
                    output_str = output_str .. tab(depth) .. key .. " = {\n"
                    table_insert(stack, node)
                    table_insert(stack, v)
                    cache[node] = cur_index + 1
                    break
                else
                    output_str = output_str .. tab(depth) .. key .. " = '" .. tostring(v) .. "'"
                end

                if (cur_index == size) then
                    output_str = output_str .. "\n" .. tab(depth - 1) .. "}"
                else
                    output_str = output_str .. ","
                end
            else
                -- close the table
                if (cur_index == size) then
                    output_str = output_str .. "\n" .. tab(depth - 1) .. "}"
                end
            end

            cur_index = cur_index + 1
        end

        if (size == 0) then
            output_str = output_str .. "\n" .. tab(depth - 1) .. "}"
        end

        if (#stack > 0) then
            node = stack[#stack]
            stack[#stack] = nil
            depth = cache[node] == nil and depth + 1 or depth - 1
        else
            break
        end
    end

    -- This is necessary for working with HUGE tables otherwise we run out of memory using concat on huge strings
    table_insert(output, output_str)
    output_str = table.concat(output)

    ngx_log(ngx_err, output_str)
end

function _M.dump(args)
    if type(args) == 'table' then
        ngx_say('type: ', type(args), "\t", 'data: ', cjson.encode(args))
        ngx_log(ngx_err, type(args), "\t", 'data: ', cjson.encode(args))
    else
        ngx_say('type: ', type(args), "\t", 'data: ', args)
        ngx_log(ngx_err, 'type: ', type(args), "\t", 'data: ', args)
    end

    return ngx.exit(200)
end

return _M
