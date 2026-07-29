import asyncio
import os

from setproctitle import setproctitle
from zues_python.component.manager import Manager

from ..share import share as shareVar


class subprocess:
    objList = {}

    def run(self, taskId, route):
        setproctitle("wp_" + route)
        # self.shareQueue = queue
        # self.shareQueue.put(multiprocessing.current_process().pid)
        Manager().setConfPath(shareVar.getConfigDir())
        asyncio.run(self.exec(taskId, route))

    async def exec(self, taskId, route):
        taskInfo = shareVar.getscheduleList(taskId)
        routeInfo = taskInfo["route"][route]
        schedules = taskInfo["schedule"][route]
        loopNum = int(routeInfo["loopnum"]) or 1
        loopsleepms = int(routeInfo["loopsleepms"]) or 1
        # 初始化类
        if len(self.objList) == 0:
            for schedule in schedules:
                self.objList[schedule] = schedule()
        while loopNum > 0:
            for schedule, instance in self.objList.items():
                try:
                    await instance.run()
                except Exception as e:
                    import traceback
                    traceback.print_exc()
                    print(f"ERROR: schedule:{schedule}, instance:{instance},{e}")
            loopNum -= 1
            await asyncio.sleep(loopsleepms / 1000)
        # 子进程结束
        os._exit(0)
