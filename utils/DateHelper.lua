local format                        = string.format
local os_time                       = os.time
local abs                           = math.abs
local ceil                          = math.ceil
local date                          = os.date
local floor                         = math.floor
local tonumber                      = tonumber

local _M = {_VERSION = '0.1',}

function _M.human_format(time, mode)
    mode = mode or 1
    time = tonumber(time) or os_time()
    local now = os_time()
    local dur = abs(now - time);

    if mode == 1 then
        if dur <= 10 then
            time = format('%s秒', dur)
        elseif dur < 50 then
            local second = ceil(dur / 10) * 10
            if second <= 0 then second = 10 end

            time = format('%s秒前', second)
        elseif dur < 3600 then
            local minutes = ceil(dur / 60)
            if minutes <= 0 then minutes = 1 end

            time = format('%s分钟前', minutes)
        elseif date("%Y%m%d") == date("%Y%m%d", time) then
            time = date('今天%H:%M', time)
        elseif date("%Y") == date("%Y", time) then
            time = date('%m月%d日', time)
        else
            time = date('%Y-%m-%d', time)
        end
    elseif mode == 2 then
        if dur < 60 then --1分钟内容
            time = "刚刚"
        elseif dur < 3600 then --一个小时以内
            local min = ceil(dur / 60)

            time = format('%s分钟前', min)
        elseif dur >= 3600 and dur <= 82800 then
            local hour = ceil(dur / 3600)

            time  = format('%s小时前', hour)
        else
            local day = ceil(dur / 86400)
            time = format('%s天前',day)
        end
    else
        local today = os_time(date('Y-m-d 00:00:00'));
        local yesterday = today - 86400;
        if dur < 600 then
            time = "刚刚"
        elseif dur < 3600 then
            local min = floor(dur / 60)
            if min <= 0 then min = 1 end
            time = format('%s分钟前', min)
        elseif dur >= 3600 and time >= today then
            local hour = floor(dur / 3600)
            if hour <= 0 then hour = 1 end
            time = format('%s小时前', hour)
        elseif time < today and time >= yesterday then
            time = date('昨天%H:%M', time);
        elseif date("%Y") == date("%Y", time) then
            time = date('%m-%d', time)
        else
            time = date('%Y-%m-%d', time)
        end
    end

    return time
end

-- 计算当前时间离当天的0点还过了多少s
-- hour：计算时间时可以向前 或向后推几个小时
function _M.range_to_0(hour)
    if not hour then hour = 0 end

    local day_s = 24*60*60
    local now = os_time() - (hour*60*60)

    return day_s - (now%day_s)
end

-- 计算当前时间从1970-01-01 08:00:00至现在经过了多少天
-- hour：计算时间时可以向前 或向后推几个小时
function _M.current_days(hour)
    if not hour then hour = 0 end

    local day_s = 24*60*60
    local now = os_time() - (hour*60*60)

    return ceil(now/day_s)
end

return _M