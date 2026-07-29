from .client.Api import Api
from .config import Config
from .db.kafka_consumer import ZuesKafkaConsumer
from .db.kafka_product import ZuesKafkaProduct
from .db.mcq import Mcq
from .db.mysqls import MysqlPool
from .db.redis import Redis
from .db.mongo import Mongo
from .log.logger import Logger
from .singleton import Singleton


class appMap(Singleton):
    __map = {
        "redis": Redis,
        "mcq": Mcq,
        "api": Api,
        "mongo": Mongo,
        "kafka_consumer": ZuesKafkaConsumer,
        "kafka_product": ZuesKafkaProduct,
        "logger": Logger,
        "mysql": MysqlPool
    }

    def getInstanceMap(self):
        return self.__map

    def config(self, instance):
        if instance in self.__map:
            conf = Config().get(instance)
            if conf is None:
                raise AttributeError(f"{instance}配置不存在")
            return conf
        else:
            raise AttributeError(f"尚未实现该{instance}交互")
