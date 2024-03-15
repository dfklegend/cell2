--------------------------------------------------------------------------------
--      Copyright (c) 2015 - 2016 , 蒙占志(topameng) topameng@gmail.com
--      All rights reserved.
--      Use, modification and distribution are subject to the "MIT License"
--------------------------------------------------------------------------------
local setmetatable = setmetatable

local list = {}
list.__index = list

function list:new()
    local t = {length = 0, _prev = 0, _next = 0}
    t._prev = t
    t._next = t
    return setmetatable(t, list)
end

function list:clear()
    self._next = self
    self._prev = self
    self.length = 0
end
---添加一个数据
---@param item 数据
function list:push(item)
    --assert(value)
    local node = {value = item, _prev = 0, _next = 0, removed = false}

    self._prev._next = node
    node._next = self
    node._prev = self._prev
    self._prev = node

    self.length = self.length + 1
    return node
end

---添加一个节点
---@param node 节点
function list:pushnode(node)
    if not node.removed then
        return
    end

    self._prev._next = node
    node._next = self
    node._prev = self._prev
    self._prev = node
    node.removed = false
    self.length = self.length + 1
end

function list:pop()
    local _prev = self._prev
    self:remove(_prev)
    return _prev.value
end

function list:unshift(item)
    local node = {value = item, _prev = 0, _next = 0, removed = false}

    self._next._prev = node
    node._prev = self
    node._next = self._next
    self._next = node

    self.length = self.length + 1
    return node
end

function list:shift()
    local _next = self._next
    self:remove(_next)
    return _next.value
end

---删除
---@param node 节点
function list:remove(node)
    if node.removed then
        return
    end

    local _prev = node._prev
    local _next = node._next
    _next._prev = _prev
    _prev._next = _next

    self.length = math.max(0, self.length - 1)
    node.removed = true
end

function list:find(item, node)
    node = node or self

    repeat
        if item == node.value then
            return node
        else
            node = node._next
        end
    until node == self

    return nil
end

---查找返回第一个,如果没找到返回空
function list:firstordefault(condition)
    for v in ilist(self) do
        if condition(v.value) then
            return v.value
        end
    end
    return nil
end

function list:findlast(item, node)
    node = node or self

    repeat
        if item == node.value then
            return node
        end

        node = node._prev
    until node == self

    return nil
end

function list:next(node)
    local _next = node._next
    if _next ~= self then
        return _next, _next.value
    end

    return nil
end

function list:prev(node)
    local _prev = node._prev
    if _prev ~= self then
        return _prev, _prev.value
    end

    return nil
end

function list:erase(item)
    local node = self:find(item)

    if node then
        self:remove(node)
    end
end

function list:insert(item, node)
    if not node then
        return self:push(item)
    end

    local tempNode = {value = item, _next = 0, _prev = 0, removed = false}

    if node._next then
        node._next._prev = tempNode
        tempNode._next = node._next
    else
        self.last = tempNode
    end

    tempNode._prev = node
    node._next = tempNode
    self.length = self.length + 1
    return tempNode
end

function list:head()
    return self._next.value
end

function list:tail()
    return self._prev.value
end

function list:clone()
    local t = list:new()

    for i, v in list.next, self, self do
        t:push(v)
    end

    return t
end

---正序循环
ilist = function(_list)
    return list.next, _list, _list
end
---倒序循环
rilist = function(_list)
    return list.prev, _list, _list
end

setmetatable(list, {__call = list.new})
return list
