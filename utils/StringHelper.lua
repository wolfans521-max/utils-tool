local str_find                      = string.find
local str_sub                       = string.sub
local str_gsub                      = string.gsub
local concat                        = table.concat
local reverse                       = string.reverse
local format                        = string.format
local match                         = ngx.re.match
local ngx_find                      = ngx.re.find
local upper                         = string.upper
local str_rep                       = string.rep
local decode_args                   = ngx.decode_args
local escape_uri                    = ngx.escape_uri
local unescape_uri                  = ngx.unescape_uri
local random                        = math.random
local tostring                      = tostring
local type                          = type
local tonumber                      = tonumber
local str_byte                      = string.byte

math.randomseed(ngx.now())

local _M = string

local ok, new_tab = pcall(require, "table.new")
if not ok or type(new_tab) ~= "function" then
    new_tab = function (narr, nrec) return {} end
end

-- 根据分隔符将字符串切割成数组
function _M.explode(d, p)
    local t, ll, index = {}, 1, 1
    if p == nil and d == nil then
        return {}
    elseif p == nil then
        return {}
    elseif d == nil then
        return {p}
    end
    if #p == 1 or #d < 1 then return {p} end
    while true do
        local startPos, endPos = str_find(p, d, ll, true)
        if startPos ~= nil then
            t[index] = str_sub(p,ll,startPos-1)
            ll = endPos + 1
            index = index + 1
        else
            t[index] = str_sub(p,ll)
            break
        end
    end

    return t
 end

-- 根据分隔符将字符串切割成数组
function _M.explode_old(delimiter, str)
    str = tostring(str)
    delimiter = tostring(delimiter)
    if 0 == #delimiter then return {str} end
    local pos, arr, index = 1, new_tab(0, 0), 1
    for st,sp in function() return str_find(str, delimiter, pos, true) end do
        local note = str_sub(str, pos, st - 1)
        if nil ~= note then
            arr[index] = note
            index = index + 1
        end
        pos = sp + 1
    end
    arr[index] = str_sub(str, pos)
    return arr
end

-- 使用分割符将数组中的元素连接成字符串
function _M.implode(delimiter, arr)
    delimiter = delimiter or ''
    return concat(arr, delimiter)
end

-- 暂不支持中文名称
local function basename(path)
    local name = str_gsub(path, "(.*/)(.*)", "%2")
    return name
end

-- 暂不支持中文名称
local function dirname(path)
    if path:match(".-/.-") then
        local name = str_gsub(path, "(.*/)(.*)", "%1")
        return name
    else
        return ''
    end
end

_M.basename = basename
_M.dirname = dirname

-- 去掉字符串两端的字符
function _M.trim_old(str, delimiter)
    if 'string' ~= type(str) then return '' end
    local patten = [[\s*(.*\S+)]]
    if nil ~= delimiter then
        patten = format([[(%s)*(.*[^%s*])]], delimiter, delimiter)
    end
    local res = match(str, patten) or {}

    return res[#res] or ''
end

function _M.trim(str, delimiter)
    if 'string' ~= type(str) then return '' end
    local patten = [[\s*(.*\S+)]]
    if nil ~= delimiter then
        patten = format([[(%s)*(.*[^%s*])]], delimiter, delimiter)
    else
        return str
    end
    local from, to, err = ngx_find(str, patten, "jo", nil, 2)
    if from then
        return str_sub(str, from, to)
    end

    return ''
end

-- 去除字符串左端匹配的字符
function _M.ltrim(str, delimiter)
    if 'string' ~= type(str) then return '' end
    local patten = [[\s*(.*)]]
    if nil ~= delimiter then
        patten = format([[(%s)*(.*)]], delimiter)
    end
    local res = match(str, patten) or {}

    return res[#res] or str
end

-- 去除字符串右端匹配的字符串
function _M.rtrim(str, delimiter)
    if 'string' ~= type(str) then return '' end
    local patten = [[(.*[^\s*])]]
    if nil ~= delimiter then
        patten = format([[(.*[^(%s)])]], delimiter)
    end
    local res = match(str, patten) or {}

    return res[#res] or str
end

-- 路径解析
function _M.pathinfo(path, key)
    if nil ~= key and 'dirname' ~= key
        and 'basename' ~= key and 'extension' ~= key
        and 'filename' ~= key then
        return nil
    end
    local ret = {filename = '',dirname = '', basename = '', extension = ''}
    -- 处理path
    if 'string' ~= type(path) or 0 == #path then
        return ret
    end

    local delimiter_pos, _ = str_find(path, '/', 1, true)
    local first_char = str_sub(path, 1, 1)
    if not delimiter_pos and '.' ~= first_char then path = './' .. path end
    -- dirname
    local dirname = dirname(path)
    if nil ~= dirname then
        ret.dirname = dirname
    end
    -- basename
    local basename = basename(path)
    if nil ~= basename then
        ret.basename = basename
    end
    -- extension and filename
    local reverse_path = reverse(basename)
    local start, _ = str_find(reverse_path, '.', 1, true)
    if nil == start then
        ret.filename = basename
    else
        local filename = str_sub(reverse_path, start + 1)
        local extension = str_sub(reverse_path, 1, start - 1)
        ret.filename = reverse(filename)
        ret.extension = reverse(extension)
    end

    if nil ~= key then
        return ret[key]
    end

    return ret
end

--url 解析
function _M.parse_url(url, module)
    local info = new_tab(0, 8)
    local res, err = match(url, [[^(?:(http[s]?):)?//([^:/\?]+)(?::(\d+))?([^\?]*)\??(.*)]], "jo")
    if err or type(res) ~= 'table' then
        res = {}
    end

    info['scheme'] = res[1] or 'http'
    info['host']   = res[2] or ''
    info['path']   = res[4] or ''
    info['query']  = res[5] or ''
    if tonumber(res[3]) then
        info['port'] = res[3]
    elseif info['scheme'] == 'https' then
        info['port'] = '443'
    elseif info['scheme'] == 'http' then
        info['port'] = '80'
    end

    if module and info[module] then return info[module] end

    return info
end

function _M.build_url(urlInfo)
    if 'table' ~= type(urlInfo) then return nil end
    local scheme = urlInfo['scheme'] or 'http'
    local ret = scheme..'://'
    --host
    if 'string' ~= type(urlInfo['host']) or #urlInfo['host'] <= 0 then
        return nil
    else
        ret = ret..urlInfo['host']
    end
    --port
    if tonumber(urlInfo['port']) then
        ret = ret..':'..urlInfo['port']
    end

    --path
    if 'string' == type(urlInfo['path']) and #urlInfo['path'] > 0 then
        ret = ret..urlInfo['path']
    else
        ret = ret..'/'
    end

    -- query
    if 'string' == type(urlInfo['query']) and #urlInfo['query'] > 0 then
        ret = ret..'?'..urlInfo['query']
    end

    return ret
end

-- 判断字符串是否以什么开头
function _M.start_with(str, with, case_sensitive)
    case_sensitive = case_sensitive or false
    if true == case_sensitive then
        return 1 == str_find(str, with, 1, true)
    end

    -- 不区分大小写
    local cmp_str = str_sub(str, 1, #with)
    return upper(cmp_str) == upper(with)
end

-- 判断字符以什么结尾
function _M.end_with(str, with, case_sensitiv)
    case_sensitiv = case_sensitiv or false
    str = reverse(str)
    with = reverse(with)

    if true == case_sensitiv then
        return 1 == str_find(str, with, 1, true)
    end

    -- 不区分大小写
    local cmp_str = str_sub(str, 1, #with)
    return upper(cmp_str) == upper(with)
end

--向右填充字符串
function _M.str_pad(str, pad_len,input)
    local str_len = #str

    if str_len < pad_len then
        input = input or ''
        str = str .. str_rep(input, pad_len - str_len)
    end

    return str
end

local function url_encode(s) return escape_uri(s) end
local function url_decode(s) return unescape_uri(s) end

_M.url_encode = url_encode
_M.url_decode = url_decode

function _M.split_param(str)
    str = url_decode(str or '')

    return decode_args(str)
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


local function chsize(char)
    if not char then
        return 0
    elseif char > 240 then
        return 4
    elseif char > 225 then
        return 3
    elseif char > 192 then
        return 2
    else
        return 1
    end
end

_M.chsize = chsize

-- 计算utf8字符串字符数, 各种字符都按一个字符计算
function _M.utf8len(str)
    local len = 0
    local currentIndex = 1
    while currentIndex <= #str do
        local char = str_byte(str, currentIndex)
        currentIndex = currentIndex + chsize(char)
        len = len +1
    end

    return len
end

-- 截取utf8 字符串
function _M.utf8sub(str, startChar, numChars)
    local startIndex = 1
    while startChar > 1 do
        local char = str_byte(str, startIndex)
        startIndex = startIndex + chsize(char)
        startChar = startChar - 1
    end

    local currentIndex = startIndex

    while numChars > 0 and currentIndex <= #str do
        local char = str_byte(str, currentIndex)
        currentIndex = currentIndex + chsize(char)
        numChars = numChars -1
    end

    return str:sub(startIndex, currentIndex - 1)
end

-- 按照字节长度截取utf8字符，多字节截取多出舍去保证不乱码
function _M.substr(str, numChars)
    local currentIndex = 1
    local currentLen   = 0
    while (currentLen < numChars) and currentIndex <= #str do
        local char = str_byte(str, currentIndex)
        local char_size = chsize(char)
        currentLen  = currentLen + char_size
        if currentLen > numChars then break end
        currentIndex = currentIndex + char_size
    end

    return str:sub(1, currentIndex - 1)
end

return _M
