local ngx          = ngx
local now          = ngx.now
local ngx_log      = ngx.log
local ngx_err      = ngx.ERR
local setmetatable = setmetatable
local format       = string.format
local mongol       = require("resty.mongol")
local confManager  = require('zues.component.Config')
local value        = require('zues.utils.ValueHelper')
local array        = require('zues.utils.ArrayHelper')

local _M           = {
    _VERSION = "0.1.0"
}
local mt           = { __index = _M }

local function _log(self, statusCode, err)
    statusCode, err = statusCode or 200, err or 'success'
    if err ~= 'success' then
        statusCode = 502
    end
    local cost = now() - self.start_time
    -- network:{protocol}|{urlKey}|{cost}|{nodeLog}|{len}|{tryTimes}|{statusCode}|{web_degrade}|msg:{err}
    local msg  = format(' network:%s|%s|%.3f|%s|%d|%d|%d|0|msg:%s',
        'mongo',
        self.busKey,
        cost,
        self.conf.port,
        0,
        1,
        statusCode,
        err
    )
    local log  = get_manager():get('log')
    log:append_log(msg, 'req', 'Info')
    log:log({
        type = 'mongo',
        conf = self.conf,
        sql = self.sql,
        busKey = self.busKey,
        timeout = self.timeout,
        cost = format('%.3f', cost),
        status = statusCode,
        err = err,
    }, 'debug', 'Debug')
end

--读取数据：cursor
local function _find(self, rule)
    local query = rule['query']
    local returnFields = rule['returnfields']
    local numEachQuery = rule['num_each_query']
    if not query then
        return nil, 'query is nil'
    end
    local cursor = self.col:find(query, returnFields, numEachQuery)

    return cursor
end

local function _exec(self, command, rule)
    if command == 'find' then
        return _find(self, rule)
    end
end


local function _keepalive(self, idle, size)
    idle, size = idle or self.conf.keepalive or 10000, size or self.conf.poolsize or 1000

    if not self['conn'] or self.is_closed then return nil end

    self.conn:set_keepalive(idle, size)
    _log(self, 200)

    return true
end

--连接数据库ok, err
local function _connect(host, port, database, collection, user, password, timeout, poolsize)
    local link = mongol:new()
    link:set_timeout(timeout)

    local pool = user .. ":" .. database .. ":" .. host .. ":" .. port
    local pool_opts = {
        pool = pool,
        pool_size = poolsize or 10,
    }

    local ok, err = link:connect(host or '', port or '', pool_opts)
    if err then
        ngx_log(ngx_err, 'failed to connect mongo host:', host, ' port:', tostring(port), err)
    end

    while not link:ismaster() do
        ok, err = link:connect(host or '', port or '')
        if not ok then
            return nil, nil, err
        end
    end

    local db = nil
    local times, err = link:get_reused_times()
    if 0 == times or nil == times then
        db = link:new_db_handle(database or '')
        if nil == db then
            ngx.log(ngx.ERR, "MongoDB new_db_handle failed 01.\n")
            -- DB_POOL.close_db_mongo_cli()
            link:close() -- 关闭此socket连接
            return nil, nil, 'MongoDB new_db_handle failed 01.'
        end

        --用户授权
        local ok, err = db:auth_scram_sha1(user or '', password or '')
        if ok ~= 1 then
            ngx.log(ngx.ERR, "MongoDB user auth failed, err=", err, ".\n")
            -- DB_POOL.close_db_mongo_cli()
            link:close() -- 关闭此socket连接
            return nil, nil, 'MongoDB auth_scram_sha1 failed: ' .. err .. " ###"
        end
    else
        ngx.log(ngx.ERR, '$$$ get_db_mongo_col() reused: ', times)
        db = link:new_db_handle(database or '')
    end

    if value.is_null(collection) then
        return nil, nil, 'collection err!'
    end

    return link, db:get_col(collection), nil
end

function _M:new()
    return setmetatable({
        timeout = 100,
        confKey = 'mongo',
        conf = 'resource',
    }, mt)
end

function _M:init()
    if 'string' == type(self.conf) then
        self.conf = confManager:load(self.conf)
    end

    self.conf = self.conf[self.confKey]
end

function _M:get_connect(args)
    local busKey = args['busKey']
    if not busKey then
        return nil, 'arg busKey not found!'
    end
    local conf = self.conf[busKey]
    if 'table' ~= type(conf) then
        return nil, 'node not found!'
    end

    local collection = conf.collection[array.get_value(args, 'colKey', '')]
    if not collection then
        return nil, 'arg collection not found!'
    end

    local start_time      = now()
    local timeout         = args['timeout'] or conf['timeout'] or self.timeout

    local mongo, col, err = _connect(conf.host, conf.port, conf.database, collection, conf.user, conf.password, timeout,
        conf.poolsize)

    local mongo_obj       = {
        conf = conf,
        busKey = busKey,
        start_time = start_time,
        timeout = timeout,
        exec = _exec,
        keepalive = _keepalive,
        conn = mongo,
        col = col
    }
    if err then
        _log(mongo_obj, 502, err)
        return nil, err
    end

    return mongo_obj
end

return _M
