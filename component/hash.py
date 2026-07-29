from ctypes import cdll
import os

crcFile = os.path.join(os.path.dirname(os.path.abspath(__file__)), "libs", "libhash.so")


class Hash:
    _libcrc = None

    @classmethod
    def crc32(cls, key, node_len):
        if cls._libcrc is None:
            cls._libcrc = cdll.LoadLibrary(crcFile)
        key = key.encode("utf-8")
        return cls._libcrc.hash_crc64(key, node_len)
