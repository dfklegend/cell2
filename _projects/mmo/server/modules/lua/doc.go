/*

已知问题:
self如果自动按照luar的自动转化(any)，会自动转化为map
再传回脚本，回导致变成userdata而不是当前对象的问题
如下:
lua(Table)->goscript(any)->lua(userData)
(明确指定为*LTable，不会有问题)
lua(Table)->goscript(*LTable)->lua(Table)
可以用闭包来避免传回go
(事件中心，后面的是定义的...any)

	--self.idSkillHit = events:SubscribeNoCheck("onSkillHit", self.onSkillHit, self )
    self.idSkillHit = events:SubscribeNoCheck("onSkillHit", function(skill)
        self:onSkillHit(skill)
    end )

local function onSkillHit(self, skill)
    -- 每2次攻击+1 attack
    --print(self)
    --print(skill)
    self.tick = self.tick + 1
    if self.tick >= 2 then
        local owner = self.owner
        owner:OffsetAttack(1)
        self.tick = 0
    end
end
*/

package lua
