--http client
local format                            = string.format
local spawn                             = ngx.thread.spawn
local wait                              = ngx.thread.wait
local setmetatable                      = setmetatable
local ipairs                            = ipairs
local tonumber                          = tonumber
local type                              = type

local stringHelper                      = require('zues.utils.StringHelper')
local valueHelper                       = require("zues.utils.ValueHelper")
local arrayHelper                       = require("zues.utils.ArrayHelper")

local random       = stringHelper.random
local value_empty  = valueHelper.empty
local new_table    = arrayHelper.new_table
local in_array     = arrayHelper.in_array

local _M = {_VERSION='0.1', timeout = {50,50,300}}
local mt = {__index = _M}

function _M:new()
    return setmetatable({timeout = {100,100,500}}, mt)
end

-- da multi call
function _M:m_curl(options)
    if value_empty(options, 'table') then return {} end
    local thread, ret = new_table(4, 4), new_table(4, 4)
    for k, option in pairs(options) do
        thread[k] = spawn(self.s_curl, self, option)
    end
    for k, v in pairs(thread) do
        local _, res = wait(v)
        ret[k] = res
    end

    return ret
end

function _M:select_node(nodes, filterNodes)
    if value_empty(nodes, 'table') then return nil end
    filterNodes = filterNodes or new_table(#nodes, 0)
    local count = 0
    local effectiveNodes = new_table(#nodes, 0)
    for _, node in ipairs(nodes) do
        local domain = node.host
        if node.port then domain = domain..':'..node.port end
        node.domain = domain
        if not in_array(domain, filterNodes) then
            local weight = tonumber(node.weight) or 1
            count = count + weight
            effectiveNodes[#effectiveNodes + 1] = node
        end
    end
    if value_empty(effectiveNodes, 'table') then return nil end
    local index, host = 0
    local randomIndex = random(1, count)
    for _, node in ipairs(effectiveNodes) do
        local weight = tonumber(node.weight) or 1
        index = index + weight
        if randomIndex <= index then
            host = node
            filterNodes[#filterNodes + 1] = node.domain
            break
        end
    end

    return host
end

function _M:log(option, result, cost, statusCode, err, tryTimes, filterNodes)
    statusCode, err, tryTimes = statusCode or 200, err or 'success', tryTimes or 1
    if err ~= 'success' then statusCode = 502 end
    local len = 0
    if type(result) == 'string' then len = #result end
    local debug = {
        option = option,
        result = result,
        rt = format('%.3f', cost),
        status = statusCode,
        len = len,
        filterNodes = filterNodes,
        tryTimes = tryTimes,
        err = err
    }
    --  network:{protocol}|{urlKey}|{cost}|{nodeLog}|{len}|{tryTimes}|{statusCode}|{web_degrade}|msg:{err}
    local msg = format(' network:%s|%s|%.3f|%s|%d|%d|%d|%s|msg:%s',
        option.protocol,
        option.urlKey,
        cost,
        option.nodeLog or '-',
        len,
        tryTimes,
        statusCode,
        option.web_degrade or 0,
        err
    )
    local log = get_manager(self):get('log')
    log:append_log(msg, 'req', 'Info')
    log:log(debug, 'debug', 'Debug')
end

return _M