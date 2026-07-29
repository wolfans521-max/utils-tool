import importlib
import os


def autoload(name):
    """
    尝试根据给定的名称动态加载模块或包。

    :param name: 要加载的模块或包的名称（不包括前缀路径）
    :return: 如果成功加载，返回模块对象；否则返回None
    """
    # 假设所有模块都位于当前目录的某个子目录（比如'modules'）中
    module_path = os.path.join(os.path.dirname(__file__), "modules", name + ".py")

    # 检查文件是否存在
    if os.path.exists(module_path):
        # 使用importlib动态加载模块
        module_spec = importlib.util.spec_from_file_location(name, module_path)
        if module_spec:
            module = importlib.util.module_from_spec(module_spec)
            module_spec.loader.exec_module(module)
            return module

    # 如果没有找到模块，可以返回None或者抛出异常
    return None
