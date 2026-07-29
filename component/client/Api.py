from typing import (
    Callable,
    Optional,
    Union,
)

from .Http import Http
from ..auth import Auth
from ..config import Config
from ..decorator import record_costtime
from ..singleton import Singleton


class Api(Singleton):
    __map = {"http": Http, "https": Http}

    def defineOptions(
            self,
            url: str = None,
            headers: Optional[list] = None,
            method: str = "GET",
            timeout: Optional[list] = None,
            args: Union[str, list] = None,
            tryTimes: int = 1,
            callbackFun: Callable = None,
    ):
        return {
            "url": url,
            "headers": headers,
            # 'method': method,
            "timeout": timeout,
            "args": args,
            "tryTimes": tryTimes,
            "callbackFun": callbackFun,
        }

    @record_costtime
    async def s_curl(self, optional: dict = None):
        urlKey = optional["urlKey"] or ""
        config = self._config[urlKey].copy() if urlKey in self._config else {}
        client = self.__parseProtocol(config["url"] or "")
        if client == None:
            return None
        client = client()
        # 签名信息
        if "auth" in config and config["auth"] == True and "user_auth" not in optional and "user_auth" not in config:
            if config.get("tauth"):
                auth_conf = config.get("tauth")
                auth_key = auth_conf.get("auth_key", None)
            else:
                auth_conf = Config().get("tauth")
                auth_key = "search"
            source, sign = Auth.sign(
                auth_key=auth_key if auth_key else "search",
                auth_conf=auth_conf,
            )
            if sign is not None:
                optional.setdefault("headers", {})["Authorization"] = sign
                optional["source"] = source
                del config["auth"]
        elif "auth" in config and config["auth"] and ("user_auth" in optional or "user_auth" in config):
            if config.get("tauth"):
                # 单独的配置文件
                auth_conf = config.get("tauth")
            else:
                # 全局认证配置文件
                auth_conf = Config().get("tauth")
            if optional.get("user_auth", {}).get("uid", ""):
                tuid = optional.get("user_auth", {}).get("uid", "")
            else:
                tuid = config.get("user_auth", {}).get("uid", "")
            source, sign = Auth.userSign(auth_conf, tuid)
            if sign is not None:
                optional.setdefault("headers", {})["Authorization"] = sign
                optional["source"] = source
        del optional["urlKey"]
        config = config | optional
        ret = await client.s_curl(config)
        return ret

    def __parseProtocol(self, url):
        if url == "":
            return None
        urls = url.split("://")
        protocol = urls[0] or ""
        if protocol in self.__map:
            return self.__map[protocol]
        return None
