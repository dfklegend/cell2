

-- to skill
local function skill_onSkillHit(script, skill)
    --print(skill)
    script:onSkillHit(skill)
end

local function skill_test(userdata)
    
end


local apis = Root.Game.skillAPIs
apis.skill_onSkillHit = skill_onSkillHit


