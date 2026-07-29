plugin_routes = {
    # "sub_logger": {
    #     "min_pnum": 1,
    #     "max_pnum": 1,
    #     "loopnum": 1000,
    #     "loopsleepms": 1000,
    #     "crontab": "* * * * * *",
    # },
}
plugin_schedules = {
    # "sub_logger": [
    #     sub_logger,
    # ],
}


class Command:
    taskId = None  # 任务唯一标记
    desc = None  # 任务描述，适用于异步通知

    def __init__(self):
        self.routeList = plugin_routes
        self.scheduleList = plugin_schedules
        self.SetRoute()
        self.SetSchedule()

    def SetRoute(self):
        return True

    def SetSchedule(self):
        return True

    def getTaskId(self):
        return self.taskId
