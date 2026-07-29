local ngx_now                           = ngx.now
local ngxLog                            = ngx.log
local ngxErr                            = ngx.ERR
local getmetatable                      = getmetatable
local setmetatable                      = setmetatable
local tonumber                          = tonumber

local jsonHelper                        = require('zues.utils.JsonHelper')
local parent                            = require('zues.component.client.ClientInterface')
local da                                = require('resty.da')
local arrayHelper                       = require("zues.utils.ArrayHelper")
local valueHelper                       = require('zues.utils.ValueHelper')

local json_encode   = jsonHelper.encode
local new_table     = arrayHelper.new_table


local _M = {
    _VERSION = "0.1.0"
}

local mt = { __index = _M }

function _M:new()
    local obj = parent:new()
    local super_mt = getmetatable(obj)
    -- 当方法在子类中查询不到时，再去父类中去查找。
    setmetatable(_M, super_mt)
    -- 这样设置后，可以通过self.super.method(self, ...) 调用父类的已被覆盖的方法。
    obj.super = setmetatable({}, super_mt)
    return setmetatable(obj, mt)
end

-- da single call
function _M:s_curl(opt)
    local startTime, res, err = ngx_now()
    -- 重定向node
    local ip    = arrayHelper.get_value(opt, 'args.ip', '')
    local port  = tonumber(arrayHelper.get_value(opt, 'args.port', 0)) or 0
    if not valueHelper.empty(ip) and port > 0 then 
        opt.nodes = {
            {
                domain = ip .. ':' .. port,
                host = ip,
                port = port,
            }
        }
    end
    local tryTimes = tonumber(opt['tryTimes']) or 0
    tryTimes = tryTimes + 1
    if tryTimes > 5 then tryTimes = 5 end
    -- 链接超时重试次数，最大3次
    local connect_timeout_times = 0
    local totalTimes, filterNodes = 0, new_table(#opt.nodes, 0)
    while tryTimes > 0 do
        tryTimes = tryTimes - 1
        totalTimes = totalTimes + 1
        local timeout_phase
        local node = self:select_node(opt.nodes, filterNodes)
        if not node then break end
        opt.node = node
        opt.nodeLog = node.host..'-'..node.port
        res, err, timeout_phase = self:handle_request(opt)
        if not err and res then
            break
        elseif timeout_phase == 1 and connect_timeout_times < 3 then
            tryTimes = tryTimes + 1
            connect_timeout_times = connect_timeout_times + 1
        else
            connect_timeout_times = 3
        end
    end
    self:log(opt, res, ngx_now() - startTime, nil, err, totalTimes, filterNodes)

    return res, err
end

function _M:handle_request(opt)
    local node       = opt.node
    local version = opt.version or 1
    if not node then
        local err = 'node-not-found'
        ngxLog(ngxErr, opt.urlKey, err)
        return nil, err
    end
    local _da = da:new()
    _da:set_version(version)
    local timeout = opt.timeout or self.timeout
    _da:set_timeouts(timeout[1],timeout[2],timeout[3])

    local args = json_encode(opt.args)
    return _da:request(node.host, node.port, args, opt.keepalive)
end

return _M
