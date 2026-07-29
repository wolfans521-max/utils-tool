import asyncio
import random
from datetime import datetime
from multiprocessing import shared_memory, Manager

import fcntl
from aiofiles import open as async_open

# ffff_path = "/data1/apache2/config/search_media_community_uids.txt"
ffff_path = "/data1/apache2/config/num.log"


# multiprocessing.manager().dict() multiprocessing.shared_memory存储的结构
class DictShared:
    def __init__(self) -> None:
        self.shared_memory_cache = shared_memory.SharedMemory(
            name="test_share", create=True, size=10
        )
        self.buffer = self.shared_memory_cache.buf
        pass

    async def load(self, dconfig):
        num = random.random()
        shm_dict = Manager().dict()
        while True:
            now = datetime.now()
            formatted_now = now.strftime("%Y-%m-%d %H:%M:%S")
            print("whilw")
            with FileLock("/tmp/lock.txt", fcntl.LOCK_EX | fcntl.LOCK_NB) as filelock:
                if filelock:
                    print("wo sasdfasf---" + str(num) + "---" + formatted_now)
                    await self.loadmem(shm_dict, num)
                    await asyncio.sleep(1)
            await self.get(5156131388, "test")
            await asyncio.sleep(5)

    async def loadmem(self, filename, call_fun):
        count = 0
        # with open(ffff_path, "r") as file:
        #     for line in file:
        #         llist = line.strip().split("\t")
        #         first = llist.pop(0)
        #         dcache.set(str(first), json.dumps(llist), tag="test")
        #         count = count + 1
        #         if count % 100000 == 0:
        #             print(count)

        async with async_open(ffff_path, mode="r") as file:
            while True:
                line = await file.readline()
                if not line:  # 如果读取到的行为空，说明已经到达文件末尾
                    break
                llist = line.strip().split("\t")
                first = llist.pop(0)
                self.buffer[first] = bytearray(llist)

                # dcache.set(str(first), json.dumps(llist), tag="test")
                count = count + 1
                if count % 100000 == 0:
                    print(count)

    async def get(self, key, tag):
        return self.buffer[key]


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
            # print(f"无法加锁文件 {self.file_path}: {e}")
            self.file.close()
            self.file = None
        return self.file

    def __exit__(self, exc_type, exc_value, traceback):
        if self.file:
            try:
                fcntl.flock(self.file.fileno(), fcntl.LOCK_UN)
                print(f"解锁成功")
            except IOError as e:
                print(f"解锁文件时出错: {e}")
            finally:
                self.file.close()
