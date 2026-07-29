import time


class Singleton(object):
    _instance = None  # 类变量
    _config = None

    def __new__(cls, *args, **kwargs):
        if cls._instance is None:
            cls._instance = super(Singleton, cls).__new__(cls)
        return cls._instance


    def setConfig(self, config: dict = {}):
        self._config = config


def check_args(func):
    def wrapper(*args, **kwargs):
        if args[0]._config == {}:
            raise AttributeError(
                f"对象{args[0].__class__.__name__}配置信息未设置，请调用setConfig"
            )
        result = func(*args, **kwargs)
        return result

    return wrapper
