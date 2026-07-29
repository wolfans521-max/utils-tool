--日志收集器
local time                      = ngx.time
local ngxTimer                  = ngx.timer.at
local type                      = type
local setmetatable              = setmetatable
local pairs                     = pairs
local ngx                       = ngx

local manager                   = require('zues.component.Manager')
local array                     = require('zues.utils.ArrayHelper')
local util_value                = require('zues.utils.ValueHelper')
local cjson                     = require "zues.utils.JsonHelper"

local new_table    = array.new_table

local function _flush_append_log(self)
    local messages = ngx.ctx[self.appendLogKey]
    if 'table' == type(messages) then
        for tag, val in pairs(messages) do
            for level, message in pairs(val) do
                self:log(message, tag, level)
            end
        end
    end

    ngx.ctx[self.appendLogKey] = nil
end

local function _seq_id(self)
    local seqHandler = self.seqHandler
    if type(seqHandler) == 'function' then return seqHandler() end

    local seqId = get_manager(self):get('request'):seqid()
    if seqId then return seqId end

    return '-'
end

local _M = {
    _VERSION = '0.1',
    LOG_LEVEL_INFO = 'Info',
    LOG_LEVEL_ERROR = 'Error',
    LOG_LEVEL_WARING = 'Waring',
    LOG_LEVEL_NOTICE = 'Notice',
    LOG_LEVEL_DEBUG  = 'Debug',
}

function _M:new()
    local o = {
        messageKey = 'message',
        appendLogKey = 'appendLog',
        targets = {},
        startTimer = true,
        uidKey = 'uid'
    }

    return setmetatable(o, {__index = _M})
end

function _M:init()
    for k, v in pairs(self.targets) do
        local target = manager:create_object(v)
        self.targets[k] = target
    end
end

-- 日志收集
function _M:log(message, tag, level)
    local uid   = ''
    local seqId = '-'
    local phase  = ngx.get_phase()
    if phase ~= 'timer' then
        uid   = get_manager(self):get('request'):get_param('login_uid')
        if util_value.empty(uid) then
            local ok, headers   = pcall(ngx.req.get_headers)
            if ok then
                uid             = array.get_value(headers, 'x-login-uid', '')
            end
        end
        if util_value.empty(uid) then
            uid     = get_manager(self):get('request'):get_param(self.uidKey)
        end
        seqId = _seq_id(self)
    end

    local messageInfo = {
        message = message,
        level = level,
        tag = tag,
        timestamp = time(),
        uid = uid,
        seqId = seqId,
    }
    local ngx_ctx = ngx.ctx
    local messages = ngx_ctx[self.messageKey] or new_table(16, 0)
    messages[#messages + 1] = messageInfo

    if not ngx_ctx[self.messageKey]  then ngx_ctx[self.messageKey] = messages end
end

-- 日志追加
function _M:append_log(message, tag, level)
    if message == nil then return end
    local ngx_ctx = ngx.ctx
    local appendLogs = ngx_ctx[self.appendLogKey] or new_table(0, 8)
    if not appendLogs[tag] then
        appendLogs[tag] = new_table(0, 8)
    end

    if not appendLogs[tag][level] then
        appendLogs[tag][level] = ''
    end

    appendLogs[tag][level] = appendLogs[tag][level] .. message

    if not ngx_ctx[self.appendLogKey] then ngx_ctx[self.appendLogKey] = appendLogs end
end

-- 日志刷出
function _M:flush()
    _flush_append_log(self)
    local ngx_ctx = ngx.ctx
    local messages = ngx_ctx[self.messageKey]
    if 'table' == type(messages) then
        for _, target in pairs(self.targets) do
            if self.startTimer then
                ngxTimer(
                    0,
                    function(_, _target, _messages)
                        _target:collect(_messages)
                    end,
                    target,
                    messages
                )
            else
                target:collect(messages)
            end
        end
    end

    ngx_ctx[self.messageKey] = nil
end

return _M