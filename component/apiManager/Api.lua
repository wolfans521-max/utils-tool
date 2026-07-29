local format        = string.format
local lower         = string.lower
local ngx_find      = ngx.re.find
local str_sub       = string.sub
local spawn         = ngx.thread.spawn
local wait          = ngx.thread.wait
local encode_base64 = ngx.encode_base64
local escape_uri    = ngx.escape_uri
local hmac_sha1     = ngx.hmac_sha1
local pow           = math.pow
local floor         = math.floor
local type          = type
local tostring      = tostring
local tonumber      = tonumber
local ngx           = ngx
local ipairs        = ipairs
local pairs         = pairs
local setmetatable  = setmetatable
local table_concat  = table.concat

local arrayHelper   = require('zues.utils.ArrayHelper')
local valueHelper   = require('zues.utils.ValueHelper')
local jsonHelper    = require('zues.utils.JsonHelper')
local stringHelper  = require('zues.utils.StringHelper')
local confManager   = require('zues.component.Config')
local manager       = require('zues.component.Manager')

local value_copy    = valueHelper.copy
local value_empty   = valueHelper.empty
local get_value     = arrayHelper.get_value
local simple_merge  = arrayHelper.simple_merge
local json_decode   = jsonHelper.decode
local explode       = stringHelper.explode
local random        = stringHelper.random
local in_array      = arrayHelper.in_array
local new_table      = arrayHelper.new_table

local function _get_token(self, uid)
    uid = uid or ''
    local authHead = ''
    local res,err = ngx.location.capture(self.tokenFile,{})
    if err then return authHead end

    local token = json_decode(res.body) or {}
    if not token['tauth_token'] or not token['tauth_token_secret'] then
        return authHead
    end
    if not uid then
        authHead = 'TAuth2 token=' .. escape_uri(token['tauth_token']);
    else
        local param  = 'uid=' .. uid;
        local sign   = encode_base64(hmac_sha1(token['tauth_token_secret'], param));
        authHead = 'TAuth2 token="' .. escape_uri(token['tauth_token']) .. '",param="' .. escape_uri(param) .. '",sign="' .. escape_uri(sign) .. '"';
    end

    return authHead
end

local function get_cace_instance(self, cache)
    local instance
    if 'string' == type(cache) then
        instance = get_manager(self):get(cache)
    elseif 'table' == type(cache) then
        instance = manager:create_object(cache)
    end
    if 'table' ~= type(instance) then instance = self.cache end

    return instance
end

-- 结果干预
local function _intervene_result(self, option, result)
    local config = self.apisConfig[option['urlKey']]
    local callBack = get_value(config, 'interveneResult')
    if 'function' == type(callBack) then
        result = callBack(option, result)
    end

    return result or {}
end

-- 缓存结果干预
local function _cache_result_intervene(self, option, result)
    local config = self.apisConfig[option['urlKey']]
    local callBack = get_value(config, 'cache.intervene')
    if 'function' == type(callBack) then
        result = callBack(option, result)
    end

    return result or {}
end

-- 获取当前服务等级
local function _current_level(self, key)
    local level = 0
    if key then key = self.degrade_key_prefix .. key end
    if 'function' == type(self['degradeFunc']) then
        level = self['degradeFunc'](key) or 0
    end

    return level
end

-- 获取当前请求seqid
local function _seq_id(self)
    local seqHandler = self.seqHandler
    if type(seqHandler) == 'function' then return seqHandler() end

    local manager = get_manager()
    if manager then
        return manager:get('request'):seqid()
    end

    return nil
end

local function _parse_group(self, group)
    if not group then return group end

    local env_group = self.env_group

    return format(group, env_group)
end

-- 协议解析
local function _parse_url(self, apiConfig, otherParam)
    local url            = apiConfig.url
    local gray_url       = apiConfig.gray_url or ''
    local proxy_host     = otherParam.proxy_host or ''
    --命中灰度使用灰度url
    if not value_empty(proxy_host, 'string') and not value_empty(gray_url, 'string') then
        url              = format(gray_url, proxy_host)
    end
    url                  = _parse_group(self, url)
    -- ocal opt = {url = url, protocol = 'http', service = apiConfig.service, group = apiConfig.group}
    local opt = new_table(0, 18)
    opt.url = url
    opt.protocol = 'http'
    opt.service = apiConfig.service
    opt.group = apiConfig.group
    local from = ngx_find(url, '://', 'jo')
    local protocol = lower(str_sub(url, 1, from - 1))
    -- http协议
    if protocol == 'http' or protocol == 'https' then
        if opt.service and opt.group then
            opt.group = _parse_group(self, opt.group)
            opt.nodes = self.discovery:naming_get(opt.service, opt.group)
        end
        return opt
    end
    local other = explode('/', str_sub(url, from+3))
    if protocol == 'motan2' then
        opt.protocol = 'motan'
        opt.service = other[2]
        opt.group = other[3]
        -- opt.path = '/'..table_concat(other, '/', 4)
    elseif protocol == 'da' then
        opt.protocol = 'da'
        opt.service = other[1]
        opt.group = other[2]
    end
    --opt.group = _parse_group(self, opt.group)
    opt.nodeLog = opt.group
    if protocol == 'da' then
        opt.nodes = self.discovery:naming_get(opt.service, opt.group)
    end

    return opt
end

local function _check_result(self, urlKey, result)
    local callBack = get_value(self.apisConfig, urlKey .. '.cache.check')
    if 'function' == type(callBack) then
        if true ~= callBack(result) then return false end
    elseif value_empty(result, 'table') then
        return false
    end

    return true
end

local function _cache_level_expire(self, urlKey, cacheConfig)
    local callBack = get_value(self.apisConfig, urlKey .. '.cache.cacheKeyExpire')
    local level = get_value(cacheConfig, 'level', 0)
    local expire = get_value(cacheConfig, 'expire', 0)
    if 'function' == type(callBack) then
        level, expire = callBack()
    end

    return level, expire
end

-- 获取缓存key
local function _cache_key(self, option, config)
    local level, _ = _cache_level_expire(self, option['urlKey'], config)
    local cacheKey = new_table(4, 0)
    cacheKey[1] = get_value(config, 'cacheKey', option['urlKey'])
    cacheKey[2] = level
    local cacheKeyExt = get_value(config, 'cacheKeyExt', {})
    local args, index = option['args'] or {}, 3
    for _,arg in ipairs(cacheKeyExt) do
        if args[arg] then
            cacheKey[index] = option['args'][arg]
            index = index + 1
        end
    end

    return table_concat(cacheKey, '_')
end

-- 根据uid指定node
local function _parse_nodes(self, nodes, uid)
    local _nodes, gray_nodes = new_table(32, 0), new_table(8, 0)
    for _, node in ipairs(nodes) do
        local rule = get_value(node,'extInfo.rule')
        if not value_empty(rule, 'table') then
            -- 判断是否命中灰度
            local uids = get_value(rule, 'uids')
            local pos = tonumber(get_value(rule,'pos')) or 1
            local num = tonumber(get_value(rule, 'num')) or 1
            local value = get_value(rule, 'value', {})
            local offset = pow(10, pos - num)
            local posNum = floor((uid / offset)) % pow(10, num)
            if not value_empty(uids, 'table') and in_array(uid, uids) then
                gray_nodes = {node}
                break
            elseif in_array(posNum, value) then
                gray_nodes[#gray_nodes + 1] = node
            end
        else
            _nodes[#_nodes + 1] = node
        end
    end

    -- 灰度命中检测
    if #gray_nodes > 0 then
        return gray_nodes
    elseif #_nodes > 0 then
        return _nodes
    end

    return nodes
end

-- 设置请求
local function _set_request(self, urlKey, params, other)
    if not urlKey then return false end
    other, params = other or {}, params or {}
    local header, timeout = other.header or {}, other.timeout

    local apiConfig = get_value(self.apisConfig, urlKey, {})
    if type(apiConfig) ~= 'table' or not apiConfig['url'] then return false end

    local args = apiConfig.args
    if value_empty(args, 'table') then
        args = params
    else
        args = value_copy(args)
        for k, v in pairs(args) do
            args[k] = get_value(params, k) or v
        end
    end

    if not apiConfig.cancel_ext_arg then
        -- 测试环境标记请求流量为压测环境
        if self.env == 'dev' then args.benckmark = '1' end
        args.web_degrade = params.web_degrade or tostring(_current_level(self, urlKey))
        args.seqid = params.seqid or tostring(_seq_id(self))
    end
    local opt = _parse_url(self, apiConfig, other)
    opt['timeout'] = timeout or get_value(apiConfig,'timeout')
    opt['method'] = get_value(apiConfig,'method') or self.defaultMethod
    opt['urlKey'] = urlKey
    opt['tryTimes'] = apiConfig['tryTimes']
    opt['args'] = args
    opt['keepalive'] = apiConfig['keepalive']
    opt['contentJson'] = apiConfig['contentJson']
    opt['web_degrade'] = args.web_degrade
    opt['version'] = apiConfig['version'] or 1
    -- 根据uid指定下游node信息
    local uid = tonumber(args['uid'])
    if not value_empty(opt.nodes, 'table') and uid then
        opt.origin_nodes = opt.nodes
        opt.nodes = _parse_nodes(self, opt.nodes, uid)
    end

    local confHeader = get_value(apiConfig, 'header', {}, true)
    header = simple_merge(confHeader,header)
    if true == get_value(apiConfig, 'auth') then
        opt['args']['source'] = self.source
        local cuid = args['uid']
        if not value_empty(args['cuid'], 'string')  or not value_empty(args['cuid'], 'number') then
            cuid = args['cuid']
        end
        local token = _get_token(self, cuid)
        header['Authorization'] = token
        if opt.protocol == 'da' then
            args['token'] = token
        end
    end
    opt['header'] = header

    return opt
end

local function _enable_cache(self, urlKey, result)
    local callBack = get_value(self.apisConfig, urlKey .. '.cache.enableCb')
    if 'function' == type(callBack) then
        if true ~= callBack(result) then return false end
    end

    return true
end

-- 缓存获取结果
local function _result_from_cache(self, option, force)
    local result
    local cacheConfig = get_value(self.apisConfig, option['urlKey']..'.cache')
    if true ~= self.enableCache or value_empty(cacheConfig, 'table') then return result end
    if force or _current_level(self, option['urlKey']) >= get_value(cacheConfig, 'level', 0) then
        local cacheKey = _cache_key(self, option, cacheConfig)
        local cache_instance = get_cace_instance(self, get_value(cacheConfig, 'instance'))
        result = cache_instance:get(cacheKey)
        -- 检查数据完整性
        if result and _check_result(self, option['urlKey'], result) then
            get_manager(self):get('log'):log({option = option, cache = true, result = result}, 'debug', 'Debug')

            return _cache_result_intervene(self, option, result)
        end
        result = nil
    end

    return result
end

-- 缓存结果
local function _cache_result(self, option, result)
    local cacheConfig = get_value(self.apisConfig, option['urlKey']..'.cache')
    if true ~= self.enableCache or value_empty(cacheConfig, 'table') then return end

    local cache_enable = _enable_cache(self, option['urlKey'], result)
    if not cache_enable then
        return true
    end

    --调用自定义的检测函数
--    if not _check_result(self, option['urlKey'], result) then return end
    local confLevel = get_value(cacheConfig, 'level', 0)
    local probability = get_value(cacheConfig, 'probability', 0)
    local randNum = random(1, self.randNumRange)

    if _current_level(self, option['urlKey']) >= confLevel or randNum <= probability then
        local cacheKey = _cache_key(self, option, cacheConfig)
        local _, expire = _cache_level_expire(self, option['urlKey'], cacheConfig)
        local cache_instance = get_cace_instance(self, get_value(cacheConfig, 'instance'))
        cache_instance:set(cacheKey, result, expire)
    end

    return true
end

-- 请求
local function _request(self, opt)
--    local opt = _set_request(self, urlKey, params, other)
    local protocol = opt.protocol
    local res, err = self.clients[protocol]:s_curl(opt)

    return _intervene_result(self, opt, json_decode(res)), err
end

local _M = {_VERSION = '0.1'}
local mt = {__index = _M}

function _M:new()
    return setmetatable({
        apisConfig = 'api',
        source = '2936099636',
        tokenFile = '',
        defaultMethod = 'GET',
        log = 'log',
        discovery = 'discovery',
        randNumRange = 10000,
        enableCache = true,
        env_group = 'yf',
        degrade_key_prefix = 'api_degrade_',
        clients = {
            http = {class = 'zues.component.client.HttpClient'},
            da = {class = 'zues.component.client.DaClient'},
            motan = {class = 'zues.component.client.MotanClient'},
        },
    }, mt)
end

-- 初始化
function _M:init()
    local get_manager = get_manager
    -- 接口配置文件
    if 'string' == type(self.apisConfig) then
        local config = self.apisConfig
        self.apisConfig = confManager:load(config)
    end
    -- client
    for k, clientConf in pairs(self.clients) do
        self.clients[k] = manager:create_object(clientConf)
    end
    -- 日志
    if type(self.log) == 'string' then self.log = get_manager(self):get(self.log) end
    --服务发现组件
    if type(self.discovery) == 'string' then self.discovery = get_manager(self):get(self.discovery) end
    
    if 'string' == type(self.cache) then
        self.cache = get_manager(self):get('cache')
    elseif 'table' == type(self.cache) then
        self.cache = manager:create_object(self.cache)
    end
    if 'table' ~= type(self.cache) then self.enableCache = false end
end

--返回认证信息
function _M:get_token(uid)
    return _get_token(self,uid)
end
-- 获取请求结果
function _M:s_request(urlKey, params, other)
    local result, err
    local get_manager = get_manager
    other = other or {}
    local force = other.force
    local opt = _set_request(self, urlKey, params, other)
    if opt == false then
        return {}
    end
    if not force then
        result = _result_from_cache(self, opt)
        if 'table' == type(result) then
            get_manager(self):get('log'):append_log(format(' custom_counter:%s_total-%s_cache', urlKey, urlKey), 'req', 'Info')
            return result
        end
    end
    result, err = _request(self, opt)
    local log = format(' custom_counter:%s_total', urlKey)
    local check = _check_result(self, urlKey, result)

    if not err and check then
        _cache_result(self, opt, result)
    else
        result = _result_from_cache(self, opt, true)
        if 'table' == type(result) then
            log = log..'-'..urlKey..'_cache'
        else
            result = {}
        end
    end

    get_manager(self):get('log'):append_log(log, 'req', 'Info')

    return result
end

--批量请求
function _M:m_request(reqs)
    local ret, reqsThread = new_table(4, 4), new_table(4, 4)
    for k, req in pairs(reqs) do
        reqsThread[k] = spawn(self.s_request, self, req.urlKey, req.params, req.other)
    end

    for k, reqThread in pairs(reqsThread) do
        local _, res = wait(reqThread)
        ret[k] = res
    end

    return ret
end

return _M