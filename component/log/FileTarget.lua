--文件类
local date                  = os.date
local worker_id             = ngx.worker.id
local popen                 = io.popen
local io_open               = io.open
local os_execute            = os.execute
local tonumber              = tonumber
local pairs                 = pairs
local type                  = type
local setmetatable          = setmetatable

local arrayHelper           = require('zues.utils.ArrayHelper')
local stringHelper          = require('zues.utils.StringHelper')
local util_file             = require('zues.utils.FileHelper')
local cjson                 = require('cjson')

local get_value             = arrayHelper.get_value
local random                = stringHelper.random
local in_array              = arrayHelper.in_array
local json_encode           = cjson.encode
local file_put_contents     = util_file.file_put_contents
local new_table             = arrayHelper.new_table

local shell_cut    = false
local ok, shell    = pcall(require, 'resty.shell')
if ok then shell_cut = true end

local function _check_file_size(self, fileName, name)
    if random(1, self.randNum) == 1 then
        local cmd = 'ls -l '..fileName..'|awk \'{print $5}\''
        local size
        if shell_cut then
            local _, file_size = shell.run(cmd)
            size = tonumber(file_size)
        else
            local fileObj = popen(cmd)
            size = tonumber(fileObj:read('*all'))
            fileObj:close()
        end
        if size and size >= self.maxSize*1024*1024*1024 then
            local newFile = self.filePath .. '/' .. self.commonPrefix .. self.delimiter .. name .. self.delimiter .. date(self.postFixCut) .. self.extension
            local cmd = 'mv '..fileName..' '..newFile
            if shell_cut then
                shell.run(cmd)
            else
                os_execute(cmd)
            end

            return true
        end
    end

    return false
end

-- 获取日志文件名称
local function _file_name(self, tag)
    local name = self.tag2finename[tag] or tag
    local fileName = self.commonPrefix .. self.delimiter .. name .. self.delimiter .. date(self.postFix) .. self.extension
    fileName = self.filePath .. '/' .. fileName
    local format = self.commonPrefix .. self.delimiter .. name .. self.delimiter .. self.postFix .. self.extension
    local cuted = false

    -- 检查文件大小
    if self.autoCut and worker_id() == 1 then
        cuted = _check_file_size(self, fileName, name)
    end

    return self.filePath .. '/' .. format, fileName, cuted
end

-- 格式化日志时间
local function _format_date(timestamp)
    return date("%Y-%m-%d %H:%M:%S", timestamp)
end

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
    local dateTime = _format_date(message['timestamp'])
    local seqId = message.seqId or '-'

    return dateTime.."\t"..seqId.."\t"..message.level.."\t"..message.tag.."\t"..message.message
    --return format("%s\t%s\t%s\t%s\%s", dateTime, seqId, message['level'], message['tag'], message['message'])
end

-- 内容导出
local function _export(self, messages)
    local _messages = new_table(0, 16)
    for _, message in pairs(messages) do
        local tag = message['tag']
        local text = _format_message(self, message)
        _messages[tag] = get_value(_messages, tag, '') .. text .. "\n"
    end

    for tag, msg in pairs(_messages) do
        local fileFormat, fileName, cuted = _file_name(self, tag)
        if self.startCache then
            local openfile, fd = self.openfiles[fileFormat] or {}, nil
            if cuted or not openfile['fd'] or openfile['fileName'] ~= fileName then
                if openfile['fd'] then openfile['fd']:close() end
                local _fd, err = io_open(fileName, 'a+')
                if not err and _fd then
                    _fd:setvbuf('no')
                    self.openfiles[fileFormat] = {fd = _fd, fileName = fileName}
                    fd = _fd
                end
            else
                fd = openfile['fd']
            end
            if fd then
                fd:write(msg)
                fd:flush()
            end
        else
            file_put_contents(fileName, msg)
        end
    end
end

local _M = {_VERSION = '0.1', openfiles = {}}

function _M:new()
    local o = {
        tags = {},
        levels = {'Error', 'Info', 'Warning', 'Trace'},
        enable = true,
        delimiter = '_',
        postFix = '%Y%m%d',
        commonPrefix = '',
        filePath = '',
        extension = '.log',
        tag2finename = {},
        maxSize = 2, --单位G
        autoCut = true,
        randNum = 1000,
        startCache = false,
        postFixCut = '%Y%m%d%H%M'
    }

    return setmetatable(o, { __index = _M })
end

-- 日志收集
function _M:collect(messages)

    local messages = _message_filter(self, messages)
    if #messages < 1 then return end

    _export(self, messages)
end

return _M