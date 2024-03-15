---@class Set : LuaObject
local Set = BaseClass("Set")

local function GetCount(self)
    return table.length(self._datas)
end

local function GetDataByIdx(self, idx)
    return self._datas[idx]
end

local function __index(self, key)
    if GetTypeName(key) == "number" then
        if key > self:GetCount() then
            Logger.LogRed("索引越界")
            return nil
        else
            return self._datas[key]
        end
    else
        return global__index(self, key)
    end
end

local function Add(self, item)
    if item == nil then
        print("需要添加的节点为Null")
    end
    
    table.insert(self._datas, item)
end

local function Find(self, predicate)
    if (predicate == nil or type(predicate) ~= "function") then
        error("predicate is invalid!", 2)
        return 
    end
    
    local count = self:GetCount()
    
    for i = 1, count do
        if predicate(self._datas[i]) then
            return self._datas[i]
        end
    end
    return nil
end

local function ForEach(self, action)
    if (action == nil or type(action) ~= "function") then
        error("action is invalid!", 2)
        return 
    end
    
    local count = self:GetCount()
    
    for i = 1, count do
        action(self._datas[i])
    end
end

local function IndexOf(self, item)
    local count = self:GetCount()
    
    for i = 1, count do
        if self._datas[i] == item then
            return i
        end
    end
    return -1
end

---查找全部
---@param predicate Func<item,bool>
local function FindAll(self, predicate)
    if (predicate == nil or type(predicate) ~= "function") then
        error("predicate is invalid!", 2)
        return 
    end
    
    local newSet = Set()
    local count = self:GetCount()
    
    for i = 1, count do
        if predicate(self._datas[i]) then
            newSet:Add(self._datas[i])
        end
    end
    return newSet
end

---排序
---@param sortFunc Func<a,b,bool>
local function Sort(self, sortFunc)
    if sortFunc then
        table.sort(self._datas, sortFunc)
    end
end

local function LastIndexOf(self, item)
    local count = self:GetCount()
    
    for i = count, 1, -1 do
        if self._datas[i] == item then
            return i
        end
    end
    return -1
end

local function Insert(self, index, item)
    table.insert(self._datas, index, item)
end

local function Remove(self, item)
    local idx = self:LastIndexOf(item)
    
    if (idx > 0) then
        table.remove(self._datas, idx)
    end
end

local function RemoveAt(self, index)
    table.remove(self._datas, index)
end

---翻转
local function Reverse(self)
    local length = self:GetCount()
    
    if length > 0 then
        local half = math.floor(length / 2)
        
        for i = 1, half do
            local temp = self[i]
            local tempLen = length - (i - 1)
            self[i] = self[tempLen]
            self[tempLen] = temp
        end
    end
end

---变成数组
---@return Table
local function ToArray(self)
    local datas = {}
    local count = self:GetCount()
    
    for i = 1, count do
        table.insert(datas, self._datas[i])
    end
    return datas
end

---清空
local function Clear(self)
    table.clear(self._datas)
end

local function __ctor(self)
    self._datas = {}
    local meta = getmetatable(self)
    meta.__pairs = function(tb)
        return pairs(tb._datas)
    end
    meta.__ipairs = function(tb)
        return ipairs(tb._datas)
    end
end

local function __dtor(self)
    self._datas = nil
end

Set.__ctor = __ctor
Set.__dtor = __dtor
Set.Add = Add
Set.Remove = Remove
Set.RemoveAt = RemoveAt
Set.Insert = Insert
Set.LastIndexOf = LastIndexOf
Set.IndexOf = IndexOf
Set.Find = Find
Set.FindAll = FindAll
Set.Sort = Sort
Set.ForEach = ForEach
Set.Clear = Clear
Set.Reverse = Reverse
Set.ToArray = ToArray
Set.GetCount = GetCount
Set.GetDataByIdx = GetDataByIdx
Set.__index = __index

---已过时,请使用ipairs 或者pairs 进行迭代
function iset(set)
    if set.__classType ~= ClassType.instance then
        Logger.LogException("只能迭代Set的实例")
    end
    
    local nCount = set:GetCount()
    local datas = set._datas
    local index = 0
    return function()
        --Logger.Log(string.format("call close func. %s", index))
        index = index + 1
        
        if index <= nCount then
            if nil ~= datas[index] then
                return index, datas[index]
            else
                while index <= nCount do
                    index = index + 1
                    
                    if nil ~= datas[index] then
                        return index, datas[index]
                    end
                end
            end
        end
    end
end
return Set
