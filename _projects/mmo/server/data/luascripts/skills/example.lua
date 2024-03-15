local this = BaseClass("ExampleSkill")
local apis = Root.goAPIs

local function __ctor(self)
    --self.x = 0
end

local function onSkillHit(self, skill)
    local userData = skill.Owner.UserData
    if userData.x == nil then
        userData.x = 0
    end
    userData.x = userData.x + 1
    if userData.x >= 3 then
        print('cast 电击')
        --skill.Owner
        apis.castSkill(skill.Owner, "电击", skill.Tar)
        userData.x = 0
    end
end

this.__ctor = __ctor
this.onSkillHit = onSkillHit

registerSkillScript("攻击", this)
return this


