local require       = require
local ngx           = ngx
local setmetatable  = setmetatable
local ngxShared     = ngx.shared

local bit           = require("bit")
local cjson         = require('zues.utils.JsonHelper')
local value         = require('zues.utils.ValueHelper')
local hash          = require('resty.hash')

local bkdr_hash     = hash.bkdr_hash

local _M = { __VERSION = 0.1 }

local mt = { __index = _M }

local int_bit_num = 31

local default_seed = 31

function _M:new()
    return setmetatable({
        dict = 'bloomfilter',
        item_count = 1000,
        false_rate = 0.001,
        expire = 0,
        seed_list = {31,739,1087,2161,3251,73,751,101,1021,127,1307,223,1553,283,37,2053,41,3469,43,4691,47,2777,53,3709,59,61,3943,67,71},
    }, mt)
end

local Ln2 = 0.693147180559945309417232121458176568075500134360255254120680009

local function get_filter_param(false_rate, item_count)
    local m = -1 * math.log(false_rate) * item_count / (Ln2 * Ln2)
    m = math.ceil(m)

    local k = m  * Ln2 / item_count
    if k > #self.seed_list then k = self.seed_list end

    return m, math.ceil(k)
end

function _M:init()
    self.m, self.hash_num = get_filter_param(self.false_rate, self.item_count)

    self.int_num = math.ceil(self.m / int_bit_num)

    self.dict = ngxShared[self.dict]
end

function _M:set(key, data, duration)
    duration = duration or self.expire
    local bitmap = self.dict:get(key)
    if bitmap then bitmap = cjson.decode(bitmap) end
    if 'table' ~= type(bitmap) then bitmap = {} end
    local itemNum = self.int_num

    for i = 1, self.hash_num do
        local hashNum = bkdr_hash(data, self.seed_list[i] or default_seed)
        local group_index = tostring(hashNum % itemNum)
        local item = bitmap[group_index] or 0
        local index = hashNum % int_bit_num
        bitmap[group_index] = bit.bor(item, bit.lshift(1, index))
    end

    -- 保存结果
    self.dict:set(key, cjson.encode(bitmap), duration)
end

function _M:check(key, data)
    local bitmap = self.dict:get(key)
    if bitmap then bitmap = cjson.decode(bitmap) end
    if value.is_empty_table(bitmap) then return false end
    local itemNum = self.int_num

    for i = 1, self.hash_num do
        local hashNum = bkdr_hash(data, self.seed_list[i] or default_seed)
        local group_index = tostring(hashNum % itemNum)
        local item = bitmap[group_index] or 0
        local index = hashNum % int_bit_num
        local res = bit.band(item, bit.lshift(1, index))
        if res == 0 then return false end
    end

    return true
end

function _M:clear(key)
    self.dict:delete(key)
end

return _M