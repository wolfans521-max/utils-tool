local format                    = string.format
local now                       = ngx.now
local ngx_log                   = ngx.log
local ngx_err                   = ngx.ERR
local type                      = type
local setmetatable              = setmetatable
local tostring                  = tostring
local ipairs                    = ipairs
local pcall                     = pcall

local redis                     = require('resty.redis')
local hash                      = require('resty.hash')
local confManager               = require('zues.component.Config')
local arrayHelper               = require('zues.utils.ArrayHelper')

local crc64_hash   = hash.get_crc64_hash
local new_table    = arrayHelper.new_table

local _M = {__VERSION = 0.1}

local mt = {__index = _M }

local function _log(self, statusCode, err)
    statusCode, err = statusCode or 200, err or 'success'
    if err ~= 'success' then statusCode = 502 end
    local cost = now() - self.start_time
    -- network:{protocol}|{urlKey}|{cost}|{nodeLog}|{len}|{tryTimes}|{statusCode}|{web_degrade}|msg:{err}
    local msg = format(' network:%s|%s|%.3f|%s|%d|%d|%d|0|msg:%s',
        'redis',
        self.busKey,
        cost,
        self.conf.port,
        0,
        1,
        statusCode,
        err
    )
    local log = get_manager():get('log')
    log:append_log(msg, 'req', 'Info')
    log:log({
        type = 'redis',
        conf = self.conf,
        db = self.db,
        request = self.request,
        busKey = self.busKey,
        timeout = self.timeout,
        cost = format('%.3f', cost),
        status = statusCode,
        err = err,
    }, 'debug','Debug')
end

local function _keepalive(self, idle, size)
    idle, size = idle or 10000, size or 1000
    if not self['conn'] or self.is_closed then return nil end

    self.conn:set_keepalive(idle, size)
    _log(self, 200)

    return true
end

local function _exec(self, command, ...)
    if not self['conn'] or self.is_closed then return nil end

    local ok, res, err = pcall(self.conn[command], self.conn, ...)

    self.request[#self.request + 1] = {
        command = command,
        args = {...},
        res = res,
        ok = ok,
    }

    if not ok or err then
        err = err or 'failed'
        _log(self, 504, err)
        self.conn:close()
        self.is_closed = true
        ngx_log(ngx_err, 'failed to exec redis cmd host:', self.conf.host, ' port:', tostring(self.conf.port), err)
    end

    return res
end

local function _gen_req(args)
    local nargs = #args

    local req = new_table(nargs * 5 + 1, 0)
    req[1] = "*" .. nargs .. "\r\n"
    local nbits = 2

    for i = 1, nargs do
        local arg = args[i]
        if type(arg) ~= "string" then
            arg = tostring(arg)
        end

        req[nbits] = "$"
        req[nbits + 1] = #arg
        req[nbits + 2] = "\r\n"
        req[nbits + 3] = arg
        req[nbits + 4] = "\r\n"

        nbits = nbits + 5
    end

    return req
end

local function send_data_mesh_hash(self, ...)
    if not self['conn'] or self.is_closed then return nil end
    local sock = self['conn']['_sock']
    if not sock then return nil end
    local args = {...}
    local req = _gen_req(args)

    sock:send(req)
end

local function _exec_hash(self, hashKey, command, ...)
    if not self['conn'] or self.is_closed then return nil end
    send_data_mesh_hash(self, 'hashkeyq', hashKey)
    local ok, res, err = pcall(self.conn[command], self.conn, ...)

    self.request[#self.request + 1] = {
        command = command,
        args = {...},
        res = res,
        ok = ok,
    }

    if not ok or err then
        err = err or 'failed'
        _log(self, 504, err)
        self.conn:close()
        self.is_closed = true
        ngx_log(ngx_err, 'failed to exec redis cmd host:', self.conf.host, ' port:', tostring(self.conf.port), err)
    end

    return res
end

local function _connect(host, port, timeout, db, auth)
    local red = redis:new()
    red:set_timeout(timeout)
    local ok, err = red:connect(host, port)
    if not ok then
        ngx_log(ngx_err, 'failed to connect redis host:', host, ' port:', tostring(port), err)
        return nil, err
    end
    -- 权限验证
    if auth then
        local ok, err = red:auth(auth)
        if not ok then
            ngx_log(ngx_err, 'failed to auth redis host:', host, ' port:', tostring(port), err)
            return nil, err
        end
    end
    red:select(db)

    return red
end

function _M:new()
    return setmetatable({
        timeout = 300,
        db = 0,
        confKey = 'redis',
        conf = 'resource',
    }, mt)
end

function _M:init()
    if 'string' == type(self.conf) then
        self.conf = confManager:load(self.conf)
    end

    self.conf = self.conf[self.confKey]
end

-- args = {
--  busKey = '', -- 对应host的key
--  hashKey = '',
--  hashNo = 0, -- 可指定的hash节点
--  timeout = 30,
--  db = 0,
--  hash_select = function, -- 可选，自定义节点选择函数，仅本次调用生效，不修改单例状态
-- }
function _M:get_connect(args)
    local busKey = args['busKey']
    if not busKey then return nil, 'arg busKey not found!' end

    local confs = self.conf[busKey]
    if 'table' ~= type(confs) then return nil, 'nodes not found!' end

    local hash_select = 'function' == type(args['hash_select']) and args['hash_select'] or nil

    local hashNo = 1
    if args['hashNo'] then
        hashNo = args['hashNo']
    elseif #confs > 1 then
        if not args['hashKey'] then return nil, 'arg hashKey not found!' end
        hashNo = self:hash(busKey, args['hashKey'], hash_select)
    end
    local conf = confs[hashNo]
    local timeout = args['timeout'] or conf['timeout'] or self.timeout
    local db = args['db'] or conf['db'] or self.db
    local start_time = now()
    local auth = conf['auth']
    local red, err = _connect(conf.host, conf.port, timeout, db, auth)
    local red_obj = {
        conf = conf,
        busKey = busKey,
        start_time = start_time,
        db = db,
        timeout = timeout,
        exec = _exec,
        exec_hash = _exec_hash,
        keepalive = _keepalive,
        conn = red,
        request = {},
    }
    if not red then
        _log(red_obj, 502, err)
        return nil, err
    end

    return red_obj
end

function _M:select_node(confs, hashKey, hash_select)
    -- 优先使用传入的局部 hash_select，其次使用单例上的 hash_select（兼容旧用法），最后使用默认 crc64
    local select_hook = hash_select or self.hash_select
    if 'function' == type(select_hook) then
        return select_hook(confs, hashKey)
    end

    return crc64_hash(hashKey, #confs) + 1
end

function _M:hash(busKey, hashKey, hash_select)
    local conf = self.conf[busKey]
    local hash_no = self:select_node(conf, hashKey, hash_select)

    return hash_no
end

function _M:hash_chunk(busKey, hashKeys)
    local confs = self.conf[busKey]
    local ret = new_table(0, #confs)
    for _, hashKey in ipairs(hashKeys) do
        local hash_no = self:select_node(confs, hashKey)
        local hash_no_str = tostring(hash_no)
        local hash_chunk = ret[hash_no_str] or new_table(16, 0)
        hash_chunk[#hash_chunk + 1] = hashKey
        if not ret[hash_no_str] then ret[hash_no_str] = hash_chunk end
    end

    return ret
end

function _M:format_zset(data, fileds)
    fileds = fileds or {'data', 'score'}
    if type(data) ~= 'table' or #data < 1 then return {} end
    local count = #data
    local ret = new_table(count/2, 0)
    local firstField, secondField = fileds[1], fileds[2]

    for i = 1, count, 2 do
        local temp = {}
        temp[firstField] = data[i]
        temp[secondField] = data[i+1]
        ret[#ret + 1] = temp
    end

    return ret
end

function _M:format_hash(data)
    if type(data) ~= 'table' or #data < 1 then return {} end

    local count = #data
    local ret = new_table(0, count/2)
    for i = 1, count, 2 do
        local key = data[i]
        local value = data[i+1]
        ret[key] = value
    end

    return ret
end

return _M