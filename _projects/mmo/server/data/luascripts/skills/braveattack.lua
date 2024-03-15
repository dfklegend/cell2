local this = BaseClass("BraveAttackSkill")
local apis = Root.goAPIs

--[[
    第3次英勇打击变成超级英勇打击
]]--

local function __ctor(self)
    --self.x = 0
end

local function changeSkill(self, owner)
    local userData = owner.UserData
    if userData.x == nil then
        return ''
    end
    if userData.x >= 3 then
        userData.x = 0
        --print('super英勇打击')
        return 'super英勇打击'
    end
end

local function onSkillHit(self, skill)   
    local userData = skill:Owner().UserData
    if userData.x == nil then
        userData.x = 0
    end
    userData.x = userData.x + 1
end

this.__ctor = __ctor
this.changeSkill = changeSkill
this.onSkillHit = onSkillHit

registerSkillScript("英勇打击", this)
return this


