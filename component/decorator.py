import time
from .log.logger import Logger

def check_args(func):
    def wrapper(*args, **kwargs):
        if args[0]._config == {}:
            raise AttributeError(f"对象{args[0].__class__.__name__}配置信息未设置，请调用setConfig")
        result = func(*args, **kwargs)
        return result

    return wrapper

def record_costtime(func):
    async def wrapper(*args, **kwargs):
        class_name = args[0].__class__.__name__.lower()
        function_name = func.__name__.lower()
        if class_name =='redis_cli':
            function_name =args[1].lower()
        st = time.perf_counter() *1000
        result = await func(*args, **kwargs)
        et = time.perf_counter() *1000
        cost_time = round(et - st)
        # Logger().log(f'请求耗时:{class_name}--{function_name}--{cost_time}--{st}--{et}','zues')
        return result

    return wrapper