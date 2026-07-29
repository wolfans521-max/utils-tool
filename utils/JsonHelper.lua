local type        = type
local pcall       = pcall

local cjson       = require('cjson')

local json_encode = cjson.encode
local json_decode = cjson.decode

local _M = {_VERSION='0.1'}

-- 编码
function _M.encode(data) return json_encode(data) end

-- 解码
function _M.decode(str)
    if 'string' ~= type(str) then return str end

    local ok, res = pcall(json_decode, str)

    if not ok then res = nil end

    return res
end

return _M