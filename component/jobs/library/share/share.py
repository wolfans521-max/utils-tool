scheduleList = {}
processscheduleMap = {}
config_path = ""


# 设置全局变量
def setSchedule(taskId, route, schedule):
    scheduleList[taskId] = {
        "taskId": taskId,
        "route": route,
        "schedule": schedule,
    }


# 获取任务
def getscheduleList(taskId):
    if taskId in scheduleList:
        return scheduleList[taskId]
    return {}


def taskExsting(taskId, schedule=None):
    scheduleListsInfo = getscheduleList(taskId)
    if schedule is None or not scheduleListsInfo:
        return len(scheduleListsInfo) != 0
    else:
        part_scheduleList = scheduleList[taskId]
        return schedule in part_scheduleList.get("route", {})


# 设置任务待执行
def setTaskReady(taskId, route, status):
    scheduleList[taskId]["route"][route]["taskReady"] = status


def addSubProcessId(taskId, route, pid, process):
    if pid not in processscheduleMap:
        processscheduleMap[pid] = {
            "taskId": taskId,
            "route": route,
            "process": process,
        }
        if "__sub_pids" not in scheduleList[taskId]["route"][route]:
            scheduleList[taskId]["route"][route]["__sub_pids"] = {}
        scheduleList[taskId]["route"][route]["__sub_pids"][pid] = 1


def getSubProcess(pid):
    if pid in processscheduleMap:
        return processscheduleMap[pid]["process"]
    return None


def getSubProcessList():
    return processscheduleMap


def delSubProcessId(pid):
    if pid in processscheduleMap:
        taskId = processscheduleMap[pid]["taskId"]
        route = processscheduleMap[pid]["route"]
        if (
                route in scheduleList[taskId]["route"]
                and pid in scheduleList[taskId]["route"][route]["__sub_pids"]
        ):
            scheduleList[taskId]["route"][route]["__sub_pids"].pop(pid)
        processscheduleMap.pop(pid)


def getSubProcessNum(taskId, route):
    if "__sub_pids" in scheduleList[taskId]["route"][route]:
        return len(scheduleList[taskId]["route"][route]["__sub_pids"])
    return 0


# 设置配置文件路径
def setConfigDir(confDir):
    global config_path
    config_path = confDir


# 获取配置文件路径
def getConfigDir():
    global config_path
    return config_path
