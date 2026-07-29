import pymongo
from pymongo import ReadPreference

from ..decorator import check_args
from ..singleton import Singleton


class Mongo(Singleton):
    __conns = {}
    _config = {}

    # 检验配置信息是否存在
    def _init_config(self, busKey: str = None):
        config = {}
        if busKey in self._config:
            config = self._config[busKey]
        if not config:
            raise AttributeError(f"Mongo对象实例化key，暂未发现配置信息")
        return config

    # 初始化参数
    def setConfig(self, config: dict = {}):
        self._config = config

    @check_args
    def getConnect(self, busKey) -> pymongo.MongoClient:
        busKey = str(busKey)
        if busKey not in self.__conns:
            config = self._init_config(busKey)
            if config.get("is_need_master"):
                self.__conns[busKey] = pymongo.MongoClient(config.get("uri"), read_preference=ReadPreference.PRIMARY)
            else:
                self.__conns[busKey] = pymongo.MongoClient(config.get("uri"))
        return self.__conns[busKey]
