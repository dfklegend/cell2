local this = BaseClass("ExampleSkill")
local apis = Root.goAPIs

local function __ctor(self)
    self.x = 0
    print 'create ExampleSkill'
end

local function onSkillHit(self, num)
    self.x = self.x + num
    --print(self.x)
    if self.x >= 3 then
        --print('hello')
        --apis.castSkill()
        self.x = 0
    end
end

this.__ctor = __ctor
this.onSkillHit = onSkillHit

registerSkillScript("example", this)
return this


