---@class Queue : LuaObject
local Queue = BaseClass("Queue")

local function setMeta(self)
    local mate = getmetatable(self)
    mate.__pairs = function(tb)
        return pairs(tb._datas)
    end
    mate.__ipairs = function(tb)
        return ipairs(tb._datas)
    end
    mate.__len = function(tb)
        return tb.n
    end
end

function Queue:__ctor()
    self._datas = {}
    self.n = 0
    setMeta(self)
end

---构造函数
---@param classType 泛型 可为空
function Queue:_ctor(classType)
    self._classType = classType
    self._datas = {}
    self.n = 0
    setMeta(self)
end

---添加队列
function Queue:Enqueue(element)
    if element == nil then
        return 
    end
    
    if self._classType then
        if not IsType(element, self._classType) then
            Logger.LogRed("添加队列类型错误,类型为:" .. self._classType.__className)
            return 
        end
    end
    
    table.insert(self._datas, element)
    self.n = self.n + 1
end

---出队列,返回并删除,如果队列为空则返回nil
function Queue:Dequeue()
    if self.n < 1 then
        return nil
    end
    
    local temp = table.remove(self._datas, 1)
    self.n = self.n - 1
    return temp
end

---出队列,返回不删除,如果队列为空则返回nil
function Queue:Peek()
    if self.n < 1 then
        return nil
    end
    return self._datas[1]
end

---清空队列
function Queue:Clear()
    table.clear(self._datas)
    self.n = 0
    -- self.size_ = 0
    -- self.head = -1
    -- self.rear = -1
end

---队列是否为空
function Queue:IsEmpty()
    return self.n == 0
    -- if self:size() == 0 then
    --     return true
    -- end
    -- return false
end

---队列长度
function Queue:GetCount()
    return self.n
end

return Queue
