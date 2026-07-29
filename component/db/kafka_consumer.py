from kafka import KafkaConsumer

from ..decorator import check_args
from ..singleton import Singleton


class ZuesKafkaConsumer(Singleton):
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
    def getConnect(self, busKey) -> KafkaConsumer:
        busKey = str(busKey)
        if busKey not in self.__conns:
            config = self._init_config(busKey)
            self.__conns[busKey] = self.init_kafka_consumer(config)
        return self.__conns[busKey]

    def init_kafka_consumer(self, kafka_conf):
        bootstrap_servers = kafka_conf.get("kserver", "")
        consumer = KafkaConsumer(
            kafka_conf.get('ktopic', ''),
            group_id=kafka_conf.get('kgroup', ''),
            bootstrap_servers=bootstrap_servers.split(",") if isinstance(bootstrap_servers, str) else bootstrap_servers,
            auto_offset_reset="latest",
            enable_auto_commit=kafka_conf.get("enable_auto_commit", True),
            fetch_min_bytes=1024 * 16,
            fetch_max_wait_ms=kafka_conf.get("fetch_max_wait_ms", 500),
            max_partition_fetch_bytes=1024 * 1024 * 10,
            max_poll_records=kafka_conf.get("max_poll_records", 1000),
            sasl_plain_username=kafka_conf.get('sasl_plain_username', ''),
            sasl_plain_password=kafka_conf.get('sasl_plain_password', ''),
            security_protocol="SASL_PLAINTEXT",
            sasl_mechanism="PLAIN",
        )
        return consumer
