--请求处理类
local ngx_gsub                  = ngx.re.gsub
local require                   = require
local ipairs                    = ipairs
local type                      = type
local pcall                     = pcall
local setmetatable              = setmetatable
local pairs                     = pairs

local arrayHelper               = require('zues.utils.ArrayHelper')
local valueHelper               = require('zues.utils.ValueHelper')

local get_value   = arrayHelper.get_value
local set_value   = arrayHelper.set_value
local value_empty = valueHelper.empty
local value_copy  = valueHelper.copy
local array_merge = arrayHelper.merge
local new_table   = arrayHelper.new_table

local function _load(self, category, language)
    local key = language..'.'..category
    local lang = get_value(self._lang, key)
    if 'table' == type(lang) then return lang end

    for _, path in ipairs(self.languagePath) do
        local lPath = path .. '.' .. language .. '.' .. category
        local ok, segment = pcall(require, lPath)
        if ok and 'table' == type(segment) then
            if not lang then
                lang = value_copy(segment)
            else
                lang = array_merge(lang, value_copy(segment))
            end
        end
    end
    set_value(self._lang, key, lang)

    return lang
end

local _M = {_VERSION='0.01'}
local mt = {__index = _M }


function _M:new()
    local o = {
        languagePath = {},
        defaultLanguage = 'zh_CN',
        _lang = new_table(0, 0),
        needTrans = true, -- 默认语言时是否需要翻译
    }

    return setmetatable(o, mt)
end

function _M:translate(category, message, params, language)
    language = language or self.defaultLanguage

    local lang
    if language ~= self.defaultLanguage or self.needTrans then
        lang = _load(self, category, language)
    end

    message = get_value(lang, message, message)

    if value_empty(params, 'table') then return message end

    for k, value in pairs(params) do
        message = ngx_gsub(message, '{'..k..'}', value)
    end

    return message
end

return _M