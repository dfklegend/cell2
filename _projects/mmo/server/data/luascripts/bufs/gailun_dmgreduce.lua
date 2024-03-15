local this = BaseClass("GailunDmgReduceBuf")

local scriptName = "盖伦技能_免伤"
local function reset(self)
    self.idPreHit2 = -1
    self.buf = nil
    self.owner = nil
end

local function __ctor(self)
    reset(self)
end

local function onStart(self, buf)
    local owner = buf:Owner()
    local events = owner:GetEvents()

    self.buf = buf
    self.owner = owner

    self.idPreHit2 = events:SubscribeNoCheck("受到结算开始2", function(dmg)
        self:onPreHit2(dmg)
    end )
end

local function doEnd(self, buf)
    local owner = buf:Owner()
    local events = owner:GetEvents()

    if self.idPreHit2 > 0 then
        events:UnsubscribeById("受到结算开始2", self.idPreHit2)
        self.idPreHit2 = -1
    end
end

local function onEnd(self, buf)
    doEnd(self, buf)
    releaseBufScript(scriptName, self)
end

local function onPreHit2(self, dmg)
    -- 法术减伤
    --print('dmgType: '..dmg.DmgType)
    --print('TarBonusDmgReduceP: '..dmg.TarBonusDmgReduceP)
    dmg.TarBonusDmgReduceP = 0.5
    --print('TarBonusDmgReduceP: '..dmg.TarBonusDmgReduceP)
end

this.__ctor = __ctor
this.reset = reset
this.onStart = onStart
this.onEnd = onEnd
this.onPreHit2 = onPreHit2

registerBufScript(scriptName, this)
return this


