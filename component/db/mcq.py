import random
import time
import traceback

import aiodns
import aiomcache

from ..decorator import record_costtime
from ..hash import Hash
from ..singleton import Singleton


class McqDns(Singleton):
    __DNSResolver = None

    def __init__(self):
        if self.__DNSResolver is None:
            self.__DNSResolver = aiodns.DNSResolver()

    async def resolve(self, domain):
        ips = []
        try:
            result = await self.__DNSResolver.query(domain, "A")
            for record in result:
                ips.append(record.host)

        except Exception:
            traceback.print_exc()
        return ips


async def my_get_flag_handler(value: bytes, flags: int) -> str:
    if flags == 32:
        return value
    # 其他flag处理逻辑
    raise Exception(f"Unsupported flag {flags}")


class McqCli:

    def __init__(self, domain: str = None, port: int = 11233, get_handler=None):
        self.__domain = domain
        self.__port = port
        self.conns = {}
        self.ips = []
        self.__checktime = 0
        self.__check_gap = 30
        self.__get_cursor = None
        self.set_cursor = None
        self.get_handler = get_handler

    async def __resolve(self):
        now = time.time()
        if now - self.__checktime > self.__check_gap:
            self.__checktime = now
            ips = await McqDns().resolve(self.__domain)
            conn_ips = self.conns.keys()
            ips = sorted(ips)
            self.ips = ips
            new_ips = set(ips) - set(conn_ips)
            invalid_ips = set(conn_ips) - set(ips)
            len(new_ips) > 0 and self.__connect(new_ips)
            len(invalid_ips) > 0 and await self.__disconnect(invalid_ips)
            self.__get_cursor = (
                self.__get_cursor
                if self.__get_cursor is not None
                else random.randint(0, len(self.ips) - 1)
            )
            self.set_cursor = (
                self.set_cursor
                if self.set_cursor is not None
                else random.randint(0, len(self.ips) - 1)
            )

    def __connect(self, ips):
        for ip in ips:
            if self.get_handler == 32:
                self.conns[ip] = aiomcache.FlagClient(
                    ip, self.__port, pool_minsize=1, pool_size=10, get_flag_handler=my_get_flag_handler
                )
            else:
                self.conns[ip] = aiomcache.FlagClient(
                    ip, self.__port, pool_minsize=1, pool_size=10
                )

    async def __disconnect(self, ips):
        for ip in ips:
            await self.conns[ip].close()
            del self.conns[ip]

    @record_costtime
    async def get(self, key):
        trytimes = 0
        await self.__resolve()
        while trytimes < len(self.ips):
            self.__get_cursor = (self.__get_cursor + 1) % len(self.ips)
            trytimes += 1
            target_ip = self.ips[self.__get_cursor]
            if target_ip in self.conns:
                try:
                    value = await self.conns[target_ip].get(key.encode())
                    if value is not None:
                        return value
                except Exception:
                    traceback.print_exc()
        return None

    @record_costtime
    async def set(self, key: str = None, value: str = None, hash_key: str | int = None):
        await self.__resolve()
        if hash_key is None:
            self.set_cursor = (self.set_cursor + 1) % len(self.ips)
            if self.ips[self.set_cursor] in self.conns:
                return await self.conns[self.ips[self.set_cursor]].set(
                    key.encode(), value.encode()
                )
        else:
            hash_key = str(hash_key)
            hash_no = Hash.crc32(hash_key, len(self.ips))
            hash_ip = ""
            if len(self.ips) > hash_no:
                hash_ip = self.ips[hash_no]
            if hash_ip in self.conns:
                return await self.conns[self.ips[hash_no]].set(
                    key.encode(), value.encode()
                )
            else:
                raise AttributeError(f"mcq hash后对象{hash_no}有问题")
        return False


class Mcq(Singleton):
    conns = {}
    _config = {}

    # 检验配置信息是否存在
    def _init_config(self, busKey: str = None):
        config = {}
        if busKey in self._config:
            config = self._config[busKey]
        if not config:
            raise AttributeError(f"mcq对象实例化key，暂未发现配置信息")
        return config

    # 初始化参数
    def setConfig(self, config: dict = {}):
        self._config = config

    def getConnect(self, busKey) -> McqCli:
        busKey = str(busKey)
        if busKey not in self.conns:
            config = self._init_config(busKey)
            self.conns[busKey] = McqCli(domain=config["host"], port=config["port"],
                                        get_handler=config.get("get_handler", None))
        return self.conns[busKey]
