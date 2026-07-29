-- cache 基类
local md5                       = ngx.md5
local type                      = type
local pcall                     = pcall
local setmetatable              = setmetatable

local jsonHelper                = require('zues.utils.JsonHelper')

local json_decode = jsonHelper.decode
local json_encode = jsonHelper.encode

-- 构建build
local function _build_key(self, key)
    if 'string' == type(key) then
        if #key > 32 then
            key = md5(key)
        end
    else
        key = md5(json_encode(key));
    end

    return self.keyPrefix .. key
end

local _M = {_VERSION = '0.1'}
local mt = {__index = _M}

function _M:init() end

-- 获取缓存值
function _M:get(key, default, needSerialize)
    if nil == needSerialize then needSerialize = true end
    key = _build_key(self,key)
    local value = self:get_value(key)

    if not value then
        return default
    end

    if needSerialize then
        if not self.serializer[2] then
            value = json_decode(value)
        else
            local _
            _, value = pcall(self.serializer[2], value)
        end
    end

    return value
end

-- 设置缓存值
function _M:set(key, value, duration, needSerialize)
    duration = duration or self.defaultDuration
    if nil == needSerialize then needSerialize = true end

    if needSerialize then
        if  not self.serializer[1] then
            value = json_encode(value)
        else
            local _
            _, value = pcall(self.serializer[1], value)
        end
    end

    key = _build_key(self,key)

    return self:set_value(key, value, duration)
end

-- 子类实现
function _M:get_value(key)
    return nil
end

-- 子类实现
function _M:set_value(key, value, duration)
    return false
end

-- 创建
function _M:new()
    return setmetatable({
        keyPrefix = '',
        defaultDuration = 0,
        serializer = {},
    }, mt)
end

return _M