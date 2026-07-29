local ngx                       = ngx
local setmetatable              = setmetatable
local type                      = type
local require                   = require
local pairs                     = pairs

-- 组建管理器
local valueHelper               = require('zues.utils.ValueHelper')
local config                    = require('zues.component.Config')
local arrayHelper               = require('zues.utils.ArrayHelper')

local value_copy  = valueHelper.copy
local new_table   = arrayHelper.new_table

function get_manager(self)
    if self and self.zues_manager then
        return self.zues_manager
    elseif ngx.ctx.zues_manager then
        return ngx.ctx.zues_manager
    end

    return GLOBAL_ZUES_MANAGER
end

local _M = {_VERSION='01'}
local mt = {__index = _M}

function _M:new(o)
    if 'string' == type(o) then o = config:load(o) end
    -- 加载核心组件
    if not o.request then o.request = {class = 'zues.component.request.Http'} end
    local obj = {components = o, componentInstances = new_table(0, 8)}

    return setmetatable(obj, mt)
end

-- 创建组建
function _M:create_object(config)
    -- 非数组
    if 'table' ~= type(config) or not config['class'] then return nil end

    -- module不存在
    local className = require(config['class'])
    local params = config['attrs'] or {}
    local object = className:new(params)
    self:configure(object, params)
    if 'function' == type(object['init']) then
        object:init() -- 对象初始化方法
    end

    return object
end

-- 配置属性
function _M:configure(object, params)
    params = params or {}

    for k, v in pairs(params) do
        object[k] = v
    end
end

-- 获取组建
function _M:get(name)
    if not name then return self.components end

    -- 判断是否有实例化
    local component = self.componentInstances[name]
    if component then return component end
    local componentConf = self.components[name]
    if not componentConf then return nil end

    component = self:create_object(value_copy(componentConf))

    if type(component) == 'table' then
        component.zues_manager = self
        self.componentInstances[name] = component
    end

    return component
end

return _M