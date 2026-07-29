--motan client
local now                               = ngx.now
local getmetatable                      = getmetatable
local setmetatable                      = setmetatable

local parent                            = require('zues.component.client.ClientInterface')
local cjson                             = require('zues.utils.JsonHelper')
local motan                             = require('resty.mesh')

local _M = {_VERSION = '0.1'}
local mt = {__index  = _M}

function _M:new()
    local obj = parent:new()
    local super_mt = getmetatable(obj)
    -- 当方法在子类中查询不到时，再去父类中去查找。
    setmetatable(_M, super_mt)
    -- 这样设置后，可以通过self.super.method(self, ...) 调用父类的已被覆盖的方法。
    obj.super = setmetatable({}, super_mt)
    return setmetatable(obj, mt)
end

-- 发送请求
function _M:s_curl(option)
    local startTime = now()
    local mesh = motan:new()
    local timeout = option.timeout or self.timeout
    mesh:set_timeouts(timeout[1],timeout[2],timeout[3])
    local url = option.url
    local header = option['header'] or {}
    local args = option['args'] or {}
    if option['contentJson'] then args = cjson.encode(args) end
    local res, _, _, err = mesh:request(url, args, header)
    if not err and res then
        res = {status = 200, body = res}
    else
        res = {status = 500, body = '[]'}
    end

    --日志记录
    self:log(option, res.body, now() - startTime, res.status, err)

    return res.body
end

return _M