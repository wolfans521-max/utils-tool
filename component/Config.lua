-- 配置文件参数获取
local ngx                   = ngx
local pairs                 = pairs
local pcall                 = pcall
local require               = require
local setmetatable          = setmetatable


local arrayHelper           = require('zues.utils.ArrayHelper')
local valueHelper           = require('zues.utils.ValueHelper')

local value_copy  = valueHelper.copy
local array_merge = arrayHelper.merge
local get_value   = arrayHelper.get_value
local new_table   = arrayHelper.new_table

local _M = {_VERSION = '0.1', config = new_table(0, 0)}
local mt = {__index = _M}

-- 配置文件加载
function _M:load(file, module)
    module = module or {'common', ngx.ctx.APPLICATION_NAME or APPLICATION_NAME or '', 'wbutil'}
    local key = module[2] .. '_' .. file
    local config = self.config[key]
    if config then return config end

    for _, m in pairs(module) do
        local mdl = 'conf.' .. file
        if m and m ~= '' then mdl = m .. '.' .. mdl end
        local ok, segment = pcall(require, mdl)
        if ok then
            if not config then
                config = value_copy(segment)
            else
                config = array_merge(config, value_copy(segment))
            end
        end
    end

    self.config[key] = config

    return config
end

-- 获取配置文件
function _M:get(name, key, default)
    local config = self:load(name)
    if not config then
        return default
    end

    if not key then return config end
    local value = get_value(config, key, default)
    if not value then return default end

    return value
end

function _M:new()
    return setmetatable({config = new_table(0, 0)}, mt)
end

return _M