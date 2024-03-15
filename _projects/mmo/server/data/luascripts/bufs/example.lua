local this = BaseClass("ExampleBuf")

local scriptName = "buf1"
local function reset(self)
    self.idSkillHit = -1
    self.buf = nil
    self.owner = nil

    self.tick = 0
    self.attackAdded = 0
end

local function __ctor(self)
    reset(self)
end

local function onStart(self, buf)
    local owner = buf.Owner()
    local events = owner:GetEvents()

    self.buf = buf
    self.owner = owner
    --self.idSkillHit = events:SubscribeNoCheck("onSkillHit", self.onSkillHit, self )
    self.idSkillHit = events:SubscribeNoCheck("onSkillHit", function(skill)
        self:onSkillHit(skill)
    end )
end

local function doEnd(self, buf)
    local owner = buf.Owner
    local events = owner:GetEvents()

    -- 清除所有额外添加的属性
    local owner = self.owner
    owner:OffsetAttack(-self.attackAdded)

    if self.idSkillHit > 0 then
        events:UnsubscribeById("onSkillHit", self.idSkillHit)
        self.idSkillHit = -1
    end
end

local function onEnd(self, buf)
    doEnd(self, buf)
    releaseBufScript(scriptName, self)
end

local function onSkillHit(self, skill)
    -- 每2次攻击+1 attack
    --print(self)
    --print(skill)
    self.tick = self.tick + 1
    if self.tick >= 2 then
        local owner = self.owner

        local offset = 1
        owner:OffsetAttack(offset)
        self.attackAdded = self.attackAdded + offset
        self.tick = 0
    end
end

this.__ctor = __ctor
this.reset = reset
this.onStart = onStart
this.onEnd = onEnd
this.onSkillHit = onSkillHit

registerBufScript(scriptName, this)
return this


