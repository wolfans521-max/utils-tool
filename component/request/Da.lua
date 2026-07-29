local ngx                   = ngx
local setmetatable          = setmetatable
local type                  = type
local tonumber              = tonumber

local da                    = require "resty.da"
local cjson                 = require "zues.utils.JsonHelper"

local json_decode           = cjson.decode


local _M = {_VERSION = '0.01'}
local mt = {__index = _M }

function _M:new()
    return setmetatable({
        seqid_param_key = 'seqid',
    }, mt)
end

function _M:parse_args(sock)
    local body, err = da:decode(sock)
    local request = json_decode(body)
    if not body or err or 'table' ~= type(request) then
        ngx.ctx.request = nil
        return nil, err
    end

    ngx.ctx.request = request

    return request
end

-- 获取解析参数
function _M:get_param(key, default)
    local request = ngx.ctx.request

    if not request then
        return nil
    end

    if not key then return nil end

    return request[key] or default
end

function _M:response(data, sock)
    if not data then
        return nil, "data is nil"
    end

    return da:send(data, sock)
end

-- 优先获取get参数
function _M:get_post(key, default) return self:get_param(key, default) end

-- 优先获取post
function _M:post_get(key, default) return self:get_param(key, default) end

function _M:get_headers() return {} end

function _M:get_header(key, default)
    return nil
end

--当前请求方法
function _M:get_method() return 'GET' end

-- 判断是否是get请求
function _M:is_get() return true end

--判断是否是post请求
function _M:is_post() return false end

-- 获取get参数
function _M:get(key,default)
    return self:get_param(key, default)
end

--获取post请求参数
function _M:post(key,default)
    return self:get_param(key, default)
end

-- 获取seqid
function _M:seqid()
    local seqid = self:get_param(self.seqid_param_key)
    if not seqid then seqid = '' end

    return seqid
end

-- 检查uid是否是一个正常登陆用户
function _M.check_uid_valid(uid)
    uid = tonumber(uid) or 0
    local ret = true

    if uid < 10000000 or uid > 9999999999 then
        ret = false
    end

    return ret
end

return _M