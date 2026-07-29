local ngx   = ngx
local unify = require "resty.unify"

local _M    = {
    _VERSION = '0.0.1'
}

-- 话题词转oid
function _M.topic2oid(topic, isunify, busid)
    return _M.tid2oid(_M.topic2tid(topic, isunify), busid)
end

-- 话题词转tid
function _M.topic2tid(topic, isunify)
    local tid

    if not topic then return tid end 

    if isunify then
        tid = unify.topicUid(topic)
    else
        tid = ngx.md5(topic)
    end

    return tid
end

-- 话题tid转oid
function _M.tid2oid(tid, busid)
    if not tid then return tid end
    
    if not busid then
        busid = '1022:231522'
    end

    return busid .. tid
end

-- 话题词归一
function _M.unify(query)

    return unify.topicUname(query)
end

return _M