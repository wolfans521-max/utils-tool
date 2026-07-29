local ngx         = ngx
local ngx_log     = ngx.log
local now         = ngx.now
local ngx_err     = ngx.ERR
local format      = string.format
local table_concat = table.concat

local value       = require('zues.utils.ValueHelper')
local array       = require('zues.utils.ArrayHelper')
local mysql       = require('resty.mysql')
local confManager = require('zues.component.Config')
local _M          = { __VERSION = 0.1 }

local mt          = { __index = _M }

local function _log(self, statusCode, err)
    statusCode, err = statusCode or 200, err or 'success'
    if err ~= 'success' then
        statusCode = 502
    end
    local cost = now() - self.start_time
    -- network:{protocol}|{urlKey}|{cost}|{nodeLog}|{len}|{tryTimes}|{statusCode}|{web_degrade}|msg:{err}
    local msg  = format(' network:%s|%s|%.3f|%s|%d|%d|%d|0|msg:%s',
            'mysql',
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
        type = 'mysql',
        conf = self.conf,
        sql = self.sql,
        busKey = self.busKey,
        timeout = self.timeout,
        cost = format('%.3f', cost),
        status = statusCode,
        err = err,
    }, 'debug', 'Debug')
end

--连接数据库ok, err, errcode, sqlstate
local function _connect(host, port, database, user, charset, password, timeout)
    local link = mysql:new()
    link:set_timeout(timeout)
    local ok, err, errcode, sqlstate = link:connect({
        host = host or '',
        port = port or '',
        database = database or '',
        user = user or '',
        charset = charset or '',
        password = password or ''
    })
    if err then
        ngx_log(ngx_err, 'failed to connect mysql host:', host, ' port:', tostring(port), err)
    end

    return link, err, errcode, sqlstate
end

--解析具体的值
local function _fun_ops(key, val)
    if type(val) == 'table' then
        local str = '';
        for _, vval in ipairs(val) do
            if str == '' then
                str = ngx.quote_sql_str(vval)
            else
                str = str .. ',' .. ngx.quote_sql_str(vval)
            end
        end
        return key .. ' in (' .. str .. ')'
    elseif type(val) == 'string' or type(val) == 'number' then
        return key .. ' = ' .. ngx.quote_sql_str(val)
    end
end

--解析嵌套
local function _fun_exp(args, op)
    local where = '';
    if value.is_null(op) or op == '' or op == '_and' then
        op = 'and';
    else
        op = 'or'
    end
    if args['_or'] then
        where       = where .. _fun_exp(args['_or'], '_or')
        args['_or'] = nil
    end
    if args['_and'] then
        if where == '' then
            where = where .. _fun_exp(args['_and'], '_and')
        else
            where = where .. ' and ' .. _fun_exp(args['_and'], '_and')
        end
        args['_and'] = nil
    end
    for key, val in pairs(args) do
        if where == '' then
            where = _fun_ops(key, val)
        else
            where = where .. ' ' .. op .. ' ' .. _fun_ops(key, val)
        end
    end
    if where == '' then
        return where;
    else
        return '(' .. where .. ')';
    end
end
--读取数据：res, err, errcode, sqlstate
local function _exec_bysql(self, sql)
    local res, err, errcode, sqlstate = self.conn:query(sql)
    self.sql = sql
    if err then
        _log(self, 500, err)
        ngx_log(ngx_err, 'mysql host:', self.conf.host, ' port:', tostring(self.conf.port), ' sql:', sql, ' err:', err)
    else
        _log(self, 200, err)
    end
    return res, err, errcode, sqlstate
end
--读取数据：res, err, errcode, sqlstate
local function _select(self, table, fields, args, order, offset, limit)
    local where = _fun_exp(args)
    local sql   = 'select ' .. fields .. ' from ' .. table
    if where then
        sql = sql .. ' where ' .. where
    end
    if order ~= '' then
        sql = sql .. ' order by ' .. order
    end
    if offset and limit then
        sql = sql .. ' limit ' .. offset .. ',' .. limit
    end

    return _exec_bysql(self, sql)
end

local function _delete(self, table, args)
    local where = _fun_exp(args)
    local sql   = 'delete from ' .. table
    if where then
        sql = sql .. ' where ' .. where
    end

    return _exec_bysql(self, sql)
end

local function _update(self, table, fields, args)
    local where = _fun_exp(args)
    local sql   = 'update ' .. table .. ' set '
    local info = {}
    for key, val in pairs(fields) do
        info[#info+1] = key .. '=' .. ngx.quote_sql_str(val)
    end
    sql = sql .. table_concat(info, ',')
    if where then
        sql = sql .. ' where ' .. where
    end

    return _exec_bysql(self, sql)
end

local function _minsert(self, table, fields)
    local sql = 'insert into '..table
    local values, filed = array.new_table(#fields, 0), array.new_table(8, 0)
    for idx, item in ipairs(fields) do
        local tmp = array.new_table(8, 0)
        if idx == 1 then
            for key, val in pairs(item) do
                filed[#filed+1] = key
                tmp[#tmp+1] = ngx.quote_sql_str(val)
            end
        else
            for _, key in ipairs(filed) do
                local val = item[key]
                tmp[#tmp+1] = ngx.quote_sql_str(val)
            end
        end
        values[#values+1] = '(' .. table_concat(tmp, ',') .. ')'
    end
    sql = sql .. '(' .. table_concat(filed,',') ..')' .. ' values ' .. table_concat(values,',')

    return _exec_bysql(self, sql)
end

local function _insert(self, table, fields)
    return _minsert(self, table, {fields})
end

local function _exec(self, command, ...)
    if command == 'select' then
        return _select(self, ...)
    elseif command == 'delete' then
        return _delete(self, ...)
    elseif command == 'update' then
        return _update(self, ...)
    elseif command == 'insert' then
        return _insert(self, ...)
    elseif command == 'minsert' then
        return _minsert(self, ...)    
    end
end

local function _keepalive(self, idle, size)
    idle, size = idle or 10000, size or 1000
    if not self['conn'] or self.is_closed then return nil end

    self.conn:set_keepalive(idle, size)
    _log(self, 200)

    return true
end

function _M:new()
    return setmetatable({
        timeout = 100,
        confKey = 'mysql',
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

    local start_time                    = now()
    local timeout                       = args['timeout'] or conf['timeout'] or self.timeout
    local mysql, err, errcode, sqlstate = _connect(conf.host, conf.port, conf.database, conf.user, conf.charset, conf.password, timeout)
    local mc_obj                        = {
        conf = conf,
        busKey = busKey,
        start_time = start_time,
        timeout = timeout,
        exec = _exec,
        exec_sql = _exec_bysql,
        keepalive = _keepalive,
        conn = mysql,
        sql = '',
    }
    if err then
        _log(mc_obj, 502, err)
        return nil, err
    end

    return mc_obj
end

return _M