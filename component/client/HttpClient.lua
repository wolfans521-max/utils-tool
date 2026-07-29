--http client
local encodeArgs                = ngx.encode_args
local now                       = ngx.now
local ngxLog                    = ngx.log
local ngxErr                    = ngx.ERR
local type                      = type
local getmetatable              = getmetatable
local setmetatable              = setmetatable
local tonumber                  = tonumber
local table_insert              = table.insert
local table_concat              = table.concat

local http                      = require('resty.http')
local parent                    = require('zues.component.client.ClientInterface')
local stringHelper              = require('zues.utils.StringHelper')
local json                      = require('zues.utils.JsonHelper')
local arrayHelper               = require("zues.utils.ArrayHelper")
local zlib                      = require('resty.zlib')

local parse_url    = stringHelper.parse_url
local build_url    = stringHelper.build_url
local json_encode  = json.encode
local new_table    = arrayHelper.new_table

local function  _resolve_dsn(self, option, filterNodes)
    local nodes = option.nodes
    if type(nodes) ~= 'table' or #nodes < 1 then return end
    local node = self:select_node(nodes, filterNodes)
    if not node or not node.host then return end
    local urlInfo, nodeLog = parse_url(option.url), node.host
    if node.port then nodeLog = nodeLog .. '-' ..node.port end
    urlInfo.host = node.host
    urlInfo.port = node.port
    option.url = build_url(urlInfo)
    option.nodeLog = nodeLog

    return
end

-- post 请求
local function _post(self, option, httpc)
    local param = {method = 'POST' }
    -- 设置请求体
    if option['args'] then
        if option['contentJson'] then
            param['body'] = json_encode(option['args'])
        else
            param['body'] = encodeArgs(option['args'])
        end
    end
    local headers = option['header'] or {}
    --设置请求头
    if not headers['Content-Type'] then
        if option['contentJson'] then
            headers['Content-Type'] = 'application/json'
        else
            headers['Content-Type'] = 'application/x-www-form-urlencoded'
        end
    end
    param['headers'] = headers

    return httpc:request_uri(option['url'], param)
end

-- get 请求
local function _get(self, option, httpc)
    local param = {method = 'GET'}
    local url = option['url']
    -- 设置请求体
    if option['args'] then
        url = url .. '?' .. encodeArgs(option['args'])
    end
    --设置请求头
    if option['header'] then param['headers'] = option['header'] end

    return httpc:request_uri(url, param)
end

local function _decompress(body)
    if not body or #body <= 0 then return body end
    -- input function
    local page, body_len = 1, #body
    local function input(bufsize)
        local start = (page - 1) * bufsize + 1
        if start > body_len then return "" end
        local finish = start + bufsize - 1
        if finish > body_len then finish = body_len end
        page = page + 1
        local data = body:sub(start, finish)
        return data
    end
    -- output function
    local output_table = new_table(8, 0)
    local function output(data)
        table_insert(output_table, data)
    end
    local ok, err = zlib.inflateGzip(input, output)
    if not ok then
        ngxLog(ngxErr, "failed to decompress body, err:", err)
        return nil, "failed to decompress body"
    end
    local output_data = table_concat(output_table, '')

    return output_data
end

local _M = {_VERSION='0.1'}
local mt = {__index = _M}

function _M:new()
    local obj = parent:new()
    local super_mt = getmetatable(obj)
    -- 当方法在子类中查询不到时，再去父类中去查找。
    setmetatable(_M, super_mt)
    -- 这样设置后，可以通过self.super.method(self, ...) 调用父类的已被覆盖的方法。
    obj.super = setmetatable({}, super_mt)
    return setmetatable(obj, mt)
end

-- 发送请求
function _M:s_curl(option)
    local startTime, res, err = now()
    -- timeout
    local httpc = http:new()
    local timeout = option['timeout'] or self.timeout
    httpc:set_timeouts(timeout[1], timeout[2], timeout[3])

    local tryTimes, filterNodes = tonumber(option['tryTimes']) or 0, new_table(0, 0)
    tryTimes = tryTimes + 1
    local totalTimes = tryTimes
    -- 判断请求方法
    while tryTimes > 0 do
        _resolve_dsn(self, option, filterNodes)
        tryTimes = tryTimes - 1
        if 'POST' == option['method'] then
            res, err = _post(self, option, httpc)
        else
            res, err = _get(self, option, httpc)
        end
        if not err and res and res.status == 200 then
            break
        else
            local args = option['args'] or {}
            local no_errlog = type(args) == "table" and args.no_errlog == 1
            if not no_errlog then
                local res_status = res and res.status or 'nil'
                local res_type = type(res)
                local err_info = err or 'nil'
                ngxLog(ngxErr, "s_curl failed, url:", option['url'],
                       ", err:", err_info,
                       ", res_type:", res_type,
                       ", res_status:", res_status,
                       ", urlKey:", option['urlKey'] or 'nil',
                       ", tryTimes:", tryTimes,
                       ", timeout:", (timeout and table_concat(timeout, ',')) or 'nil')
            end
        end
    end

    -- 返回结果
    if err or 'table' ~= type(res) or 200 ~= res.status then
        err = err or 'unknow'
        self:log(option, '', now() - startTime, 502, err, totalTimes - tryTimes, filterNodes)
        return nil, err
    end
    -- gzip解压
    if res.headers and res.headers['Content-Encoding'] == 'gzip' and res.body then
       res.body, err =  _decompress(res.body)
       if err then res.status = 500 end
    end
    
    --日志记录
    self:log(option, res.body, now() - startTime, res.status, err, totalTimes - tryTimes, filterNodes)

    return res.body
end

return _M