local format                    = string.format
local now                       = ngx.now
local ngx_log                   = ngx.log
local ngx_err                   = ngx.ERR
local type                      = type
local setmetatable              = setmetatable

local confManager               = require('zues.component.Config')
local array                     = require('zues.utils.ArrayHelper')
local producer                  = require "resty.kafka.producer"

local simple_merge              = array.simple_merge
local _M = {__VERSION = 0.1}

local mt = {__index = _M }

local function _log(self, statusCode, err)
    statusCode, err = statusCode or 200, err or 'success'
    if err ~= 'success' then statusCode = 502 end
    local cost = now() - self.start_time
    -- network:{protocol}|{urlKey}|{cost}|{nodeLog}|{len}|{tryTimes}|{statusCode}|{web_degrade}|msg:{err}
    local msg = format(' network:%s|%s|%.3f|%s|%d|%d|%d|0|msg:%s',
        'kafka',
        self.busKey,
        cost,
        0,
        0,
        1,
        statusCode,
        err
    )
    local log = get_manager():get('log')
    log:append_log(msg, 'req', 'Info')
    log:log({
        type = 'kafka',
        conf = self.conf,
        request = self.request,
        busKey = self.busKey,
        cost = format('%.3f', cost),
        producer_config = self.producer_config,
        status = statusCode,
        err = err,
    }, 'debug','Debug')
end

local function _send(self, topic, key, message)
    local pd = self['pd']
    if not pd then return nil end
    local offset, err = pd:send(topic, key, message)

    self.request[#self.request + 1] = {
        args = {topic = topic, key = key, message = message},
        ok = offset,
        err = err
    }

    if not offset or err then
        err = err or 'failed'
        _log(self, 504, err)
        ngx_log(ngx_err, 'failed to exec kafka cmd busyKey:', self.conf.busKey, err)
    end

    return offset, err
end

function _M:new()
    return setmetatable({
        timeout = 300,
        db = 0,
        confKey = 'kafka',
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
--  producer_config = nil,
--  cluster_name = nil,
-- }
function _M:get_producer(args)
    local busKey = args['busKey']
    if not busKey then return nil, 'arg busKey not found!' end

    local conf = self.conf[busKey]
    if 'table' ~= type(conf) then return nil, 'nodes not found!' end

    local broker_list = conf.broker_list
    local producer_config = simple_merge(conf.producer_config, args.producer_config, true)
    local cluster_name = args.cluster_name or conf.cluster_name
    local start_time = now()
    local pd = producer:new(broker_list, producer_config, cluster_name)
    local pd_obj = {
        conf = conf,
        busKey = busKey,
        start_time = start_time,
        producer_config = producer_config,
        send = _send,
        pd = pd,
        request = {},
    }

    return pd_obj
end

return _M