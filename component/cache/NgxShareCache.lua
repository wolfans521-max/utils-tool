local ngxShared                         = ngx.shared
local getmetatable                      = getmetatable
local setmetatable                      = setmetatable
local type                              = type

local parent                            = require('zues.component.cache.Cache')

local _M = {__VERSION = '0.1'}
local mt = {__index = _M }


function _M:new()
    local obj = parent:new()
    local super_mt = getmetatable(obj)
    -- 当方法在子类中查询不到时，再去父类中去查找。
    setmetatable(_M, super_mt)
    -- 这样设置后，可以通过self.super.method(self, ...) 调用父类的已被覆盖的方法。
    obj.super = setmetatable({}, super_mt)
    obj.cacheInstance = 'my_cache'

    return setmetatable(obj, mt)
end

function _M:init()
    if 'string' == type(self.cacheInstance) then
        local instance = self.cacheInstance
        self.cacheInstance = ngxShared[instance]
    end
end

-- 获取值
function _M:get_value(key)
    return self.cacheInstance:get(key)
end

-- 种值
function _M:set_value(key, value, duration)
    return self.cacheInstance:set(key, value, duration)
end

return _M