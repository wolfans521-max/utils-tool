local sub_str                   = string.sub
local str_find                  = string.find
local concat                    = table.concat
local table_remove              = table.remove
local floor                     = math.floor
local tonumber                  = tonumber
local type                      = type
local pairs                     = pairs
local next                      = next
local tostring                  = tostring
local ipairs                    = ipairs
local ngx                       = ngx

local stringHelper              = require('zues.utils.StringHelper')
local valueHelper               = require('zues.utils.ValueHelper')
local numberHelper              = require('zues.utils.NumberHelper')

local random            = stringHelper.random
local value_empty       = valueHelper.empty
local value_copy        = valueHelper.copy
local is_empty_table    = valueHelper.is_empty_table
local explode           = stringHelper.explode

local ok, new_tab = pcall(require, "table.new")
if not ok or type(new_tab) ~= "function" then
    new_tab = function (narr, nrec) return {} end
end

local fetch = function(tag, narr, nrec, noclear) return new_tab(narr, nrec) end
local release = function(endtag, obj, noclear) return end

local ok, tablepool = pcall(require, "tablepool")
if ok and type(tablepool) == 'table' then
    fetch = function(tag, narr, nrec, noclear, add_ctx)
        local ngx_ctx = ngx.ctx
        local tab = tablepool.fetch(tag, narr, nrec)

        if add_ctx then
            local pools = ngx_ctx.zues_tablepool or new_tab(8, 0)
            pools[#pools + 1] = {tag = tag, table = tab, noclear = noclear }
            if not ngx_ctx.zues_tablepool then ngx_ctx.zues_tablepool = pools end
        end

        return tab
    end
    release = tablepool.release
end

local _M = {
    _VERSION = '0.1',
    fetch = fetch,
    release = release,
    new_table = new_tab,
}

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

_M.nkeys = nkeys

-- 清除tablepool 需要在log阶段调用
function _M.clear_tablepool()
    local ngx_ctx = ngx.ctx
    local pools = ngx_ctx.zues_tablepool
    if pools then
        for _, item in ipairs(pools) do
            release(item.tag, item.table, item.noclear)
        end

        ngx_ctx.zues_tablepool = nil
    end
end

-- 获取数组中某个key值,可以嵌套获取
local function get_value(array, key, default, copy)
    if 'table' ~= type(array) or not key then return default end

    -- 判断key是否点号分割
    local start, _ = str_find(key, '.', 1, true)
    -- 获取一纬key
    if not start then
        key = tonumber(key) or key
        local ret = array[key]
        if not ret then return default end

        if copy and 'table' == type(ret) then return value_copy(ret) end

        return ret
    end

    -- 获取多纬key
    local keys, ret = explode('.', key), array
    local keys_count = #keys
    for i = 1, keys_count - 1 do
        local key = tonumber(keys[i]) or keys[i]
        if 'table' ~= type(ret) or not ret[key] then return default end

        ret = ret[key]
    end

    local key = tonumber(keys[keys_count]) or keys[keys_count]
    if 'table' ~= type(ret) or not ret[key] then return default end
    ret = ret[key]

    if copy and 'table' == type(ret) then return value_copy(ret) end

    return ret
--    local pop_key = sub_str(key, 1, start - 1)
--    pop_key =
--    local ret = array[pop_key]
--    if not ret then return default end
--
--    return get_value(ret, sub_str(key, start + 1), default, copy)
end
_M.get_value = get_value

-- 数组中设置某个值
function _M.set_value(arr, path, value)
    local keys
    if 'table' == type(path) then keys = path end
    if 'string' == type(path) then keys = explode('.', path) end

    if nil == keys then
        arr = value
        return
    end
    local temp, key = arr
    while #keys > 1 do
        key = table_remove(keys, 1)
        key = tonumber(key) or key
        if nil == temp[key] then
            temp[key] = {}
        end

        if 'table' ~= type(temp[key]) then
            temp[key] = {temp[key]}
        end

        temp = temp[key]
    end

    key = table_remove(keys, 1)
    key = tonumber(key) or key
    temp[key] = value
end

-- 获取二维数组中的某一列
function _M.get_column(array, key, copy)
    if 'table' ~= type(array) or not key then return {} end

    local ret, index = new_tab(nkeys(array), 0), 1
    for _, item in pairs(array) do
        if 'table' ~= type(item) then return {} end

        local value = item[key]
        if nil ~= value then
            if copy and 'table' == type(value) then value = value_copy(value) end
            ret[index] = value
            index = index + 1
        end
    end

    return ret
end

-- 数组合并, 多维数组会递归合并
local function merge(arr1, arr2, copy)
    if 'table' ~= type(arr1) and 'table' ~= type(arr2) then
        return {}
    elseif 'table' ~= type(arr1) then
        if copy then return value_copy(arr2) end

        return arr2
    elseif 'table' ~= type(arr2) then
        if copy then return value_copy(arr1) end

        return arr1
    end

    local ret = arr1
    if copy then ret = value_copy(arr1) end

    for k, v in pairs(arr2) do
        if copy and 'table' == type(v) then v = value_copy(v) end
        if 'number' == type(k) then
            if nil ~= ret[k] then -- 存在index 插入
                ret[#ret + 1] = v
            else -- 不存在，使用k作为index
                ret[k] = v
            end
        elseif('table' == type(v) and 'table' == type(ret[k])) then
            ret[k] = merge(ret[k], v, copy)
        else
            ret[k] = v
        end
    end

    return ret
end
_M.merge = merge

-- 数组简单合并
function _M.simple_merge(arr1, arr2, copy)
    if 'table' ~= type(arr1) and 'table' ~= type(arr2) then
        return {}
    elseif 'table' ~= type(arr1) then
        if copy then return value_copy(arr2) end

        return arr2
    elseif 'table' ~= type(arr2) then
        if copy then return value_copy(arr1) end

        return arr1
    end

    local ret = arr1
    if copy then ret = value_copy(arr1) end

    for k, v in pairs(arr2) do
        if copy and 'table' == type(v) then v = value_copy(v) end
        if 'number' == type(k) then
            if nil ~= ret[k] then -- 存在index 插入
                ret[#ret + 1] = v
            else -- 不存在，使用k作为index
                ret[k] = v
            end
        else
            ret[k] = v
        end
    end

    return ret
end
-- 根据给定的值，返回数组中对应的key，找不到时返回false
local function search(arr, value)
    if 'table' ~= type(arr) then return false end

    for k, v in pairs(arr) do
        if v == value then
            return k
        end
    end

    return false
end
_M.search = search

-- 获取数组中的所有key
function _M.keys(arr)
    if 'table' ~= type(arr) then return {} end
    local ret, index = new_tab(nkeys(arr), 0), 1
    for k, _ in pairs(arr) do
        ret[index] = k
        index = index + 1
    end

    return  ret
end

-- 获取表中的所有value值
function _M.values(arr, copy)
    if 'table' ~= type(arr) then return {} end

    local ret, index = new_tab(nkeys(arr), 0), 1
    for _, v in pairs(arr) do
        if copy and 'table' == type(v) then v = value_copy(v) end
        ret[index] = v
        index = index + 1
    end

    return ret
end

-- 键值互换
local function flip(arr, copy)
    if 'table' ~= type(arr) then return {} end
    local ret = new_tab(0, nkeys(arr))
    for k, v in pairs(arr) do
        if copy then
            if 'table' == type(k) then k = value_copy.copy(k) end
            if 'table' == type(v) then v = value_copy.copy(v) end
        end
        ret[v] = k
    end

    return ret
end
_M.flip = flip

-- 数组差集 存在于arr1中不存在于arr2中
function _M.diff(arr1, arr2)
    if 'table' ~= type(arr1) or 'table' ~= type(arr2) then return {} end
    local lc_arr2, arr1Cop = flip(value_copy(arr2)), value_copy(arr1)
    local ret, index = new_tab(0, 0), 1

    for _, v in pairs(arr1Cop) do
        if nil == lc_arr2[v] then -- 不存在于arr2
            ret[index] = v
            index = index + 1
        end
    end

    return ret
end

-- 数组过滤
function _M.filter(arr, call_back, copy)
    if 'table' ~= type(arr) then return {} end
    local ret = new_tab(0, 0)
    for k, v in pairs(arr) do
        if('function' == type(call_back) and true == call_back(k, v)) -- 自定义过滤函数
            or ('function' ~= type(call_back) and false == value_empty(v)) -- 默认判断值是否为空
        then
            if copy and 'table' == type(v) then v = value_copy(v) end
            if 'number' ~= type(k) then
                ret[k] = v
            else
                ret[#ret + 1] = v
            end
        end
    end

    return ret;
end

-- 数组切割
function _M.slice(arr, offset, length, copy)
    if 'table' ~= type(arr) or nil == next(arr) then return {} end
    offset  = tonumber(offset) or 1
    length  = tonumber(length) or 1
    local current_pos = 1
    local ret = new_tab(length, 0)
    for k, v in ipairs(arr) do
        if current_pos >= offset and current_pos < offset + length then
            if copy and 'table' == type(v) then v = value_copy(v) end
            if 'number' == type(k) then
                ret[#ret + 1] = v
            else
                ret[k] = v
            end
        end
        current_pos = current_pos + 1
    end

    return ret
end

-- 数组分块
function _M.chunk(arr, size, copy)
    local ret, tmp = new_tab(0, 0), new_tab(size, 0)
    if 'table' ~= type(arr) or nil == next(arr) then return ret end
    size = tonumber(size) or #arr
    local length = #arr

    if #arr <= size then
        return {arr}
    end

    local current_pos = 1
    for k, v in pairs(arr) do
        if copy and 'table' == type(v) then v = value_copy(v) end
        if 'number' == type(k) then
            tmp[#tmp + 1] = v
        else
            tmp[k] = v
        end

        if 0 == current_pos % size or current_pos == length then
            ret[#ret + 1] = tmp
            tmp = new_tab(size, 0)
        end

        current_pos = current_pos + 1
    end

    return ret
end

-- http_build_query
function _M.http_build_query(arr)
    local tmp = new_tab(0, 0)
    if 'table' ~= type(arr) then return '' end

    local index = 1
    for k, v in pairs(arr) do
        tmp[index] = tostring(k) .. '=' .. tostring(v)
        index = index + 1
    end

    return concat(tmp, '&')
end

-- 判断数据中是否存在某个值
function _M.in_array(item, arr)
    if 'table' ~= type(arr) then return false end

    for k, v in pairs(arr) do
        if v == item then
            return true
        end
    end

    return false
end

function _M.array_map(func, arr)
    if 'table' ~= type(arr) then return {} end

    for k, v in pairs(arr) do
        arr[k] = func(v)
    end

    return arr
end

-- 数组插入元素 效率比table.insert高
function _M.insert(arr, item)
    arr = arr or new_tab(4, 0)
    arr[#arr + 1] = item

    return arr
end

function _M.reverse(arr)
    if 'table' ~= type(arr) then return {} end
    local len = #arr
    local mid = floor(#arr / 2)
    for i = 1, mid, 1 do
        local tmp = arr[i]
        arr[i] = arr[len - i + 1]
        arr[len - i + 1] = tmp
    end

    return arr
end

-- 数组去重
function _M.uniq(arr, filed, copy)
    if value_empty(arr, 'table') then return {} end
    local ret = new_tab(#arr, 0)
    local list = new_tab(0, #arr)
    for _, item in ipairs(arr) do
        if filed and not value_empty(item, 'table') then
            local filedValue = item[filed] or 'nil'
            if not list[filedValue] then
                if copy then item = value_copy(item) end
                ret[#ret + 1] = item
            end
            list[filedValue] = true
        else
            if not list[item] then
                if copy and 'table' == type(item) then item = value_copy(item) end
                ret[#ret + 1] = item
            end
            list[item] = true
        end
    end

    return ret
end

function _M.shuffle(arr)
    if type(arr) ~= 'table' then return arr end

    local tmp, index
    for i = 1, #arr-1 do
        index = random(i, #arr)
        if i ~= index then
            arr[index], arr[i] = arr[i], arr[index]
        end
    end

    return arr
end

function _M.field2key(arr, field, copy)
    if 'table' ~= type(arr) then return {} end
    local ret = new_tab(0, #field)

    for _, v in ipairs(arr) do
        if copy then v = value_copy(v) end
        if not value_empty(v, 'table') then
            local key = v[field] or 'nil'
            ret[key] = v
        end
    end

    return ret
end

--  随机从数组取出一个（多个?）元素
function _M.random(arr)
    if is_empty_table(arr) then
        return {}
    end
    local rand = numberHelper.random(1, #arr)

    return arr[rand]
end

return _M
