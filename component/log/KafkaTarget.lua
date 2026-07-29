--文件类
local ngx                   = ngx
local pairs                 = pairs
local type                  = type
local ngxTimer              = ngx.timer
local setmetatable          = setmetatable

local arrayHelper           = require('zues.utils.ArrayHelper')
local cjson                 = require('cjson')
local producer              = require "resty.kafka.producer"

local in_array              = arrayHelper.in_array
local json_encode           = cjson.encode
local new_table             = arrayHelper.new_table


-- 日志过滤
local function _message_filter(self, messages)
    local _messages, index = new_table(#messages, 0), 1
    for _, message in pairs(messages) do
        if in_array(message['tag'], self.tags) and in_array(message['level'], self.levels) then
            if 'string' ~= type(message['message']) and 'number' ~= type(message['message'])then
                message['message'] = json_encode(message['message'])
            end
            _messages[index] = message
            index = index + 1
        end
    end

    return _messages
end

-- 格式化日志
local function _format_message(self, message)
    return message.message
end

local id = 0

local function partition_key()
    local ret = tostring(id)
    id = id + 1
    if id > 1073741824 then id = 0 end

    return ret
end

-- 内容导出
local function _export(self, messages)
    ngxTimer.at(0, function(_, target, messages)
        local pd = producer:new(self.broker_list, self.producer_config, self.cluster_name)
        for _, msg in ipairs(messages) do
            local message = _format_message(target, msg)
            local topic = self.tag2topic[msg['tag']]
            if topic then
                local _, _err = pd:send(topic, partition_key(), message)
                if _err then
                    ngx.log(ngx.ERR, 'failed send to topic:' .. topic, _err, msg['tag'], message)
                end
            else
                ngx.log(ngx.ERR, 'not found topic for tag:' .. msg['tag'])
            end
        end
    end, self, messages)
end

local function error_handle(topic, partition_id, message_queue, index, err, retryable)
    ngx.log(ngx.ERR, "send msg to kafka failed, topic:", topic, " num:", index, " err:", err)
end

local _M = {_VERSION = '0.1', openfiles = {}}

-- local config = {
--     broker_list = {
--         {
--             host = "127.0.0.1",
--             port = 9092,
--         }
--     },
--     sasl_config = {
--         mechanism = "PLAIN",
--         user = "USERNAME",
--         password = "PASSWORD"
--     },
--     cluster_name = 'miniapp',
-- }

function _M:new()
    local o = {
        tags = {},
        levels = {'Error', 'Info', 'Warning', 'Trace'},
        enable = true,
        delimiter = '_',
        tag2topic = {},
        broker_list = {},
        sasl_config = {
            mechanism = "PLAIN",
        },
        producer_config = {},
        cluster_name = '',
    }

    return setmetatable(o, { __index = _M })
end

function _M:init()
    -- 初始化broker_list配置
    local broker_list = self.broker_list
    if not broker_list then self.enable = false end
    local sasl_config = self.sasl_config
    for _, broker in ipairs(broker_list) do
        broker.sasl_config = sasl_config
    end

    self.producer_config.error_handle = error_handle
end

-- 日志收集
function _M:collect(messages)

    local messages = _message_filter(self, messages)
    if #messages < 1 then return end

    _export(self, messages)
end

return _M