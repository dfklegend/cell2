if Root.bufMgr == nil then
    Root.bufMgr ={}
end

local Stat = {}
Stat.hit = 0
Stat.total = 0



function registerBufScript(name, cls)
    local mgr = Root.bufMgr
    mgr[name] = cls
    print('registerBufScript: '..name)
    print(cls)
end

local function reportStat()
    print("buf script pool total: "..Stat.total.." hit: "..Stat.hit/Stat.total)
end

-- call by goscript
function createBufScript(name)
    local mgr = Root.bufMgr
    --print(name)
    local cls = mgr[name]
    --print('createSkillScript: '..name..' called')
    --print(cls)
    if cls == nil then
        return nil
    end

    if Stat.total > 0 and Stat.total%10000 == 0 then
        reportStat()
    end

    Stat.total = Stat.total + 1
    local obj = getBufScriptFromPool(name)
    if obj ~= nil then
        --print("hit pool: "..name)
        Stat.hit = Stat.hit + 1
        return obj
    end
    --print("miss pool: "..name)

    local obj = cls()
    return obj
end

function getBufScriptFromPool(name)
    local pools = Root.pools
    return pools:Get(name)
end

function putBufScriptToPool(name, obj)
    local pools = Root.pools
    pools:Put(name, obj)
    --print("put to pool: "..name)
end

function releaseBufScript(name, obj)
    obj:reset()
    putBufScriptToPool(name, obj)
end





