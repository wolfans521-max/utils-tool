local format                    = string.format
local random                    = math.random
local ngxTimer                  = ngx.timer
local date                      = os.date
local tonumber                  = tonumber
local type                      = type
local ipairs                    = ipairs
local ngx                       = ngx
local setmetatable              = setmetatable

local arrayHelper               = require('zues.utils.ArrayHelper')
local manager                   = require('zues.component.Manager')
local cjson                     = require('cjson')

local json_encode = cjson.encode
local in_array    = arrayHelper.in_array
local new_table   = arrayHelper.new_table

-- 获取tcp服务器信息
local function _get_server(self)
    local index = 1
    local hostsNum = #self.hosts
    if hostsNum > 1 then
        index = random(1,#self.hosts)
    end
    local server = self.hosts[index]

    return server.host, tonumber(server.port), server.timeout or 500
end

local function _filter_check(self, message)
    local filter = self.filter
    if 'function' == type(filter) then
        if filter(self, message) then return true end
    end

    local uid = tonumber(message['uid'])
    if not uid then return false end

    -- 白名单
    local whiteList = self.register:config_get(self.key)
    if 'table' ~= type(whiteList) or #whiteList <= 0 then return false end

    if in_array(uid, whiteList) then return true end

    return false
end

-- 格式化日志时间
local function _format_date(timestamp)
    return date("%Y-%m-%d %H:%M:%S", timestamp)
end

-- 格式化
-- [消息格式text/json][时间][来源ip][索引][doc][用户uid]消息内容
-- source 来源(es中的索引名)
-- message 数据内容
-- format 数据格式，有text/json两种格式，默认text
-- uid 用户id，用于筛选
-- type 分类类型，用于筛选
-- time 时间，比如2019-03-29 14:23:03.888
-- '[text]['.$sdate.']['.$ip.']['.$ext['source'].']['.$ext['type'].']['.$params['cuid'].']'.$ext['message']
local function _format_message(self, message)
    -- 不能为空字符串
    local uid = message['uid']
    local msg = message['message']
    local tag = message['tag'] or ""
    local time = _format_date(message.timestamp)
    local seqId = message['seqId'] or '-'
    local source = self.commonPrefix .. "-" .. tag
    if #msg > 60000 then msg = string.sub(msg, 1, 60000) .. '...' end

    return format('[text][%s][%s][%s][%s][%s]%s',time, seqId, source, tag, uid, msg)
end

-- 日志过滤
local function _message_filter(self, messages)
    local _messages, index = new_table(#messages, 0), 1
    if not self.canUse then return _messages end
    for _, message in ipairs(messages) do
        if in_array(message['tag'], self.tags)
                and in_array(message['level'], self.levels)
                and _filter_check(self, message) then
            if 'string' ~= type(message['message']) and 'number' ~= type(message['message']) then
                message['message'] = json_encode(message['message'])
            end
            _messages[index] = message
            index = index + 1
        end
    end

    return _messages
end

local function _export(self, messages)
    ngxTimer.at(0, function(_, target, messages)
        local host, port, timeout = _get_server(target)
        local sock = ngx.socket.udp()
        sock:settimeout(timeout)
        local ok, err = sock:setpeername(host, port)
        if not ok then
            ngx.log(ngx.ERR, 'failed connect to host:'.. host, err)
            return false
        end
        for _, msg in ipairs(messages) do
            local message = _format_message(target, msg)
            local _, _err = sock:send(message)
            if _err then
                ngx.log(ngx.ERR, 'failed send to host:' .. host, _err, message)
            end
        end
        sock:close()
    end, self, messages)
end

local _M = {_VERSION = '0.1'}
local mt = { __index = _M }

function _M:new()
    local o = {
        tags = {},
        levels = {'Error', 'Info', 'Warning', 'Trace'},
        enable = true,
        hosts = {},
        key = 'elk_white',
        register = 'discovery'
    }

    return setmetatable(o, mt)
end

-- 日志收集
function _M:collect(messages)

    local messages = _message_filter(self, messages)
    if #messages < 1 then return end

    _export(self, messages)
end

function _M:init()
    local cType = type(self.register)
    if cType == 'table' then
        self.register = manager:create_object(self.register)
    elseif cType == 'string' then
        self.register = get_manager(self):get(self.register)
    end

    if 'table' == type(self.register) then self.canUse = true end
end

return _M