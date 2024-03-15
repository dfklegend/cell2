if Root.skillMgr == nil then
    Root.skillMgr ={}
end

function registerSkillScript(name, cls)
    local mgr = Root.skillMgr
    mgr[name] = cls
    print('register: '..name)
    print(cls)
end

-- call by goscript
function createSkillScript(name)
    local mgr = Root.skillMgr
    --print(name)
    local cls = mgr[name]
    --print('createSkillScript: '..name..' called')
    --print(cls)
    if cls == nil then
        return nil
    end
    local obj = cls()
    --print(obj.onSkillHit)
    return obj
end





