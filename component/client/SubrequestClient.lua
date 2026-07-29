--sub request client
local now                               = ngx.now
local ngxLog                            = ngx.log
local ngxErr                            = ngx.ERR
local capture                           = ngx.location.capture
local getmetatable                      = getmetatable
local setmetatable                      = setmetatable
local ngx                               = ngx
local type                              = type

local parent                            = require('zues.component.client.ClientInterface')
local jsonHelper                        = require('zues.utils.JsonHelper')
local array                             = require('zues.utils.ArrayHelper')

local json_decode  = jsonHelper.decode
local new_table    = array.new_table

local _M = {_VERSION = '0.1'}
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
--[[
    option = {
        url = '', //自请求地址
        returnStr = true, //bool 是否返回原始值，不为true时，返回json_decode后的值
        -- 一下参数参考ngx.location.capture的option参数
        args = {} or '', //请求参数 table or string
        ctx = {}, //传递给子请求的ngx.ctx表
    }
 ]]
function _M:s_curl(option)
    local ngx_ctx = ngx.ctx

    -- 处理子请求的ctx信息。保证子请求的日志信息能被追加到父请求中。
    if not ngx_ctx.message then ngx_ctx.message = new_table(16, 0) end
    if not ngx_ctx.appendLog then ngx_ctx.appendLog = new_table(0, 8) end
    if not ngx_ctx.zues_tablepool then ngx_ctx.zues_tablepool = new_table(8, 0) end

    local opt = option.option or new_table(0, 4)
    local ctx = opt.ctx or new_table(0, 4)
    if not ctx.message then ctx.message = ngx_ctx.message end
    if not ctx.appendLog then ctx.appendLog = ngx_ctx.appendLog end
    if not ctx.zues_tablepool then ctx.zues_tablepool = ngx_ctx.zues_tablepool end
    opt.ctx = ctx
    -- 参数赋值
    opt.args    = option.args
    opt.method  = option.method
    opt.body   = option.body

    -- 处理urlKey信息，复用self:log()
    local url = option.url
    option.protocol = 'subrequest'
    option.urlKey = url
    local startTime = now()

    local res, err = capture(url, opt)
    -- 返回结果
    if err or 'table' ~= type(res) or 200 ~= res.status then
        err = err or 'unknow'
        self:log(option, '', now() - startTime, 502, err, 1)
        ngxLog(ngxErr, url, err)
        return nil, err
    end
    --日志记录
    self:log(option, res.body, now() - startTime, res.status, err, 1)

    if option.returnStr then return res.body end

    return json_decode(res.body)
end

return _M