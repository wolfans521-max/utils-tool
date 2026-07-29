import asyncio, time
from ..share import share as shareVar
from ..core import PareCrontab


class crontab:
    def __init__(self):
        self.last_runtime = 0

    def add_task(self, taskId, route, routeInfo, last_time):
        cronstatus = PareCrontab.parse(routeInfo["crontab"], last_time)
        if cronstatus[0] and len(cronstatus[1]) > 0:
            shareVar.setTaskReady(taskId, route, True)

    async def run(self, taskId, schedule=None):
        while True:
            routes = shareVar.scheduleList.get(taskId).get("route", {})
            now_time = int(time.time())
            if self.last_runtime == 0:
                self.last_runtime = now_time
            for i in range(self.last_runtime + 1, int(time.time()) + 1):
                self.last_runtime = i
                if schedule is None:
                    for route, routeInfo in routes.items():
                        self.add_task(taskId, route, routeInfo, i)
                else:
                    routeInfo = routes[schedule]
                    self.add_task(taskId, schedule, routeInfo, i)
            await asyncio.sleep(1)
