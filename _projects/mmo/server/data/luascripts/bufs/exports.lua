
--
local function buf_onAdd(script, buf)
    script:onAdd(buf)
end

local apis = Root.Game.skillAPIs
apis.buf_onAdd = buf_onAdd


