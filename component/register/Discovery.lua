local getmetatable                  = getmetatable
local setmetatable                  = setmetatable
local type                          = type
local ipairs                        = ipairs
local tonumber                      = tonumber

local parent                        = require('zues.component.cache.Cache')
local jsonHelper                    = require('zues.utils.JsonHelper')
local discovery                     = require('resty.discovery')
local stringHelper                  = require('zues.utils.StringHelper')
local arrayHelper                   = require('zues.utils.ArrayHelper')
local confManager                   = require('zues.component.Config')

local config_get   = discovery.configget
local naming_get   = discovery.namingget
local json_decode  = jsonHelper.decode
local explode      = stringHelper.explode
local get_value    = arrayHelper.get_value
local new_table    = arrayHelper.new_table

local discovery_node_list = new_table(0, 32)
local discovery_config_list = new_table(0, 32)

local _M = {}
local mt = {__index = _M}

function _M:new()
    local obj = parent:new()
    local super_mt = getmetatable(obj)
    -- 当方法在子类中查询不到时，再去父类中去查找。
    setmetatable(_M, super_mt)
    -- 这样设置后，可以通过self.super.method(self, ...) 调用父类的已被覆盖的方法。
    obj.super = setmetatable({
        group = '',
        defaultConfig = nil,
    }, super_mt)

    return setmetatable(obj, mt)
end

function _M:init()
    if 'string' == type(self.defaultConfig) then
        self.defaultConfig = confManager:load(self.defaultConfig)
    end
end

-- 获取值
function _M:config_get(key, default, needSerise, group)
    group = group or self.group
    if needSerise == nil then needSerise = true end
    local cacheKey, sign
    if needSerise then
        cacheKey = group .. '_' .. key
        local sign = discovery.configgetsign(group, key)
        local old_item = discovery_config_list[key] or {}
        local old_sign = old_item['sign']
        local old_value = old_item['value']
        if sign and sign == old_sign and old_value then
            return old_value
        end
    end
    local config = config_get(group, key)

    if not config then return default end

    if needSerise then
        config = json_decode(config)
        discovery_config_list[key] = {sign = sign, value = config}
    end

    return config
end

function _M:naming_get(service, cluster)
    -- 获取sign值
    local key = service..'_'..cluster
    local sign = discovery.naminggetsign(service, cluster)
    local old_item = discovery_node_list[key] or {}
    local old_sign = old_item['sign']
    local old_nodes = old_item['nodes']
    -- 数据无更新
    if sign and sign == old_sign
            and 'table' == type(old_nodes) then
        return old_nodes
    end

    local nodes = get_value(json_decode(naming_get(service, cluster)), 'working')
    if 'table' == type(nodes) and #nodes > 0 then
        for _, item in ipairs(nodes) do
            local extInfo = {}
            if type(item.extInfo) == 'string' and #item.extInfo > 0 then
                extInfo = json_decode(item.extInfo)
                if type(extInfo) ~= 'table' then extInfo = {} end
            end
            item.extInfo = extInfo
            item.weight = extInfo.weight or 1
            item.domain = item.host
            local hostInfo = explode(':', item.host)
            item.host = hostInfo[1]
            item.port = tonumber(hostInfo[2])
        end
        discovery_node_list[key] = {sign = sign, nodes = nodes }

        return nodes
    elseif 'table' == type(old_nodes) then -- 兼容注册中心获取失败的情况使用老数据
        return old_nodes
    end

    -- 获取默认配置
    if 'table' ~= type(self.defaultConfig) then return {} end
    nodes = self.defaultConfig[service]
    if nodes then return nodes[cluster] or {} end

    return {}
end

return _M