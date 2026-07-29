import base64
import hashlib
import hmac
import json
import os
import time
from urllib.parse import quote

wb_header = None
wb_header_time = 0
zhisou_wb_header = None
zhisou_wb_header_time = 0

header_dict = {}


class Auth:
    @staticmethod
    def get_token_by_discovery(auth_conf):
        from daemon.transfer_python.vendor.zues_python.component.discovery import configget
        discovery_key = auth_conf.get("discovery_key", "")
        project_key = auth_conf.get("project_key", "")
        assert auth_conf.get("style", "") == "discovery" and discovery_key and project_key
        token = configget(project_key, discovery_key)
        if token and isinstance(token, str):
            token = json.loads(token)
        return token

    @staticmethod
    def get_header(auth_conf, uid=''):
        if not uid:
            uid = auth_conf.get("uid", "")
        if not auth_conf.get("style", ""):
            tauth_file = auth_conf.get("tauth_file", "")
            # 判断fileName是否存在
            if not tauth_file or not os.path.exists(tauth_file):
                return None
            # 读取fileName内容，jsondecode，获取到密钥，计算签名
            with open(tauth_file, "r") as f:
                token = f.read(2048)
            token = json.loads(token)
        else:
            token = Auth.get_token_by_discovery(auth_conf)
        # 如果token是dict结构，则获取key
        if isinstance(token, dict) and "tauth_token" in token and "tauth_token_secret" in token:
            if uid:
                param = f"uid={uid}"
                sign = base64.b64encode(
                    hmac.new(
                        token["tauth_token_secret"].encode(),
                        param.encode(),
                        hashlib.sha1,
                    ).digest()
                ).decode()
                header = f'TAuth2 token="{quote(token["tauth_token"])}",param="{quote(param)}",sign="{quote(sign)}"'
            else:
                header = f'TAuth2 token="{token["tauth_token"]}"'
        else:
            header = None
        return header, uid

    @staticmethod
    def sign(auth_conf, auth_key='search'):
        # 超过5s 并且 wb_header不为空则重新获取
        global header_dict
        uid = auth_conf.get("uid", None)
        if uid is None:
            raise Exception("tauth config error")
        if header_dict.get(auth_key, {}):
            if time.time() - header_dict[auth_key]["time"] < 5 and not header_dict[auth_key]["header"] is None:
                return uid, header_dict[auth_key]["header"]
        header, uid = Auth.get_header(auth_conf)
        header_dict.setdefault(auth_key, )
        header_dict[auth_key] = {"time": time.time(), "header": header}
        return uid, header

    @staticmethod
    def userSign(auth_conf, tuid):
        header, uid = Auth.get_header(auth_conf, tuid)
        return uid, header
