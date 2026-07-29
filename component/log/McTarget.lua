local ngx                       = ngx
local random                    = math.random
local timer_every               = ngx.timer.every
local timer_at                  = ngx.timer.at
local type                      = type
local ipairs                    = ipairs
local setmetatable              = setmetatable

local arrayHelper               = require('zues.utils.ArrayHelper')
local cjson                     = require('cjson')
local mc                        = require('resty.memcached')
local sendbuffer                = require('zues.component.log.sendbuffer')

local json_encode = cjson.encode
local in_array    = arrayHelper.in_array
local new_table   = arrayHelper.new_table

local buffers = new_table(0, 4)

-- 获取tcp服务器信息
local function _get_server(self)
    local index = 1
    local hostsNum = #self.hosts
    if hostsNum > 1 then
        index = random(1,#self.hosts)
    end
    local server = self.hosts[index]
    server.timeout = server.timeout or 500

    return server
end

local function send(broker_conf, topic, msg)
    local _mc = mc:new()
    _mc:set_timeout(broker_conf.timeout or 3000)
    local ok, err = _mc:connect(broker_conf.host, broker_conf.port)
    if not ok then
        ngx.log(ngx.ERR, 'failed to connect MC host:', broker_conf.host, ' port:', tostring(broker_conf.port), err)
        return nil, err
    end
    local ok, err = _mc:set(topic, msg)
    if ok then
        _mc:set_keepalive(10000, 100)
        return true
    else
        _mc:close()
    end

    return nil, err
end

-- 日志过滤
local function _message_filter(self, messages)
    local _messages, index = new_table(#messages, 0), 1
    for _, message in ipairs(messages) do
        if in_array(message['tag'], self.tags)
                and in_array(message['level'], self.levels) then
            if 'string' ~= type(message['message']) and 'number' ~= type(message['message']) then
                message['message'] = json_encode(message['message'])
            end
            _messages[index] = message
            index = index + 1
        end
    end

    return _messages
end

local _flush_buffer

local function _export(self, messages)
    if self.producer_type == 'async' then
        local buffer = buffers[self.cluster_name]
        local over = false
        for _, msg in ipairs(messages) do
            local message = msg['message']
            local tag = msg['tag']
            local topic = self.tag2topic[tag] or tag
            over = buffer:add(topic, 1, message)
        end
        if over then
            _flush_buffer(self)
        end
    else
        timer_at(0, function(_, target, messages)
            local broker_conf = _get_server(target)
            for _, msg in ipairs(messages) do
                local message = msg['message']
                local tag = msg['tag']
                local topic = target.tag2topic[tag] or tag
                send(broker_conf, topic, message)
            end
        end, self, messages)
    end
end

local function _flush_lock(self)
    if not self.flushing then
        if debug then
            ngx.log(ngx.DEBUG, "flush lock accquired")
        end
        self.flushing = true
        return true
    end
    return false
end


local function _flush_unlock(self)
    if debug then
        ngx.log(ngx.DEBUG, "flush lock released")
    end
    self.flushing = false
end

local function _flush(premature, self)
    if not _flush_lock(self) then
        if debug then
            ngx.log(ngx.DEBUG, "previous flush not finished")
        end
        return
    end

    local sendbuffer = buffers[self.cluster_name] or {}
    local buffer = {}
    for topic, partition_id, queue in sendbuffer:loop() do
        local msg = buffer[topic] or new_table(16,0)
        for _, item in ipairs(queue.queue) do
            msg[#msg+1] = item
        end
        buffer[topic] = msg
        sendbuffer:clear(topic, partition_id)
    end
    -- send_data
    if next(buffer) then
        for topic, msgList in pairs(buffer) do
            local msg = table.concat(msgList, self.delimiter)
            local broker_conf = _get_server(self)
            send(broker_conf, topic, msg)
        end
    end

    _flush_unlock(self)

    -- reset _timer_flushing_buffer after flushing complete
    self._timer_flushing_buffer = false

    return true
end

_flush_buffer = function (self)
    if self._timer_flushing_buffer then
        if debug then
            ngx.log(ngx.DEBUG, "another timer is flushing buffer, skipping it")
        end

        return
    end

    local ok, err = timer_at(0, _flush, self)
    if ok then
        self._timer_flushing_buffer = true
        return
    end

    ngx.log(ngx.ERR, "failed to create timer_at timer, err:", err)
end

local function _timer_flush(premature, self)
    self._timer_flushing_buffer = false
    _flush_buffer(self)
end

local cluster_inited = new_table(0, 4)

local _M = {_VERSION = '0.1'}
local mt = { __index = _M }

function _M:new()
    local o = {
        tags = {},
        levels = {'Error', 'Info', 'Warning', 'Trace'},
        enable = true,
        hosts = {},
        tag2topic = {},
        batch_num = 200,
        batch_size = 50000,
        producer_type = false,
        cluster_name = 'default',
        max_retry = 3,
        delimiter = ','
    }

    return setmetatable(o, mt)
end

function _M:init()
    if self.producer_type == 'async' then
        -- 判断是否初始化过
        if not cluster_inited[self.cluster_name] then
            -- 初始化buffer池子
            buffers[self.cluster_name] = sendbuffer:new(self.batch_num, self.batch_size, 0)
            local ok, err = timer_every((self.flush_time or 1000) / 1000, _timer_flush, self) -- default: 1s
            if not ok then
                ngx.log(ngx.ERR, "failed to create timer_every, err: ", err)
            end
            cluster_inited[self.cluster_name] = true
        end
    end
end

-- 日志收集
function _M:collect(messages)

    local messages = _message_filter(self, messages)
    if #messages < 1 then return end

    _export(self, messages)
end

return _M