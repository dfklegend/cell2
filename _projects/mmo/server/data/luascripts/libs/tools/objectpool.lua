---@class ObjectPool : LuaObject
---@return ObjectPool
local this = BaseClass("ObjectPool")

---获得一个物体
---@param self ObjectPool
local function Get(self)
    return self._pool:Dequeue()
end

---回收一个物体
---@param self ObjectPool
local function Put(self, obj)
    self._pool:Enqueue(obj)
end

---构造函数
---@param self ObjectPool
local function __ctor(self)
    self._pool = Queue()
end

local function __dtor(self)
    self._pool = nil
end

this.__ctor = __ctor
this.__dtor = __dtor
this.Put = Put
this.Get = Get
return this
