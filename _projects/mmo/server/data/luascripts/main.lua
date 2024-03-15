print("main init")
-- 初始化公共属性
require('env')
require('libs.init')
require('skills.init')
require('bufs.init')

-- 初始化完毕

local function doTest()
    print('---- test ----')
    local apis = Root.goAPIs

    local self = {}
    apis.test1(function(t)
        print('test1')
        print(t)
    end, self)
    apis.test2(function(t)
        print('test2')
        print(t)
    end, self)
    apis.test3(function(t)
        print('test3')
        print(t)
    end, self)
end

-- 启动函数
function start(goEnv)
    print('main.start: ')
    --print(goEnv)
    appEnv.goEnv = goEnv
    --doTest()
end

function update()
end

function stop()
end

function globalhello()
    print('hello')
end








