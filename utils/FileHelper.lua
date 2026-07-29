local fopen                         = io.open
local re_match                      = ngx.re.match
local substr                        = string.sub
local str_byte                      = string.byte
local tonumber                      = tonumber

local ok, new_tab = pcall(require, "table.new")
if not ok or type(new_tab) ~= "function" then
    new_tab = function (narr, nrec) return {} end
end

local _M = {
    _VERSION = '0.1',
    FILE_APPEND = 'a', -- 文件追加
    FILE_WRITE = 'w', -- 覆盖写
}


-- 文件内容读取
function _M.file_get_content(filename)
    local fp, err = fopen(filename)
    if err then return '' end

    local content = fp:read('*a')
    fp:close()

    return content
end

-- 文件写入
function _M.file_put_contents(filename, content, mod)
    mod = mod or  'a+'
    local fp, err = fopen(filename, mod)
    if err then return false end
    fp:setvbuf('no')
    
    local n = fp:write(content)
    fp:close()

    return n
end

function _M.parse_ini_file(filename)
    local section_pattern = [[ \A \[ ([^ \[ \] ]+) \] \z ]]
    local keyvalue_pattern = [[ \A \s* ( [\w_\.]+ ) \s* = \s* ( ' [^']* ' | " [^"]* " | \S+ ) (?:\s*)? \z ]]
    local data = new_tab(0, 0)
    local section = "default"
    local fp, err = fopen(filename)
    if not fp then
        return data, "failed to open file: " .. (err or "")
    end

    for line in fp:lines() do
        local m = re_match(line, section_pattern, "jox")
        if m then
            section = m[1]
        else
            local m = re_match(line, keyvalue_pattern, "jox")
            if m then
                if not data[section] then
                    data[section] = {}
                end
                local key, value = m[1], m[2]
                local val = tonumber(value)
                if val then
                    value = val
                elseif value == "true" then
                    value = true
                elseif value == "false" then
                    value = false
                else
                    local fst = str_byte(value, 1)
                    if fst == 34 or fst == 39 then  -- ' or "
                        value = substr(value, 2, -2)
                    end
                end
                data[section][key] = value
            end
        end
    end
    fp:close()

    return data
end

return _M