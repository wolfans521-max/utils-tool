import asyncio
import fcntl
import os
import time
from datetime import datetime

import aiofiles
import diskcache


# import open as async_open


class Dict_DiskCache:
    __dconfig = {}
    __disk_instance = {}  # diskcache实例
    __default_update_key = '__default_update_time__'  # 上次入cache的时间
    __default_file_mtime_key = '__default_file_mtime__'  # 上次更新的文件变更时间

    def __init__(self, dconfig) -> None:
        self.__dconfig = dconfig
        # 实例化各个diskcache
        for key, config in dconfig.items():
            disk_path = config['disk_path']
            self.__disk_instance[key] = diskcache.Cache(disk_path)
        pass

    async def load(self) -> None:
        while True:
            for key, config in self.__dconfig.items():
                file_lock = config['file_lock']
                now = datetime.now()
                with FileLock(file_lock, fcntl.LOCK_EX | fcntl.LOCK_NB) as filelock:
                    if filelock:
                        await self.__loadToCache(key, config)
            await asyncio.sleep(5)

    async def __loadToCache(self, key, config) -> None:
        file_path = config['file_path']
        ttl = config['ttl']
        parse_line_callbackfn = config['parse_line_callbackfn']
        stat_info = os.stat(file_path)
        file_mtime = stat_info.st_mtime
        # 没有导入 或者源文件变更
        last_file_mtime = await self.get(key, self.__default_file_mtime_key)
        if last_file_mtime is not None and file_mtime <= last_file_mtime:
            return
        # 追加更新时间
        self.__disk_instance[key].set(self.__default_file_mtime_key, time.time(), tag="")
        async with aiofiles.open(file_path, mode="r") as file:
            while True:
                line = await file.readline()
                if not line:  # 如果读取到的行为空，说明已经到达文件末尾
                    break
                l_status, l_key, l_value = parse_line_callbackfn(line)
                if l_status:
                    self.__disk_instance[key].set(str(l_key), l_value, tag="", expire=ttl)

    async def get(self, group, key):
        if group in self.__disk_instance:
            return self.__disk_instance[group].get(key)
        return None

    # 暂不允许写入数据
    async def set(self, ):
        pass


class FileLock:
    def __init__(self, file_path, lock_type=fcntl.LOCK_EX):
        self.file_path = file_path
        self.lock_type = lock_type
        self.file = None

    def __enter__(self):
        self.file = open(self.file_path, "w+")
        try:
            fcntl.flock(self.file.fileno(), self.lock_type)
        except IOError as e:
            self.file.close()
            self.file = None
        return self.file

    def __exit__(self, exc_type, exc_value, traceback):
        if self.file:
            try:
                fcntl.flock(self.file.fileno(), fcntl.LOCK_UN)
            except IOError as e:
                print(f"解锁文件时出错: {e}")
            finally:
                self.file.close()
