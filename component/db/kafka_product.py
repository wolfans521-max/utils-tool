from kafka import KafkaProducer

from ..decorator import check_args
from ..singleton import Singleton


class ZuesKafkaProduct(Singleton):
    __conns = {}
    _config = {}

    # 检验配置信息是否存在
    def _init_config(self, busKey: str = None):
        config = {}
        if busKey in self._config:
            config = self._config[busKey]
        if not config:
            raise AttributeError(f"kafka对象实例化key，暂未发现配置信息")
        return config

    # 初始化参数
    def setConfig(self, config: dict = {}):
        self._config = config

    @check_args
    def getConnect(self, busKey) -> KafkaProducer:
        busKey = str(busKey)
        if busKey not in self.__conns:
            config = self._init_config(busKey)
            self.__conns[busKey] = self.init_kafka_producer(config)
        return self.__conns[busKey]

    def init_kafka_producer(self, kafka_conf):
        bootstrap_servers = kafka_conf.get("kserver", "")
        producer = KafkaProducer(
            bootstrap_servers=bootstrap_servers.split(",") if isinstance(bootstrap_servers,str) else bootstrap_servers,
            value_serializer=lambda v: str(v).encode('utf-8'),  # 消息内容编码为 UTF-8 字节流
            compression_type=kafka_conf.get("compression_type", "gzip"),  # 启用gzip压缩
            # Kafka broker 地址
            sasl_plain_username=kafka_conf.get("sasl_plain_username", ""),
            sasl_plain_password=kafka_conf.get("sasl_plain_password", ""),
            security_protocol="SASL_PLAINTEXT",  # 如果是启用 SASL 认证，则需要设置协议
            sasl_mechanism="PLAIN",  # 设置认证机制
            max_request_size=kafka_conf.get("max_request_size", 10485760)  # 设置最大请求大小为 10MB (10485760 bytes)
        )
        return producer
