local format                    = string.format
local now                       = ngx.now
local ngx_log                   = ngx.log
local ngx_err                   = ngx.ERR
local pcall                     = pcall
local type                      = type
local setmetatable              = setmetatable
local tostring                  = tostring

local mc                        = require('resty.memcached')
local confManager               = require('zues.component.Config')

local _M = {__VERSION = 0.1}

local mt = {__index = _M }

local function _log(self, statusCode, err)
    statusCode, err = statusCode or 200, err or 'success'
    if err ~= 'success' then statusCode = 502 end
    local cost = now() - self.start_time
    -- network:{protocol}|{urlKey}|{cost}|{nodeLog}|{len}|{tryTimes}|{statusCode}|{web_degrade}|msg:{err}
    local msg = format(' network:%s|%s|%.3f|%s|%d|%d|%d|0|msg:%s',
        'mc',
        self.busKey,
        cost,
        self.conf.port,
        0,
        1,
        statusCode,
        err
    )
    local log = get_manager():get('log')
    log:append_log(msg, 'req', 'Info')
    log:log({
        type = 'mc',
        conf = self.conf,
        request = self.request,
        busKey = self.busKey,
        timeout = self.timeout,
        cost = format('%.3f', cost),
        status = statusCode,
        err = err,
    }, 'debug','Debug')
end

local function _keepalive(self, idle, size)
    idle, size = idle or 10000, size or 1000
    if not self['conn'] or self.is_closed then return nil end

    self.conn:set_keepalive(idle, size)
    _log(self, 200)

    return true
end

local function _exec(self, command, ...)
    if not self['conn'] or self.is_closed then return nil end

    local ok, res, err = pcall(self.conn[command], self.conn, ...)
    self.request[#self.request + 1] = {
        command = command,
        args = {...},
        res = res,
        ok = ok,
    }
    if not ok or err then
        err = err or 'failed'
        _log(self, 504, err)
        self.conn:close()
        self.is_closed = true
        ngx_log(ngx_err, 'failed to exec MC cmd host:', self.conf.host, ' port:', tostring(self.conf.port), err)
    end

    return res
end

local function _connect(host, port, timeout)
    local _mc = mc:new()
    _mc:set_timeout(timeout)
    local ok, err = _mc:connect(host, port)
    if not ok then
        ngx_log(ngx_err, 'failed to connect MC host:', host, ' port:', tostring(port), err)
        return nil, err
    end

    return _mc
end

function _M:new()
    return setmetatable({
        timeout = 300,
        confKey = 'memc',
        conf = 'resource',
    }, mt)
end

function _M:init()
    if 'string' == type(self.conf) then
        self.conf = confManager:load(self.conf)
    end

    self.conf = self.conf[self.confKey]
end

-- args = {
--  busKey = '', -- 对应host的key
--  timeout = 30,
-- }
function _M:get_connect(args)
    local busKey = args['busKey']
    if not busKey then return nil, 'arg busKey not found!' end

    local conf = self.conf[busKey]
    if 'table' ~= type(conf) then return nil, 'node not found!' end

    local timeout = args['timeout'] or conf['timeout'] or self.timeout
    local start_time = now()
    local mc, err = _connect(conf.host, conf.port, timeout)
    local mc_obj = {
        conf = conf,
        busKey = busKey,
        start_time = start_time,
        timeout = timeout,
        exec = _exec,
        keepalive = _keepalive,
        conn = mc,
        request = {},
    }
    if not mc then
        _log(mc_obj, 502, err)
        return nil, err
    end

    return mc_obj
end

return _M