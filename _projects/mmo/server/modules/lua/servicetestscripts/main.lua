print("main init")
-- 初始化公共属性
require('env')
require('core.init')
require('skills.init')

-- 初始化完毕

-- 启动函数
function start(goEnv)
    print('main.start: ')
    print(goEnv)
    goEnv.V = 100
    goEnv.F32 = 0.5
    goEnv.F64 = 0.5
    appEnv.goEnv = goEnv
    print(BaseClass)
end

function update()
end

function stop()
end


function testEnv(goEnv)
    goEnv.V = 88
    goEnv.F32 = 0.5
    goEnv.F64 = 0.5
end

function testInitEvents(events)
    events:SubscribeNoCheck("hello", function(data)
        data.V = 99
    end )
end

function testUserData(data)
    print(data)
    print(data.Data1)
    data.V = 6
    data.Data1.V = 7
end





