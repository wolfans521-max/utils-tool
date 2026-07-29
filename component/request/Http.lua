local format                    = string.format
local ngx_now                   = ngx.now
local get_pid                   = ngx.worker.pid
local get_worker_id             = ngx.worker.id
local ngx                       = ngx
local setmetatable              = setmetatable
local type                      = type
local tonumber                  = tonumber

local arrayHelper               = require('zues.utils.ArrayHelper')
local valueHelper               = require('zues.utils.ValueHelper')
local numberHelper              = require('zues.utils.NumberHelper')
local cjson                     = require('zues.utils.JsonHelper')

local array_merge               = arrayHelper.merge
local simple_merge              = arrayHelper.simple_merge
local value_empty               = valueHelper.empty
local random                    = numberHelper.random
local get_value                 = arrayHelper.get_value
local json_decode               = cjson.decode

-- 纯Lua手动查找纯文本（彻底避开ngx.re.find的坑，兼容+号/二进制/大文件）
local function find_plain_text(str, substr, start_pos)
    start_pos = start_pos or 1
    local str_len = #str
    local substr_len = #substr
    
    -- 边界检查
    if substr_len == 0 or start_pos > str_len or str_len < substr_len then
        return nil
    end
    
    -- 逐字符对比（纯文本，+号就是普通字符）
    for i = start_pos, str_len - substr_len + 1 do
        local match = true
        for j = 1, substr_len do
            if str:sub(i + j - 1, i + j - 1) ~= substr:sub(j, j) then
                match = false
                break
            end
        end
        if match then
            return i, i + substr_len - 1  -- 返回和string.find一致的起止位置
        end
    end
    return nil
end
-- 解析Content-Type=multipart/form-data请求参数（上传图片版）
local function _parse_multipart_upload_args(content_type)
    -- 1. 读取请求体（兼容内存/临时文件）
    local body_data = ngx.req.get_body_data()
    if value_empty(body_data) then
        local body_file = ngx.req.get_body_file()
        if body_file then
            local file, err = io.open(body_file, "r")
            if file then
                body_data = file:read("*a")  -- 读取全部二进制内容
                file:close()
            else
                ngx.log(ngx.ERR, "Failed to open temp file: ", body_file, " error: ", err)
                return {}
            end
        else
            return {}
        end
    end

    -- 2. 提取并清理boundary
    local boundary = string.match(content_type, 'boundary=([^;]+)')
    if not boundary then
        return {}
    end
    -- 清理空格/引号（兼容各种格式）
    boundary = string.gsub(boundary, "%s+", "")
    boundary = string.gsub(boundary, "['\"]", "")
    local boundary_delimiter = "--" .. boundary  
    -- 3. 分割multipart parts（用纯Lua手动匹配，避开ngx.re.find坑）
    local parts = {}
    local current_pos = 1
    local body_len = #body_data
    local match_count = 0

    while true do
        -- 手动查找分隔符（纯文本，+号无问题）
        local start_idx, end_idx = find_plain_text(body_data, boundary_delimiter, current_pos)
        if not start_idx then
            break
        end
        match_count = match_count + 1

        -- 查找下一个分隔符，确定当前part的结束位置
        local next_start_idx = find_plain_text(body_data, boundary_delimiter, end_idx + 1)
        local part
        if next_start_idx then
            -- 提取当前part（不含分隔符）
            part = string.sub(body_data, end_idx + 1, next_start_idx - 1)
            current_pos = next_start_idx
        else
            -- 处理最后一个part（结束符：--Boundary+...--）
            local end_marker = boundary_delimiter .. "--"
            local end_marker_start = find_plain_text(body_data, end_marker, end_idx + 1)
            if end_marker_start then
                part = string.sub(body_data, end_idx + 1, end_marker_start - 1)
            end
            break
        end

        -- 过滤空part
        if part and part ~= "" then
            table.insert(parts, part)
        end
    end

    -- 4. 解析每个part的参数
    local params = {}
    for idx, part in ipairs(parts) do
        if part == "" or part == "--\r\n" then
            goto continue
        end

        -- 兼容\r\n\r\n（Windows）和\n\n（Linux）换行符
        local header_end = string.find(part, "\r\n\r\n") or string.find(part, "\n\n")
        if not header_end then
            goto continue
        end

        -- 提取header和body
        local header = string.sub(part, 1, header_end - 1)
        local offset = string.find(part, "\r\n\r\n") and 4 or 2  -- 动态偏移量
        local body = string.sub(part, header_end + offset)
        
        -- 清理body末尾的多余换行/空格
        body = string.gsub(body, "\r\n$", "")
        body = string.gsub(body, "\n$", "")
        body = string.gsub(body, "%s+$", "")

        -- 提取参数名（比如uid/file）
        local name = string.match(header, 'name="([^"]+)"')
        if name then
            name = string.gsub(name, "%s+", "")  -- 清理名称中的空格
            params[name] = body
        end
        ::continue::
    end

    return params
end

-- 解析Content-Type=multipart/form-data请求参数
local function _parse_multipart_args(content_type)
    local body_data = ngx.req.get_body_data()
    local boundary = string.match(content_type, 'boundary=(.+)')
    local parts = {}
    local start_idx, end_idx = 1, 1
    while true do
        start_idx, end_idx = string.find(body_data, "--" .. boundary, end_idx)
        if not start_idx then
            break
        end
        local part = string.sub(body_data, end_idx + 1)
        local part_end_idx = string.find(part, "--" .. boundary, 1)
        if part_end_idx then
            part = string.sub(part, 1, part_end_idx - 3)
        end
        table.insert(parts, part)
        end_idx = end_idx + 1
    end
    local params = {}
    for _, part in ipairs(parts) do
        if part ~= "" and part ~= "--\r\n" then
            local header_end = string.find(part, "\r\n\r\n")
            local header = string.sub(part, 1, header_end - 1)
            local body = string.sub(part, header_end + 4)
            local name = string.match(header, 'name="(.-)"')
            if name then
                name = string.gsub(name, "\r\n", "")
                params[name] = body
            end
        end
    end
    return params
end

local function _parse_args()
    local request = ngx.ctx.request
    local ngx_req = ngx.req
    local ngx_re  = ngx.re
    if 'table' == type(request) then return request end
    local uri_args, uri_truncated = ngx_req.get_uri_args(200)
    if uri_truncated == true then
        ngx.log(ngx.ERR, "uri args exceed 200 limit")
    end
    request = {
        getArgs = uri_args,
        headers = ngx_req.get_headers(),
        postArgs = {}, parsedArgs = {}
    }
    request.parsedArgs = request.getArgs
    if ngx_req.get_method() == 'POST' then
        ngx_req.read_body()
        local content_type = get_value(request, 'headers.content-type', '')
        local request_uri = ngx.var.document_uri
        if ngx_re.find(content_type, 'application/json') then
            request.postArgs = json_decode(ngx.req.get_body_data()) or {}
            if value_empty(request.postArgs) then
                local body_file = ngx_req.get_body_file()
                if body_file then
                    local file, err = io.open(body_file, "r")
                    if file then
                        local content = file:read("*a")
                        request.postArgs = json_decode(content) or {}
                        file:close()
                    else
                        ngx.log(ngx.ERR, "Failed to open temp file: ", body_file, " error: ", err)
                    end
                end
            end
        elseif content_type:match('multipart/form%-data') then 
            if request_uri and request_uri == '/search/picture/upload.json'then
                request.postArgs = _parse_multipart_upload_args(content_type)
            else
                request.postArgs = _parse_multipart_args(content_type)
            end
            
        else
            -- get_post_args() 只能读内存 buffer，body 超出 client_body_buffer_size 写入临时文件时会返回空 {}
            -- 此时需手动读取临时文件并用 ngx.decode_args() 解析
            local post_args, truncated = ngx_req.get_post_args(200)
            request.postArgs = post_args or {}
            if truncated == true then
                ngx.log(ngx.ERR, "post args exceed 200 limit")
            end
            if value_empty(request.postArgs, 'table') then
                local body_file = ngx_req.get_body_file()
                if body_file then
                    local file, err = io.open(body_file, "r")
                    if file then
                        local content = file:read("*a")
                        file:close()
                        request.postArgs = ngx.decode_args(content) or {}
                    else
                        ngx.log(ngx.ERR, "Failed to open post body temp file: ", body_file, " error: ", err)
                    end
                end
            end
        end

        request.parsedArgs = array_merge(request.getArgs, request.postArgs)
    end
    local gatewayParamsStr = get_value(request, 'headers.x-gateway-params', '')
    if gatewayParamsStr ~= '' then
        local gatewayParams = json_decode(gatewayParamsStr) or {}
        request.parsedArgs = array_merge(request.parsedArgs, gatewayParams)
    end
    ngx.ctx.request = request

    return request
end


local _M = {_VERSION='0.01'}
local mt = {__index = _M }


function _M:new()
    return setmetatable({
        seqid_param_key = 'seqid',
        seqid_header_key = 'seqid',
    }, mt)
end

-- 获取解析参数
function _M:get_param(key, default)
    local parsedArgs = _parse_args(self).parsedArgs
    if not key then return parsedArgs end
    local ret = parsedArgs[key] or default

    return ret
end

-- 优先获取get参数
function _M:get_post(key, default)
    local request = _parse_args(self)
    if not key then
        local ret = simple_merge(request.postArgs, request.getArgs, true)
        return ret
    end
    local ret = self:get(key)
    if ret then return ret end

    return self:post(key, default)
end

-- 优先获取post
function _M:post_get(key, default) return self:get_param(key, default) end

function _M:get_headers() return _parse_args(self).headers end

function _M:get_header(key, default)
    local headers = _parse_args(self).headers

    return headers[key] or default
end

--当前请求方法
function _M:get_method() return ngx.req.get_method() end

-- 判断是否是get请求
function _M:is_get() return 'GET' == ngx.req.get_method() end

--判断是否是post请求
function _M:is_post() return 'POST' == ngx.req.get_method() end

-- 获取get参数
function _M:get(key,default)
    local args = _parse_args(self).getArgs
    if key and args then
        return args[key] or default
    else
        return args
    end
end

--获取post请求参数
function _M:post(key,default)
    local args = _parse_args(self).postArgs
    if key and args then
        return args[key] or default
    else
        return args
    end
end

-- 获取seqid
function _M:seqid()
    --http header获取seqid信息
    if ngx.ctx['seqid'] then return ngx.ctx['seqid'] end
    local seqid = self:get_param(self.seqid_param_key)
    if value_empty(seqid,'string') then
        seqid = self:get_header(self.seqid_header_key)
    end

    if value_empty(seqid, 'string') then
        seqid = format("%14d%04d%d%d"
            , ngx_now()*10000, get_pid()
            , get_worker_id(), random(1,9))
    end

    if not seqid then seqid = '' end
    if seqid then ngx.ctx.seqid = seqid end

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