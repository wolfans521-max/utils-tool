import asyncio
import multiprocessing
import os
import signal
import sys
from multiprocessing import Process

from . import subprocess as subprocess_module
from ..share import share as shareVar


class ProcessManager:
    def __init__(self):
        self.shareQueue = multiprocessing.Queue()

    def sigchld_handler(self, sig, frame):
        try:
            while True:
                if sig == signal.SIGINT:
                    # 只能接收到主进程的退出信号
                    sys.stderr.write("接收到中断信号，正在退出...\n")
                    os._exit(0)
                pid, exitcode = os.waitpid(-1, os.WNOHANG)
                if pid == 0:
                    break
                process = shareVar.getSubProcess(pid)
                if process:
                    shareVar.delSubProcessId(pid)
                    process.join()
        except ChildProcessError:
            pass

    async def exceptJoins(self):
        while True:
            plist = shareVar.getSubProcessList()
            pids = list(plist.keys())
            colist = []
            for pid in pids:
                pro = (
                    plist[pid]["process"]
                    if pid in plist and "process" in plist[pid]
                    else None
                )
                if pro and not pro.is_alive():
                    colist.append(asyncio.to_thread(pro.join))  # 回收资源
                    shareVar.delSubProcessId(pid)
            if len(colist) > 0:
                done, colist = await asyncio.wait(colist, timeout=2)
            await asyncio.sleep(1)

    async def join(self):
        while True:
            try:
                pid = self.shareQueue.get(timeout=1)
                process = shareVar.getSubProcess(pid)
                if process:
                    await asyncio.to_thread(process.join)
                    shareVar.delSubProcessId(pid)
            except Exception as e:
                await asyncio.sleep(1)

    def add_schedule(self, taskId, routeInfo, route):
        if routeInfo.get("taskReady"):
            shareVar.setTaskReady(taskId, route, False)
            current_num = shareVar.getSubProcessNum(taskId, route)
            for _ in range(current_num, routeInfo.get("max_pnum", 0)):
                subp = Process(
                    target=subprocess_module.subprocess().run,
                    name=route,
                    args=(taskId, route),
                )
                # 不建议设为守护进程，避免未清理资源
                subp.start()
                shareVar.addSubProcessId(taskId, route, subp.pid, subp)

    async def run(self, taskId, schedule=None):
        signal.signal(signal.SIGCHLD, self.sigchld_handler)
        signal.signal(signal.SIGTERM, self.sigchld_handler)
        signal.signal(signal.SIGINT, self.sigchld_handler)
        routes = shareVar.scheduleList.get(taskId, {}).get("route", {})
        while True:
            if schedule is None:
                for route, routeInfo in routes.items():
                    self.add_schedule(taskId, routeInfo, route)
            else:
                routeInfo = routes.get(schedule)
                if routeInfo:
                    self.add_schedule(taskId, routeInfo, schedule)
            await asyncio.sleep(1)
