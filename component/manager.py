import sys
from typing import Union

from .appMap import appMap
from .client.Http import Http
from .config import Config
from .db.kafka_consumer import ZuesKafkaConsumer
from .db.kafka_product import ZuesKafkaProduct
from .db.mcq import Mcq
from .db.mysqls import MysqlPool
from .db.redis import Redis
from .db.mongo import Mongo
from .log.logger import Logger
from .singleton import Singleton


class Manager(Singleton):
    __pool = {}
    taskId = ""

    def setConfPath(self, dir):
        sys.path.append(dir)

    def __initLogger(self):
        self.getInstance("logger")

    def getInstance(
            self, instance: str = None, module_name="zues_conf"
    ) -> Union[Http, Redis, Mcq, Mongo, Config, ZuesKafkaConsumer, ZuesKafkaProduct, MysqlPool, Logger, None]:
        instance = instance.lower()
        if instance not in self.__pool:
            if instance == "config":
                self.__pool[instance] = Config(module_name=module_name)
                return self.__pool[instance]
            elif instance == "logger":
                instance_map = appMap().getInstanceMap()
                args = appMap().config(instance)
                self.__pool[instance] = instance_map[instance]()
                self.__pool[instance].taskId = self.taskId
                self.__pool[instance].setConfig(args)
            else:
                instance_map = appMap().getInstanceMap()
                args = appMap().config(instance)
                self.__pool[instance] = instance_map[instance]()
                self.__pool[instance].setConfig(args)
        if instance != "logger":
            self.__initLogger()
        return self.__pool[instance]

    def getConfig(self):
        return self.__pool["config"]
