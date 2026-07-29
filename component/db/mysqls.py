#!/usr/bin/env python
# -*- coding:utf-8 -*-
from databases import Database


class MysqlPool:
    __conns = {}
    _config = {}

    def setConfig(self, config: dict):
        self._config = config

    def getConnect(self, busKey: str):
        if busKey in self.__conns:
            return self.__conns[busKey]
        cfg = self._config.get(busKey)
        if not cfg:
            raise ValueError(f"数据库配置未找到，key={busKey}")

        url = (
            f"mysql+asyncmy://{cfg['user']}:{cfg['password']}@"
            f"{cfg['host']}:{cfg['port']}/{cfg['db']}?charset={cfg.get('charset', 'utf8mb4')}"
        )
        db = Database(url)
        self.__conns[busKey] = db
        return db
