import importlib

from .singleton import Singleton


class Config(Singleton):
    _global_config = {}

    def __init__(self, module_name="zues_conf"):
        self.module_name = module_name
        self.__initConfig()

    def get(self, key: str):
        if key is None:
            return None
        return self._global_config.get(key, None)

    def __initConfig(self):
        if len(self._global_config) > 0:
            return
        self._global_config = {}
        st = self.__check_module_exists(self.module_name)
        if st:
            config = importlib.import_module(self.module_name)
            self._global_config = config.config
            return
        raise AttributeError(f"没有可调用配置文件,请完善{self.module_name}")

    def __check_module_exists(self, module_name):
        try:
            __import__(module_name)
            return True
        except ImportError:
            return False
