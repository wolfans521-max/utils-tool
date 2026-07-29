local type                      = type
local tostring                  = tostring
local tonumber                  = tonumber
local ngx_null                  = ngx.null
local next                      = next
local pairs                     = pairs

local _M = {_VERSION = '0.1' }

-- 判断是否是数字
function _M.is_numeric(value)
    if 'number' == type(value) then
        return true
    elseif 'string' == type(value) then
        local newValue = tonumber(value)

        if newValue and value == tostring(newValue) then return true end
    end

    return false
end

function _M.is_int(value) return 'number' == type(value) end

-- 判断是否是数组
function _M.is_array(value) return 'table' == type(value) end

-- 判断是否是函数(可调用)
function _M.is_callable(value) return 'function' == type(value) end

-- 判断是否是字符串
function _M.is_string(value) return 'string' == type(value) end

-- 判断redis返回是否是null
function _M.is_null(retJson)
    if retJson and retJson ~= 'null' and retJson ~= ngx_null and retJson ~= '' then
        return false
    else
        return true
    end
end

-- 判断是否为空
local function empty(value, v_type)
    if 'table' == v_type then
        if 'table' ~= type(value) or not next(value) then
            return true
        end
    elseif 'number' == type(value) then
        if 'number' ~= type(value) or 0 == value then return true end
    elseif 'function' == v_type and 'function' ~= type(value) then
        return true
    elseif 'string' == v_type then
        if 'string' ~= type(value) or '' == value then return true end
    end


    if nil == value or ngx_null == value or '' == value or 0 == value then
        return true
    end

    -- 判断是否表
    if 'table' == type(value) and not next(value) then
        return true
    end

    return false
end
_M.empty = empty

_M.is_empty_table = function(tb) return empty(tb, 'table') end

local new_tab
do
    local ok
    ok, new_tab = pcall(require, "table.new")
    if not ok or type(new_tab) ~= "function" then
        new_tab = function (narr, nrec) return {} end
    end
end

local nkeys
do
    local ok
    ok, nkeys = pcall(require, "table.nkeys")
    if not ok or type(nkeys) ~= "function" then
        nkeys = function(tab)
            local count = 0
            for _, _ in pairs(tab) do
                count = count + 1
            end
            return count
        end
    end
end

local function copy(var)
    if 'table' ~= type(var) then return var end

    local len = #var
    local res = new_tab(#var, nkeys(var) - len)

    for k, v in pairs(var) do
        if 'table' == type(v) then
            res[k] = copy(v)
        else
            res[k] = v
        end
    end

    return res
end
_M.copy = copy

--- 浅拷贝：仅复制第一层，用于参数隔离
function _M.shallow_copy(t)
    if 'table' ~= type(t) then return t end
    local res = {}
    for k, v in pairs(t) do
        res[k] = v
    end
    return res
end

return _M