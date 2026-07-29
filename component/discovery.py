#!/usr/bin/env python
# -*- coding:utf-8 -*-
import os
from cffi import FFI

base_dir = os.path.dirname(os.path.abspath(__file__))
ffi = FFI()

ffi.cdef(
    """
        bool configget(const char *group, const char *raw_key, char **result, bool verbose);
        bool configgetsign(const char *group, const char *raw_key, char **sign, bool verbose);
        bool namingget(const char *service, const char *cluster, char **result, bool verbose);
        bool naminggetsign(const char *service, const char *cluster, char **sign, bool verbose);
        size_t strlen(const char *s);
        void free(void *ptr);
    """
)

libdconf_path = os.path.join(os.path.dirname(os.path.abspath(__file__)), "libs", "libdconf.so")

# 修改为libdconf.so相对路径，该so文件位于discovery/ext/CPP/lib/so下
C = ffi.dlopen(libdconf_path)


def configget(group, raw_key):
    if len(group) == 0 or len(raw_key) == 0:
        return None
    result = ffi.new("char **")
    group_char = group.encode("utf-8")
    raw_key_char = raw_key.encode("utf-8")
    ret = C.configget(group_char, raw_key_char, result, 0)
    if not ret:
        return None
    res = ffi.string(result[0])
    ffi.gc(result[0], C.free)
    return res.decode("utf-8")


def configgetsign(group, raw_key):
    if len(group) == 0 or len(raw_key) == 0:
        return None
    result = ffi.new("char **")
    group_char = group.encode("utf-8")
    raw_key_char = raw_key.encode("utf-8")
    ret = C.configgetsign(group_char, raw_key_char, result, 0)
    if not ret:
        return None
    res = ffi.string(result[0])
    ffi.gc(result[0], C.free)
    return res.decode("utf-8")


def namingget(service, cluster):
    if len(service) == 0 or len(cluster) == 0:
        return None
    result = ffi.new("char **")
    se_char = service.encode("utf-8")
    cl_char = cluster.encode("utf-8")
    ret = C.namingget(se_char, cl_char, result, 0)
    if not ret:
        return None
    res = ffi.string(result[0])
    ffi.gc(result[0], C.free)
    return res.decode("utf-8")


def naminggetsign(service, cluster):
    if len(service) == 0 or len(cluster) == 0:
        return None
    result = ffi.new("char **")
    se_char = service.encode("utf-8")
    cl_char = cluster.encode("utf-8")
    ret = C.configgetsign(se_char, cl_char, result, 0)
    if not ret:
        return None
    res = ffi.string(result[0])
    ffi.gc(result[0], C.free)
    return res.decode("utf-8")


if __name__ == '__main__':
    print(configget("admin_ai", "comment_dispatch_is_open"))
