local this = BaseClass("RageBladeCoreBuf")

--[[
鬼索的狂暴之刃
每次普攻，攻速增加6%
    +10% 攻速
    +10% 法强
]]--

local scriptName = "鬼刀_机制"
local BufIdToAdd = '鬼刀_叠加'

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

    if self.idSkillHit > 0 then
        events:UnsubscribeById("onSkillHit", self.idSkillHit)
        self.idSkillHit = -1
    end
end

local function onEnd(self, buf)
    doEnd(self, buf)
    releaseBufScript(scriptName, self)
end

-- 如果需要，可以加入时间间隔
local function onSkillHit(self, skill)
    -- 加攻速buf
    if skill:IsBGSkill() then
        return
    end

    local owner = self.owner
    owner:AddBuf(BufIdToAdd, 1, 1)
end

this.__ctor = __ctor
this.reset = reset
this.onStart = onStart
this.onEnd = onEnd
this.onSkillHit = onSkillHit

registerBufScript(scriptName, this)
return this


