local table_sort                = table.sort
local ngx                       = ngx
local package                   = package
local pairs                     = pairs

local jsonHelper                = require('zues.utils.JsonHelper')
local arrayHelper               = require('zues.utils.ArrayHelper')

local new_table     = arrayHelper.new_table
require('resty.core.shdict')

local json_decode = jsonHelper.decode

local _M = {_VERSION='0.01'}

--查看系统状态
function _M:system_info()
    local ret = {
        --共享内存信息
        shared_dict = new_table(0, 16),
        --timers
        timers = {
            running = ngx.timer.running_count(),
            pedding = ngx.timer.pending_count(),
        },
        config = {
            prefix = ngx.config.prefix(),
            debug = ngx.config.debug,
            ngx_version = ngx.config.nginx_version,
            lua_version = ngx.config.ngx_lua_version,
            ngx_configure = ngx.config.nginx_configure(),
            subsystem = ngx.config.subsystem,
        },
        path = {
            path = package.path,
            cpath = package.cpath,
        },
        worker_count = ngx.worker.count(),
        current_worker_pid = ngx.worker.pid(),
        current_worker_id = ngx.worker.id(),
    }
    for k, _ in pairs(ngx.shared) do
        local shareCache = ngx.shared[k]
        --共享内存信息
        ret['shared_dict'][k] = {
            total = shareCache:capacity() / 1024 / 1024 .. 'M',
            free = shareCache:free_space() / 1024 / 1024 .. 'M',
        }
    end

    return ret
end

function _M:show_cache(params)
    params = params or {}
    local ret = new_table(0, 0)
    local cacheKey = params['cacheKey']
    local dictKey = params['dictKey'] or 'my_cache'
    local shareCache = ngx.shared[dictKey]

    if cacheKey then
        local res = shareCache:get(cacheKey)
        if res then res = json_decode(res) end
        ret[cacheKey] = {value = res, ttl = shareCache:ttl(cacheKey)}
    else
        local keys = shareCache:get_keys(0)
        for _, key in pairs(keys) do
            local res = json_decode(shareCache:get(key))
            ret[key] = {value = res, ttl = shareCache:ttl(key)}
        end
    end

    return ret
end

function _M:clear_cache(params)
    params = params or {}
    local ret = new_table(0, 0)
    local cacheKey = params['cacheKey']
    local dictKey = params['dictKey'] or ''
    local shareCache = ngx.shared[dictKey]

    if cacheKey then
        ret[cacheKey] = shareCache:delete(cacheKey)
    else
        local keys = shareCache:get_keys(0) or {}
        table_sort(keys)
        ret['keys'] = keys
        shareCache:flush_all()
        local totalMem = shareCache:capacity()
        local free = shareCache:free_space()
        ret['delete_num'] = shareCache:flush_expired()
        ret['total'] = totalMem .. 'B/'..totalMem / 1024 / 1024 ..'M'
        ret['free'] =  free .. 'B/' .. free / 1024 / 1024 .. 'M'
    end

    return ret
end

-- 更新cache接口
function _M:update_cache(params)
    params = params or {}
    local ret = {code = 200, msg = 'ok'}
    local params = self:_get_param()
    local cache = self.app:get('config_cache')
    local items = json_decode(params['items']) or {}

    for key, value in pairs(items) do
    cache:set(key, value, nil, false)
    end

    return ret
end

return _M