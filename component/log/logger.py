import logging
import os
from typing import Dict

from concurrent_log_handler import ConcurrentTimedRotatingFileHandler

from ..singleton import Singleton


class TagFilter(logging.Filter):
    def __init__(self, name, tag):
        self.tag = tag
        super().__init__(name)

    def filter(self, record):
        return getattr(record, 'tag', None) == str(self.tag)


class MillisecondFormatter(logging.Formatter):
    def formatTime(self, record, datefmt=None):
        from datetime import datetime
        ct = datetime.fromtimestamp(record.created)
        if datefmt:
            s = ct.strftime(datefmt)
            return s[:-3]  # 保留毫秒（去掉微秒的后三位）
        else:
            return super().formatTime(record, datefmt)


class Logger(Singleton):
    _logger: Dict[str, logging.Logger] = {}
    _config = {}
    _logqueue = {}
    taskId = ""

    # 设置配置文件
    def setConfig(self, config: list):
        if isinstance(config, list):
            for item in config:
                tag = item.get("tag")
                if tag:
                    self._config[tag] = item

        self.init_logger()

    def create_log_dir(self, filename):
        log_dir = os.path.dirname(filename)
        task_file_path = os.path.join(log_dir, self.taskId)
        file_name = os.path.basename(filename)

        # 如果目录不存在，则创建
        if task_file_path and not os.path.exists(task_file_path):
            os.makedirs(task_file_path, exist_ok=True)
        return os.path.join(task_file_path, file_name)

    def init_logger(self):
        # init queue
        for tag, item in self._config.items():
            handler = ConcurrentTimedRotatingFileHandler(
                filename=self.create_log_dir(item['filename']),
                # maxBytes=item.get('maxBytes', 50 * 1024 * 1024),  # 默认 50MB
                when="H",  # Rotate every hour
                interval=1,
                encoding="utf-8",
                backupCount=item.get('backupCount', 5)
            )
            formatter = MillisecondFormatter(item['format'], datefmt='%Y-%m-%d %H:%M:%S.%f')
            handler.setFormatter(formatter)
            handler.addFilter(TagFilter('myfilter', tag))

            logger = logging.getLogger(tag)
            logger.setLevel(item['level'])
            logger.addHandler(handler)
            self._logger[tag] = logger

    def get_logger(self, tag):
        return self._logger.get(tag, None)

    def log(self, message: str = None, tag: str = None, level=logging.INFO):
        llog = self.get_logger(tag)
        if not llog:
            return

        # 支持字符串级别如 "INFO"
        if isinstance(level, str):
            level = logging.getLevelName(level.upper())

        llog.log(level, message, extra={'tag': tag})
