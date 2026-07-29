local math              = math
local value             = require "zues.utils.ValueHelper"
local debug             = require 'zues.utils.DebugHelper'
local string            = require "zues.utils.StringHelper"

-- 62进制转换工具模块（OpenResty 风格实现）
local _M = {
    _VERSION = "1.0",
    _string = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ',
    _encode_block_size = 7,
    _decode_block_size = 4
}

-- 批量10进制转62进制
function _M.multi_from10to62(mids)
    local res = {}
    for _, mid in ipairs(mids) do
        mid = string.format('%.0f', mid)
        res[mid] = _M.from10to62(mid)
    end
    return res
end

-- 批量62进制转10进制
function _M.multi_from62to10(mids, compat, for_mid)
    local res = {}
    compat = compat or false
    for_mid = for_mid == nil and true or for_mid
    for _, mid in ipairs(mids) do
        mid = string.format('%.0f', mid)
        res[mid] = _M.from62to10(mid, compat, for_mid)
    end
    return res
end

-- 10进制转62进制
function _M.from10to62(mid_str)
    assert(type(mid_str) == "string", "mid_str must be string")
    local result = ""
    local mid_len = #mid_str
    local segments = math.ceil(mid_len / _M._encode_block_size)
    local start = mid_len

    for i = 1, segments - 1 do
        start = start - _M._encode_block_size
        local seg = mid_str:sub(start + 1, start + _M._encode_block_size)
        seg = _M._encode_segment(tonumber(seg))
        result = string.rep('0', _M._decode_block_size - #seg) .. seg .. result
    end
    result = _M._encode_segment(tonumber(mid_str:sub(1, start))) .. result
    return result
end

-- 62进制转10进制
function _M.from62to10(to_mid, compat, for_mid)
    assert(type(to_mid) == "string", "to_mid must be string")
    local mid = ""
    local s_len = #to_mid
    local segments = math.ceil(s_len / _M._decode_block_size)
    local start = s_len

    for i = 1, segments - 1 do
        start = start - _M._decode_block_size
        local seg = to_mid:sub(start + 1, start + _M._decode_block_size)
        seg = _M._decode_segment(seg)
        mid = string.rep('0', _M._encode_block_size - #seg) .. seg .. mid
    end
    mid = _M._decode_segment(to_mid:sub(1, start)) .. mid

    for_mid = for_mid == nil and true or for_mid
    if for_mid then
        local mid_len = #mid
        local first = mid:sub(1, 1)
        if (mid_len == 16 and (first == '3' or first == '4')) or (mid_len == 19 and first == '5') then
            return mid
        end
    end

    compat = compat or false
    if compat and not string.match(mid:sub(1, 3), '^(109|110|201|211|221|231|241)') then
        mid = _M._decode_segment(s:sub(1, 4)) .. _M._decode_segment(s:sub(5))
    end

    if for_mid and mid:sub(1, 1) == '1' and #mid > 8 and mid:sub(8, 8) == '0' then
        mid = mid:sub(1, 7) .. mid:sub(9)
    end

    return mid
end

-- 单个段10进制转62进制
function _M._encode_segment(num)
    if num == 0 then
        return '0'
    end
    local out = ""
    while num > 0 do
        local idx = (num % 62) + 1
        out = _M._string:sub(idx, idx) .. out
        num = math.floor(num / 62)
    end
    return out
end

-- 单个段62进制转10进制
function _M._decode_segment(s)
    local out = 0
    local base = 1
    for i = #s, 1, -1 do
        local char = s:sub(i, i)
        local idx = _M._string:find(char, 1, true) - 1
        out = out + idx * base
        base = base * 62
    end
    out = string.format('%.0f', out)
    return out
end

return _M
