import re
import time
from datetime import datetime


def _parse_cron_number(str, min, max):
    res = {}
    str_arr = re.split(",", str)
    for item in str_arr:
        md = re.split("/", item)
        step = 1
        if len(md) == 2:
            step = int(md[1])
        m_arr = re.split("-", md[0])
        if len(m_arr) == 2:
            _min = m_arr[0]
            _max = m_arr[1]
        else:
            if m_arr[0] == '*':
                _min = min
                _max = max
            else:
                _min = _max = m_arr[0]
        for pos in range(int(_min), int(_max) + 1):
            if pos % step == 0:
                res[pos] = pos
    return res


# return list [bool,list/string] bool：代表是否解析成功，成功话第二个参数为秒级list
def parse(cron_string, start_time=None):
    mm = re.match(
        '^((\*(\/[0-9]+)?)|[0-9\-\,\/]+)\s+((\*(\/[0-9]+)?)|[0-9\-\,\/]+)\s+((\*(\/[0-9]+)?)|[0-9\-\,\/]+)\s+((\*(\/[0-9]+)?)|[0-9\-\,\/]+)\s+((\*(\/[0-9]+)?)|[0-9\-\,\/]+)\s+((\*(\/[0-9]+)?)|[0-9\-\,\/]+)$',
        cron_string)
    if mm == None:
        mm = re.match(
            '^((\*(\/[0-9]+)?)|[0-9\-\,\/]+)\s+((\*(\/[0-9]+)?)|[0-9\-\,\/]+)\s+((\*(\/[0-9]+)?)|[0-9\-\,\/]+)\s+((\*(\/[0-9]+)?)|[0-9\-\,\/]+)\s+((\*(\/[0-9]+)?)|[0-9\-\,\/]+)$',
            cron_string)
    if mm == None:
        return False
    cron_arr = re.split('[\s]+', cron_string)
    date = {}
    if len(cron_arr) == 6:
        date = {
            'second': _parse_cron_number(cron_arr[0], 1, 59),
            'minutes': _parse_cron_number(cron_arr[1], 0, 59),
            'hours': _parse_cron_number(cron_arr[2], 0, 23),
            'day': _parse_cron_number(cron_arr[3], 1, 31),
            'month': _parse_cron_number(cron_arr[4], 1, 12),
            'week': _parse_cron_number(cron_arr[5], 0, 6),
        }
    elif len(cron_arr) == 5:
        date = {
            'second': {1: 1},
            'minutes': _parse_cron_number(cron_arr[0], 0, 59),
            'hours': _parse_cron_number(cron_arr[1], 0, 23),
            'day': _parse_cron_number(cron_arr[2], 1, 31),
            'month': _parse_cron_number(cron_arr[3], 1, 12),
            'week': _parse_cron_number(cron_arr[4], 0, 6),
        }
    else:
        return False, 'crontab任务解析错误'
    if not isinstance(start_time, int):
        if start_time != None:
            return False, 'start_time need int'
        start_time = int(time.time())
    now_date = datetime.fromtimestamp(start_time)
    # print(date)
    # print(now_date.second, now_date.minute, now_date.hour, now_date.day, now_date.month, now_date.weekday())
    if now_date.weekday() in date['week'] and now_date.month in date['month'] and now_date.day in date[
        'day'] and now_date.hour in date['hours'] and now_date.minute in date['minutes'] and now_date.second in date[
        'second']:
        return True, date['second']
    return True, {}
