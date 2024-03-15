---class:对象池系统
---@class PoolSystem : BaseSystem
---@field private _pools table<any,ObjectPool>
---@return PoolSystem

local ObjectPool = GetClass("ObjectPool")
local this = BaseClass("PoolSystem")

---GetItem:获取一个物体
---@param self PoolSystem
---@param key string 
local function Get(self, key)
    if self._pools[key] then
        return self._pools[key]:Get()
    end
    return nil
end

---GetItem:获取一个物体
---@param self PoolSystem
---@param key string 
---@param item
local function Put(self, key, item)
    if not self._pools[key] then
        self._pools[key] = ObjectPool()
    end
    self._pools[key]:Put(item)
end

local function __ctor(self)
    self._pools = {}
end

local function __dtor(self)
    self._pools = nil
end

this.__ctor = __ctor
this.__dtor = __dtor
this.Get = Get
this.Put = Put

if Root.pools == nil then
    Root.pools = this()
end

return this
