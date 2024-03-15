---检测lua类实例是否是某个类型
---@param instance lua类实例
---@param ty 类型|类型名
---@return true|false
function IsType(instance, ty)
    local luaType = type(instance)
    
    if luaType ~= "table" then
        return false
    end
    
    if instance.__classType ~= nil and instance.__classType ~= ClassType.instance then
        Logger.LogError("不是一个lua类实体")
        return false
    end
    
    local typeName = ty.__className
    local cls = instance
    
    while cls do
        if cls.__className == typeName then
            return true
        end
        
        cls = cls.super
    end
    return false
end

---检测lua类的父子关系
---@param childType 子类型
---@param parentType 父类型
---@return true|false
function IsSubType(childType, parentType)
    if parentType then
        local cls = childType
        local parentName = parentType.__className
        
        while cls do
            if cls.__className == parentName then
                return true
            end
            
            cls = cls.super
        end
    end
    return false
end

local unpack = unpack or table.unpack

---解决原生pack的nil截断问题，SafePack与SafeUnpack要成对使用
function SafePack(...)
    local params = { ... }
    params.n = select("#", ...)
    return params
end

---解决原生unpack的nil截断问题，SafePack与SafeUnpack要成对使用
function SafeUnpack(safe_pack_tb)
    return unpack(safe_pack_tb, 1, safe_pack_tb.n)
end

-- --- 对两个SafePack的表执行连接
-- function ConcatSafePack(safe_pack_l, safe_pack_r)
--     local concat = {}
--     for i = 1, safe_pack_l.n do
--         concat[i] = safe_pack_l[i]
--     end
--     for i = 1, safe_pack_r.n do
--         concat[safe_pack_l.n + i] = safe_pack_r[i]
--     end
--     concat.n = safe_pack_l.n + safe_pack_r.n
--     return concat
-- end

-- 对两个SafePack的表执行连接
function ConcatSafePack(safe_pack_l, safe_pack_r)
    local concat = {}
    
    for i = 1, safe_pack_l.n do
        concat[i] = safe_pack_l[i]
    end
    
    for i = 1, safe_pack_r.n do
        concat[safe_pack_l.n + i] = safe_pack_r[i]
    end
    
    concat.n = safe_pack_l.n + safe_pack_r.n
    return concat
end

---获取一个对象的类型
---@param obj 一个对象
---@return 类型名
function GetTypeName(obj)
    if obj == nil then
        return "nil"
    end
    
    local _type = type(obj)
    
    if _type == "userdata" then
        local _meta = getmetatable(obj)
        
        if (_meta ~= nil and _meta.__name ~= nil) then
            local strs = string.split(_meta.__name, ".")
            return strs[#strs]
            --return _meta.__name
        end
        
        if obj.GetTypeName ~= nil then
            return obj:GetTypeName()
        end
    elseif _type == "table" then
        if obj ~= nil and obj.__className ~= nil then
            return obj.__className
        end
    end
    return _type
    -- local _type = type(obj)
    -- if (_type ~= "table" and _type ~= "userdata") then
    --     return _type
    -- elseif _type == "userdata" then
    --     local _meta = getmetatable(obj)
    --     if (_meta ~= nil and _meta.__name ~= nil) then
    --         local strs = string.split(_meta.__name, ".")
    --         return strs[#strs]
    --     --return _meta.__name
    --     end
    --     if obj.GetTypeName ~= nil then
    --         return obj:GetTypeName()
    --     end
    -- elseif (obj ~= {} and obj ~= nil and obj.__className ~= nil) then
    --     return obj.__className
    -- end
    -- return "Unknown"
end

---BeginSample:用于检测一个方法的执行效率
---@param num number 执行次数
---@param func function 需要检测的方法
---@param ... any 方法参数
function BeginSample(num, func, ...)
    local temp = os.clock()
    
    for i = 1, num, 1 do
        func(...)
    end
    
    print(os.clock() - temp)
end

---深拷贝不拷贝meta
function DeepCopy(originTable)
    local copy 
    
    if type(originTable) == "table" then
        --setmetatable(copy, DeepCopy(getmetatable(originTable)))
        copy = {}
        
        for orig_key, orig_value in next, originTable, nil do
            copy[DeepCopy(orig_key)] = DeepCopy(orig_value)
        end
    else
        copy = originTable
    end
    return copy
end

---深拷贝且拷贝meta
function DeepCopyWithMeta(originTable)
    local copy 
    
    if type(originTable) == "table" then
        copy = {}
        
        for orig_key, orig_value in next, originTable, nil do
            copy[DeepCopy(orig_key)] = DeepCopy(orig_value)
        end
        
        setmetatable(copy, DeepCopy(getmetatable(originTable)))
    else
        copy = originTable
    end
    return copy
end
