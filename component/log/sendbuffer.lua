local setmetatable = setmetatable
local next = next


local ok, new_tab = pcall(require, "table.new")
if not ok then
    new_tab = function (narr, nrec) return {} end
end


local _M = {}
local mt = { __index = _M }

function _M.new(self, batch_num, batch_size, max_reuse)
    local sendbuffer = {
        topics = {},
        queue_num = 0,
        batch_num = batch_num,
        batch_size = batch_size,
        max_reuse = max_reuse or 10000,
    }
    return setmetatable(sendbuffer, mt)
end


function _M.add(self, topic, partition_id, msg)
    local topics = self.topics

    if not topics[topic] then
        topics[topic] = {}
    end

    if not topics[topic][partition_id] then
        topics[topic][partition_id] = {
            queue = new_tab(self.batch_num, 0),
            index = 0,
            used = 0,
            size = 0,
        }
    end

    local buffer = topics[topic][partition_id]
    local index = buffer.index
    local queue = buffer.queue

    if index == 0 then
        self.queue_num = self.queue_num + 1
    end

    queue[index + 1] = msg

    buffer.index = index + 1
    buffer.size = buffer.size + #msg

    if (buffer.size >= self.batch_size) or (buffer.index >= self.batch_num) then
        return true
    end
end


function _M.clear(self, topic, partition_id)
    local buffer = self.topics[topic][partition_id]
    buffer.index = 0
    buffer.size = 0
    buffer.used = buffer.used + 1

    if buffer.used >= self.max_reuse then
        buffer.queue = new_tab(self.batch_num, 0)
        buffer.used = 0
    end

    self.queue_num = self.queue_num - 1
end


function _M.done(self)
    return self.queue_num == 0
end


function _M.loop(self)
    local topics, t, p = self.topics

    return function ()
        if t then
            for partition_id, queue in next, topics[t], p do
                p = partition_id
                if queue.index > 0 then
                    return t, partition_id, queue
                end
            end
        end


        for topic, partitions in next, topics, t do
            t = topic
            p = nil
            for partition_id, queue in next, partitions, p do
                p = partition_id
                if queue.index > 0 then
                    return topic, partition_id, queue
                end
            end
        end

        return
    end
end


return _M
