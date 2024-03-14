function OnOpen(node, tick)
    log.debug("LoginRq OnOpen , name = " .. tostring(node:GetName()) ..", title = ".. tostring(node:GetTitle()))
    tick.Blackboard:SetMem("login", true)
end

function OnTick(node, tick)
    log.debug("LoginRq OnTick , name = " .. tostring(node:GetName()) ..", title = ".. tostring(node:GetTitle()))
    return "success"
end
