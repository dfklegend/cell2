local this = BaseClass("ElectricBladeCoreBuf")

--[[
电刀
    3次攻击后，释放一次电击伤害
    +10% 攻速
    +10% 物强
]]--

local scriptName = "电刀_机制"
local triggleConditionTimes = 3

local function reset(self)
    self.idSkillHit = -1
    self.buf = nil
    self.owner = nil

    self.tick = 0
end

local function __ctor(self)
    reset(self)
end

local function onStart(self, buf)
    local owner = buf:Owner()
    local events = owner:GetEvents()

    self.buf = buf
    self.owner = owner
    self.idSkillHit = events:SubscribeNoCheck("onSkillHit", function(skill)
        self:onSkillHit(skill)
    end )
end

local function doEnd(self, buf)
    local owner = buf:Owner()
    local events = owner:GetEvents()

    --print('电刀.onEnd')
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
    --
    if skill:IsBGSkill() then
        --print('bgSkill:'..skill:GetId())
        return
    end
    self.tick = self.tick + 1
    if self.tick >= triggleConditionTimes then
        self.tick = 0

        local owner = self.owner
        --print('电刀伤害')
        owner:CallbackSkill('电刀伤害', 1, owner, owner:GetTarId())
    end
end

this.__ctor = __ctor
this.reset = reset
this.onStart = onStart
this.onEnd = onEnd
this.onSkillHit = onSkillHit

registerBufScript(scriptName, this)
return this


