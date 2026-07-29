local now                           = ngx.now
local random                        = math.random
local format                        = string.format
local floor                         = math.floor
local tonumber                      = tonumber

math.randomseed(now())
local _M = {
    _VERSION = "0.0.1"
}

function _M.round(num, numDecimalPlaces)
    return tonumber(format("%." .. (numDecimalPlaces or 0) .. "f", num))
end

function _M.random(startNum, endNum)
    local num

    if not endNum then
        num =  random(startNum)
    else
        num = random(startNum, endNum)
    end

    return num
end

function _M.format_fans_number(num, prec)
    num = tonumber(num) or 0
    prec = tonumber(prec) or 1
    local formatNum
    if prec >= 0 then prec = prec end
    if num >= 100000000 then
        formatNum  = format('%.'..prec..'f亿', num / 100000000)
    elseif num >= 10000 then
        formatNum  = format('%.'..prec..'f万', num / 10000)
    else
        formatNum  = num
    end

    return formatNum
end

function _M.format_star(num, mode)
    mode = mode or 1
    num = tonumber(num) or 0
    num = num * 10
    local star = '[空星][空星][空星][空星][空星]'
    if mode == 1 then
        if num <= 0 then
            return star
        elseif num >= 50 then
            star = '[星星][星星][星星][星星][星星]'
        else
            star = ''
            local total, helf = floor(num / 10), num % 10
            for i = 1, total do
                star = star..'[星星]'
            end
            if helf > 0 then
                total = total + 1
                star = star..'[半星]'
            end
            if total < 5 then
                for i = 1, 5 - total do
                    star = star..'[空星]'
                end
            end
        end
    end

    return star
end

return _M