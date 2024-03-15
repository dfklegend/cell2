--[[
-- Lua面向对象设计
--]]
--已经定义过的类
local _classList = {}

---自定义类型
ClassType = {
    class = 1,
    instance = 2
}
---关键字
Keywords = {
    __call = 1,
    __index = 1,
    __newindex = 1,
    __origin = 1,
    __originName = 1,
    __delete = 1,
    __classType = 1,
    __className = 1,
    super = 1,
    --无参构造函数
    __ctor = 1,
    --有参构造函数
    _ctor = 1,
    --析构函数
    __dtor = 1,
    ---设置原方法
    __setMetaFunc = 1,
    ---
    __tostring = 1
}

---@class LuaAction @回调方法封装,可变参数处于最后一位
---@field public Obj table 当前事件的承载表
---@field public Func function 方法
---@field public Params ... 可变参数

---创建Lua事件
---@param action 方法
---@param tb table
---@param param 参数
function CreateLuaAction(action, tb, param)
    local luaAction = {}
    luaAction.Obj = tb
    luaAction.Func = action
    luaAction.Params = param
    ---@param self LuaAction
    luaAction.Call = function(self, ...)
        local status, err 
        
        if self.Params then
            --pcall(self.Func, self.Obj, ..., self.Params)
            if SafePack(...).n == 0 then
                -- self.Func(self.Obj, self.Params)
                status, err = pcall(self.Func, self.Obj, self.Params)
            else
                -- self.Func(self.Obj, ..., self.Params)
                status, err = pcall(self.Func, self.Obj, ..., self.Params)
            end
        else
            -- self.Func(self.Obj, ...)
            status, err = pcall(self.Func, self.Obj, ...)
        end
        
        if not status then
            Logger.LogRed(self.Obj)
            event_err_handle(err .. "\n" .. debug.traceback())
        end
        return status
    end
    return luaAction
end

---所有类的基类
---@class LuaObject
LuaObject = {
    super = nil,
    ---类名
    __className = "Object",
    __classType = ClassType.class,
    __init = false,
    __delete = false,
    ---构造函数
    __ctor = function(self)
        self._luaActions = {}
    end,
    ---析构函数
    __dtor = function(self)
        self._luaActions = nil
    end,
    -- ---设置原方法
    -- __setMetaFunc = function(self, funcName, func)
    --     if type(func) ~= "function" then
    --         error("传入的func参数错误", 2)
    --     end
    --     local meta = getmetatable(self)
    --     meta[funcName] = func
    -- end,
    ---获取类中的所有方法
    __getAllFunc = function(self)
        local meta = getmetatable(self)
        local str = "AllFunc:{\n"
        
        for key, value in pairs(meta.__functions.__lastFuncs) do
            str = str .. key .. "\n"
        end
        
        str = str .. "}"
        Logger.LogGreen(str)
    end,
    ---获取类中的所有字段
    __getAllField = function(self)
        local meta = getmetatable(self)
        local str = "AllField:{\n"
        
        for key, value in pairs(meta.__fields) do
            str = str .. key .. "\n"
        end
        
        str = str .. "}"
        Logger.LogGreen(str)
    end,
    __getAllSuper = function(self)
        local str = "Supper:{\n"
        local tmp = self
        
        while tmp do
            str = str .. "====>" .. tmp.__originName .. "\n"
            tmp = tmp.super
        end
        
        str = str .. "}"
        Logger.LogGreen(str)
    end,
    __index = nil,
    __newindex = nil,
    __tostring = function(self)
        return self.__className
    end,
    ---获取到当前类型
    ---@param self LuaObject
    GetType = function(self)
        local meta = getmetatable(self)
        return meta.class
    end,
    ---获取一个LuaAction
    ---@param self LuaObject
    ---@param action 方法
    ---@param param 其他参数<在最后一位>
    ---@return LuaAction
    GetLuaAction = function(self, action, param)
        if param == nil then
            local luaAction = self._luaActions[action]
            
            if luaAction == nil then
                luaAction = CreateLuaAction(action, self, param)
                self._luaActions[action] = luaAction
            end
            return luaAction
        else
            return CreateLuaAction(action, self, param)
        end
    end
}
_classList["LuaObject"] = LuaObject

---验证默认方法是否存在
local function checkDefaultFunc(class)
    local cls = class
    local nilFunc = function()
    end
    
    while cls ~= nil do
        --构造函数
        local tor = rawget(cls, "__ctor")
        
        if tor == nil then
            rawset(cls, "__ctor", nilFunc)
        elseif GetTypeName(tor) ~= "function" then
            rawset(cls, "__ctor", nilFunc)
        end
        
        --析构函数
        tor = rawget(cls, "__dtor")
        
        if tor == nil then
            rawset(cls, "__dtor", nilFunc)
        elseif GetTypeName(tor) ~= "function" then
            rawset(cls, "__dtor", nilFunc)
        end
        
        cls = cls.super
    end
end

---全局访问属性和方法(类)
function global__index(table, key)
    if key == "super" or key == "__originName" then
        return rawget(table, key)
    end
    
    local meta = getmetatable(table)
    
    local value = meta.__fields[key]
    
    if value ~= nil then
        return value
    end
    
    --最后检查是不是方法
    local originName = table.__originName
    --如果使用self调用
    if originName == meta.instance.__originName then
        value = meta.__functions.__lastFuncs[key]
        
        if value ~= nil then
            return value
        end
    else
        --使用base调用
        value = meta.__functions[originName][key]
        
        if value ~= nil then
            return value
        else
            error("父类:[" .. originName .. "]中没有方法:[" .. key .. "]")
        end
    end
    return value
end

---全局修改属性和方法(类)
function global__newindex(table, key, value)
    if Keywords[key] then
        error(key .. " 是关键字,不能修改 !!", 2)
    end
    
    local newType = GetTypeName(value)
    
    if newType ~= "function" then
        newType = "field"
    end
    
    local meta = getmetatable(table)
    -- local val = meta.__fields[key]
    
    -- if val ~= nil then
    --     --当新类型不是字段的时候
    --     if value and newType ~= "field" then
    --         error("正在给字段[" .. key .. "]赋值其他类型,请检查", 2)
    --         return 
    --     end
        
    --     meta.__fields[key] = value
    --     return 
    -- end
    
    -- val = meta.__functions.__lastFuncs[key]
    
    -- if val ~= nil then
    --     --当新类型不是字段的时候
    --     if value and newType ~= "function" then
    --         error("正在给方法[" .. key .. "]赋值其他类型,请检查", 2)
    --         return 
    --     end
        
    --     meta.__functions.__lastFuncs[key] = value
    --     return 
    -- end
    
    --当基础变量中不存在则新增变量
    if newType == "function" then
        meta.__functions.__lastFuncs[key] = value
    else
        meta.__fields[key] = value
    end
end

---设置一个类事件的元表
local function setEventMetaTable(meta)
    for key, value in pairs(meta.__events) do
        setmetatable(meta.__events[key], meta)
    end
end

---解析一个类
local function parseClass(class)
    local cls = class
    --查找类链
    local _clsListTemp = {}
    
    while (cls ~= nil) do
        table.insert(_clsListTemp, cls)
        cls = cls.super
    end
    
    -- -- 所有字段
    local fields = {} --只存最后一个<名称,Value>
    -- 所有方法
    local functions = {} -- Map<类名, Map<方法名,func> >
    
    local parentFuncs = {}
    --先从上往下依次调用类链中的所有方法.
    local _cls = nil
    
    for i = #_clsListTemp, 1, -1 do
        _cls = _clsListTemp[i]
        local className = _cls.__className
        functions[className] = {}
        
        for k, v in pairs(parentFuncs) do
            functions[className][k] = v
        end
        
        --找出方法,保存到方法列表中
        for k, v in pairs(_cls) do
            local typeName = type(v)
            
            if typeName == "function" then
                functions[className][k] = v
                
                if k ~= "_ctor" and k ~= "__dtor" and k ~= "__ctor" then
                    parentFuncs[k] = v
                end
            elseif typeName == "table" then
                local copy = DeepCopy(v)
                
                if Keywords[k] == nil then
                    fields[k] = copy
                end
            else
                if Keywords[k] == nil then
                    fields[k] = v
                end
            end
        end
    end
    
    --保存最后实现的那些方法. 实现多态.
    functions.__lastFuncs = functions[class.__className]
    --return fields, attributes, events, functions
    return fields, functions
end

--创建一个新类
local function newClass(class, ...)
    --obj是一个对象
    checkDefaultFunc(class)
    local obj = {}
    local meta = { instance = obj, class = class }
    meta.__mode = "k"
    obj.__originName = class.__className
    local fields, functions = parseClass(class)
    meta.__fields = fields
    meta.__functions = functions
    meta.__fields.__className = class.__className
    meta.__fields.__classType = ClassType.instance
    fields = nil
    functions = nil
    local createBase 
    createBase = function(o, c)
        if c.super then
            o.super = {
                __originName = c.super.__className --__origin = o
            }
            createBase(o.super, c.super)
        end
    end
    createBase(obj, class)
    
    meta.__index = class.__index and class.__index or global__index
    meta.__newindex = class.__newindex and class.__newindex or global__newindex
    meta.__tostring = class.__tostring
    -- meta.__index = global__index
    -- meta.__newindex = function(t, k, v)
    --     global__newindex(t, k, v)
    --     t.__newindex(t, k, v)
    -- end
    -- meta.__tostring = function(t)
    --     return t.__tostring(t)
    -- end
    
    setmetatable(obj, meta)
    local temp = obj.super
    
    while temp do
        setmetatable(temp, meta)
        temp = temp.super
    end
    
    -- 调用无参构造函数
    -- 无参构造函数是递归调用的
    local create 
    create = function(c)
        if c then
            if c.super then
                create(c.super)
            end
            
            c.__ctor(obj)
        end
    end
    create(obj.super)
    local len = select("#", ...)
    
    if len == 0 then
        obj.__ctor(obj)
    else
        obj._ctor(obj, ...)
    end
    
    --设置这个方便调试查看
    rawset(obj, "__meta", meta)
    obj.__init = true
    return obj
end

---基类,类似C#的object
---__ctor 为构造函数
---__dtor 为析构函数
---@param classname 类名
---@param super 父类表
---@class BaseClass
---@return BaseClass
function BaseClass(classname, super)
    assert(type(classname) == "string" and #classname > 0)
    assert(_classList[classname] == nil, '"' .. tostring(classname) .. '" 重复定义!!')
    local clsSource = {}
    --基类
    local baseClass = {}
    --当有父类
    if super then
        baseClass = _classList[super.__className]
        
        if baseClass == nil then
            error('"' .. tostring(super.__className) .. '" 父类不存在!!', 2)
        end
    else
        baseClass = LuaObject
    end
    
    setmetatable(
        clsSource,
        {
            __tostring = function(t)
                return t.__className
            end,
            __call = newClass,
            __index = function(t, k)
                local cls = t
                local val = nil
                
                while cls ~= nil do
                    val = rawget(cls, k)
                    
                    if val then
                        return val
                    end
                    
                    cls = rawget(cls, "super")
                end
                return val
            end
        }
    )
    
    clsSource.super = baseClass
    clsSource.__className = classname
    clsSource.__classType = ClassType.class
    clsSource.__init = false
    clsSource.__delete = false
    _classList[classname] = clsSource
    
    return clsSource
end

local function del(ins, c)
    c.__dtor(ins)
    
    if c.super then
        del(ins, c.super)
    end
end

---删除一个实例,调用完方法请执行 class = nil
---@param class 一个实例
function delete(class)
    if class == nil then
        return 
    end
    
    if type(class) ~= "table" then
        return 
    end
    
    local temp = rawget(class, "__delete")
    
    if temp then
        --已经删除过了
        return 
    end
    
    if (class.__classType ~= ClassType.instance) then
        error("要删除的不是一个实例")
    end
    
    rawset(class, "__delete", true)
    del(class, class)
    local meta = getmetatable(class)
    meta.__mode = "kv"
    meta.__fields = nil
    meta.__functions = nil
    meta.instance = nil
    meta.class = nil
    meta.__call = nil
    meta.__index = function(t, k)
        error("你正在访问一个删除的对象:" .. t.__originName .. "====" .. k)
    end
    meta.__newindex = function()
        error("你正在访问一个删除的对象")
    end
    class = nil
end
