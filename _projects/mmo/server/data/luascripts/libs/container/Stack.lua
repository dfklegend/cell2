---@class Stack : LuaObject
local Stack = BaseClass("Stack")

function Stack:__ctor()
    self._datas = {}
    self._size = 0
    local mate =  getmetatable(self)
    mate.__pairs = function (tb)
        return pairs(tb._datas)
    end
    mate.__ipairs = function (tb)
        return ipairs(tb._datas)
    end
    mate.__len = function(tb)
        return tb._size
    end
end

function Stack:__dtor()
    self._datas = nil
end

function Stack:GetCount()
    return self._size
end

---压栈
function Stack:Push(element)
    if element == nil then
        return 
    end
    
    local size = self._size
    size = size + 1
    table.insert(self._datas, size, element)
    --self.Stack_table[size + 1] = element
    self._size = size
end

---出栈,如果没有则返回空
function Stack:Pop()
    local size = self._size
    
    if self._size == 0 then
        return nil
    end
    
    self._size = size - 1
    local res = self._datas[size]
    self._datas[size] = nil
    return res
end

---判断栈是否为空
function Stack:IsEmpty()
    if self._size == 0 then
        return true
    end
    
    return false
end

---栈转成table
function Stack:ToTable()
    local temp = {}
    
    for i, v in ipairs(self) do
        table.insert(temp, i, v)
    end
    return temp
end

---清空所有栈
function Stack:Clear()
    self._size = 0
    table.clear(self._datas)
end

---弹出第一个但不删除,如果没有返回nil
function Stack:Peek()
    --弹栈顶 不删除
    local size = self._size
    
    if size == 0 then
        --Logger.LogRed("Error:Stack is emtpy")
        return nil
    end
    
    return self._datas[size]
end

-- function istack(stack)
--     if stack.__classType ~= ClassType.instance then
--         Logger.LogException("只能迭代Stack的实例", 2)
--     end
    
--     local nCount = stack:GetCount()
--     local index = 0
--     return function()
--         --Logger.Log(string.format("call close func. %s", index))
--         index = index + 1
        
--         if index <= nCount then
--             if nil ~= stack._datas[index] then
--                 return index, stack._datas[index]
--             else
--                 while index <= nCount do
--                     index = index + 1
                    
--                     if nil ~= stack._datas[index] then
--                         return index, stack._datas[index]
--                     end
--                 end
--             end
--         end
--     end
-- end
return Stack
