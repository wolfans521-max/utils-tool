local require       = require
local ngx           = ngx
local ngx_log       = ngx.log
local ngx_err       = ngx.ERR
local ngx_sleep     = ngx.sleep
local ngx_timer     = ngx.timer
local setmetatable  = setmetatable
local ngxShared     = ngx.shared
local ngx_now       = ngx.now
local type          = type
local pairs         = pairs

local confManager   = require('zues.component.Config')

local _M = { __VERSION = 0.1 }

local mt = { __index = _M }

function _M:new()
    return setmetatable({
        delay = 6,
        ttl = 600,
        times = 100,
        expire_rate = 0.75,
        conf_file = 'dict',
        conf = nil,
    }, mt)
end

function _M:init()
    if 'string' == type(self.conf_file) then
        self.conf = confManager:load(self.conf_file)
    end
    if not self.conf then
        ngx_log(ngx_err, "lost dict component conf!")
        return
    end
    for k, cf in pairs(self.conf) do
        if 'string' == type(cf.shared_name) and ngxShared[cf.shared_name] then
            cf.cacheInstance = ngxShared[cf.shared_name]
            cf.last_mod_time = 0
            cf.ttl = cf.ttl or self.ttl
            cf.delay = cf.delay or self.delay
            cf.times = cf.times or self.times
            cf.expire_rate = cf.expire_rate or self.expire_rate
        else
            self.conf[k] = nil
            ngx_log(ngx_err, "lost dict shared mem config! dict name :", cf.file)
        end
    end
end

local function _load(self)
    local file = io.open(self.file, "r")
    if not file then
        ngx_log(ngx_err, "open dict " , self.file , " failed")
        return false
    end
    if "function" ~= type(self.parse_line_callbackfn) then
        ngx_log(ngx_err, "lost parse_line_callbackfn conf, dict: " , self.file)
        file:close()
        return false
    end

    local times = 0
    for line in file:lines() do
        local ok, key, value = pcall(self.parse_line_callbackfn, line)
        if ok then
            if 'string' == type(key) and 'string' == type(value) then
                self.cacheInstance:set(key, value, self.ttl)
            end
        else
            ngx_log(ngx_err, "parse line format error")
        end
        if (times >= self.times) then
            -- switch coroutine every 100 lines
            times = 0
            ngx_sleep(0.01)
        end
        times = times + 1
    end
    file:close()
    return true
end

local function _check(self)
    if ngx_now() >= self.last_mod_time + self.ttl * self.expire_rate then
        return true
    end
    local fstat = io.popen("stat -c %Y " .. self.file)
    if fstat then
        local last_modified = fstat:read()
        if tonumber(last_modified) > self.last_mod_time then
            return true
        end
    else
        ngx_log(ngx_err, "check dict " , self.file, " failed")
    end

    return false
end

function _M:check_dict_load(cf)
    if _check(cf) then
        ngx.update_time()
        ngx_log(ngx_err, "dict ", cf.file, " load start")
        if _load(cf) then
            ngx.update_time()
            cf.last_mod_time = ngx_now()
            ngx_log(ngx_err, "dict ", cf.file, " load success")
        end
    end
end

local _handler
_handler = function (premature, self, cf)
    if premature then
        return
    end
    self:check_dict_load(cf)
    self:create_timer(cf)
end

function _M:create_timer(cf, delay)
    local ok, err = ngx_timer.at(delay or cf.delay, _handler, self, cf)
    if not ok then
        ngx_log(ngx_err, 'create dict timer failed! ', err, cf.file)
    end
end

function _M:run_timer()
    for _, cf in pairs(self.conf) do
        self:create_timer(cf, 0)
    end
end

function _M:get(dict_index, dkey)
    if not self.conf[dict_index] then
        ngx_log(ngx_err, "get none dict index ", dict_index)
        return nil
    end
    local cacheInstance = self.conf[dict_index].cacheInstance
    if not cacheInstance or not dkey then
        ngx_log(ngx_err, "lost dict index ", dict_index, " cache instance or key")
        return nil
    end

    return cacheInstance:get(dkey)
end

return _M
