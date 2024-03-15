-- 可以将对象挂在Root上，防止命名冲突
-- 保证第一个require
if Root ~= nil then
    return
end

Root = Root or {}
Root.Game = Root.Game or {}
Root.goAPIs = Root.goAPIs or {}