-- cache file
local spawn        = ngx.thread.spawn
local wait         = ngx.thread.wait
local capture      = ngx.location.capture
local type         = type
local pairs        = pairs
local ipairs       = ipairs
local setmetatable = setmetatable

local arrayHelper  = require('zues.utils.ArrayHelper')
local valueHelper  = require('zues.utils.ValueHelper')
local stringHelper = require('zues.utils.StringHelper')
local confManager  = require('zues.component.Config')
local manager      = require('zues.component.Manager')

local explode      = stringHelper.explode
local get_value    = arrayHelper.get_value
local value_empty  = valueHelper.empty
local trim         = stringHelper.trim
local basename     = stringHelper.basename
local array_filter = arrayHelper.filter
local random       = stringHelper.random
local new_table    = arrayHelper.new_table

local function _file_name(self, config, params)
    local fileName = ''
    local callBack = get_value(config, 'generateFileName')
    if 'function' == type(callBack) then
        fileName = callBack(config, params)
    else
        local file = get_value(config, 'fileName', '')
        fileName = self.filePath .. '/' .. trim(file, '/')
    end

    return fileName
end

local function _temp_file(self, config, params)
    local file = _file_name(self, config, params)

    return self.tempFilePath .. '/' .. basename(file)
end

-- cache
local function _cache_key(self, key, params, cacheConfig)
    local config = self.fileConfig[key]
    local cacheKeyExt = get_value(config, 'cache.cacheKeyExt', {})
    local level = get_value(cacheConfig, 'level', 0)
    local cacheKey = get_value(config, 'cache.cacheKey', _file_name(self, config, params))
    cacheKey = cacheKey .. '_' .. level

    for _, item in ipairs(cacheKeyExt) do
        local value = get_value(params, item, '')
        cacheKey = cacheKey .. '_' ..value
    end

    return cacheKey
end

-- 获取当前服务等级
local function _current_level(self)
    local level = 0
    if 'function' == type(self['degradeFunc']) then
        level =  self['degradeFunc']() or 0
    end

    return level
end

local function _check_result(self, key, result)
    local callBack = get_value(self.fileConfig, key..'.cache.check')
    if 'function' == type(callBack) then
        if true ~= callBack(result) then return false end
    elseif value_empty(result, 'table') then
        return false
    end

    return true
end

-- 缓存获取结果
local function _result_from_cache(self, key ,params, force)
    local result
    local cacheConfig = get_value(self.fileConfig, key..'.cache')
    if true ~= self.enableCache or value_empty(cacheConfig, 'table') then return result end

    if force or _current_level(self) >= get_value(cacheConfig, 'level', 0) then
        local cacheKey = _cache_key(self, key, params, cacheConfig)
        result = self.cache:get(cacheKey)
        -- 检查数据完整性
        local callBack = get_value(cacheConfig, 'check')
        if 'function' == type(callBack) then
            if nil ~= result and true == callBack(result) then
                get_manager(self):get('log'):log({key = key, params = params, cache = true, result = result}, 'debug', 'Debug')
                return result
            end
        elseif not value_empty(result, 'table') then
            get_manager(self):get('log'):log({key = key, params = params, cache = true, result = result}, 'debug', 'Debug')
            return result
        end
    end
    result = nil

    return result
end

-- 缓存结果
local function _cache_result(self, key, params, result)
    local cacheConfig = get_value(self.fileConfig, key..'.cache')
    if true ~= self.enableCache or value_empty(cacheConfig, 'table') then return end

    -- 调用自定义的检测函数
--    if not _check_result(self, key, result) then return end

    local confLevel = get_value(cacheConfig, 'level', 0)
    local probability = get_value(cacheConfig, 'probability', 0)
    local randNum = random(1, self.randNumRange)

    if _current_level(self) >= confLevel or  randNum <= probability then
        local cacheKey = _cache_key(self, key, params,cacheConfig)
        local expire = get_value(cacheConfig, 'expire')
        self.cache:set(cacheKey, result, expire)
    end

    return
end

-- 获取内容
local function _content(self, key, params, other)
    local tmp = other.temp
    local data = {}
    local config = get_value(self.fileConfig, key, {})
    if value_empty(config, 'table') then return data end

    local fileName = _file_name(self, config, params)

    local res,err = capture(fileName,{})
    local content = res.body
    if err or not content then return data end
    data = explode("\n", content)
    if value_empty(data, 'table') and tmp then
        local tmpFile = _temp_file(self, config, params)
        res,err = capture(tmpFile,{})
        content = res.body
        data = explode("\n", content)
        data = array_filter(data)
    end

    if value_empty(data, 'table') then return {} end

    -- 数据的二次切割
    local delimiter = get_value(config, 'delimiter')
    if delimiter then
        for k, v in pairs(data) do
            v = explode(delimiter, v)
            data[k] = v
        end
    end

    local callBack = get_value(config, 'callBack')
    if 'function' == type(callBack) then
        data = callBack(data, params)
    end

    return data
end

local _M = {_VERSION = '0.1',}
local mt = {__index = _M}


function _M:new(params)
    params = params or {}
    local o = {
        filePath = params['filePath'] or '',
        tempFilePath = params['tempFilePath'] or '',
        fileConfig = params['fileConfig'] or {},
        cache = 'cache',
        randNumRange = 10000,
        enableCache = true,
    }

    return setmetatable(o, mt)
end

function _M:init()
    if 'string' == type(self.cache) then
        self.cache = get_manager(self):get('cache')
    elseif 'table' == type(self.cache) then
        self.cache = manager:create_object(self.cache)
    end

    if 'string' == type(self.fileConfig) then
        self.fileConfig = confManager:load(self.fileConfig)
    end

    if 'table' ~= type(self.cache) then self.enableCache = false end
end

-- 获取内容
function _M:content(key, params, other)
    other = other or {}
    local force = other.force
    local data
    if not force then
        data = _result_from_cache(self, key, params)
        if 'table' == type(data) then return data end
    end

    data = _content(self, key, params, other)

    if not _check_result(self, key, data) then
        data = _result_from_cache(self, key, params, true) or {}
    else
        _cache_result(self, key, params, data)
    end

    return data
end

function _M:m_content(reqs)
    local ret, reqsThread = new_table(4, 4), new_table(4, 4)
    for k, req in pairs(reqs) do
        reqsThread[k] = spawn(self.content, self, req.key, req.params, req.other)
    end

    for k, reqThread in pairs(reqsThread) do
        local _, res = wait(reqThread)
        ret[k] = res
    end

    return ret
end

return _M