

-- to skill
local function skill_onSkillHit(skill, num)
    print(skill)
    skill:onSkillHit(num)
end

local function skill_test(userdata)
   -- print(userdata)
    userdata.Str = ""
end


local apis = {}
Root.Game.skillAPIs = apis

apis.skill_onSkillHit = skill_onSkillHit
apis.skill_test = skill_test


