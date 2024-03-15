local this = BaseClass("TitanCoreBuf")

--[[
泰坦的坚决
每次攻击和被攻击，增加2物攻和2%法术强度，最多叠加25次。叠加满之后，
增加25护甲和25法抗
        +10% 攻速
        +20 护甲
]]--

local scriptName = "泰坦_机制"
local BufIdToAdd = '泰坦_叠加'
local BufIdFull = '泰坦_叠加1'

local function reset(self)
    self.idSkillHit = -1
    self.idBeHit = -1
    self.buf = nil
    self.owner = nil

    self.addTimes = 0
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

    self.idBeHit = events:SubscribeNoCheck("onSkillBeHit", function(skill)
        self:onBeHit(skill)
    end )
end

local function doEnd(self, buf)
    local owner = buf:Owner()
    local events = owner:GetEvents()

    if self.idSkillHit > 0 then
        events:UnsubscribeById("onSkillHit", self.idSkillHit)
        self.idSkillHit = -1
    end

    if self.idBeHit > 0 then
        events:UnsubscribeById("onSkillBeHit", self.idBeHit)
        self.idBeHit = -1
    end
end

local function onEnd(self, buf)
    doEnd(self, buf)
    releaseBufScript(scriptName, self)
end

local function addOnce(self)
    if self.addTimes >= 25 then
        return
    end

    local owner = self.owner
    self.addTimes = self.addTimes + 1
    owner:AddBuf(BufIdToAdd, 1, 1)

    if self.addTimes == 25 then
        --print('泰坦，堆叠满')
        owner:AddBuf(BufIdFull, 1, 1)
    end
end

--
local function onSkillHit(self, skill)
    if skill:IsBGSkill() then
        return
    end
    addOnce(self)
end

local function onBeHit(self, skill)
    addOnce(self)
end

this.__ctor = __ctor
this.reset = reset
this.onStart = onStart
this.onEnd = onEnd
this.onSkillHit = onSkillHit
this.onBeHit = onBeHit

registerBufScript("泰坦_机制", this)
return this


