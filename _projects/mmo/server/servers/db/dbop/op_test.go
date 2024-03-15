package dbop

import (
	"errors"
	"log"
	"testing"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/magiconair/properties/assert"

	mymsg "mmo/messages"
)

// 	start a redis server before run testing
func Test_InitDB(t *testing.T) {
	c := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	TryInitDB(c)

	_, err := c.Get(AutoIncrUIdKey).Result()
	assert.Equal(t, true, err == nil)
	c.Close()
}

func Test_AutoIncr(t *testing.T) {
	c := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	TryInitDB(c)

	uid := Auth(c, "dfk", "")
	c.Close()

	assert.Equal(t, true, uid > 0)
}

func Test_Err(t *testing.T) {
	c := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	c.Del(AutoIncrUIdKey)
	c.Del("dfk.id")
	c.Set(AutoIncrUIdKey, "fef", 0)

	uid := Auth(c, "dfk", "")
	c.Close()

	assert.Equal(t, true, uid == 0)
}

// 测试redis 是否会自动连接
// 结论是，每次操作都会获取连接，如果没有就自动连接

func TestRedisCon(t *testing.T) {
	c := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	running := true
	go func() {
		for running {
			time.Sleep(time.Second)
			log.Println("begin get")
			_, err := c.Get("hello").Result()
			log.Println("end get")
			if err != nil {
				log.Println(err)
			}

		}
	}()

	time.Sleep(9 * time.Second)
}

func TestRedisIdleConn(t *testing.T) {
	c := redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:6379",
		Password:     "", // no password set
		DB:           0,  // use default DB
		MinIdleConns: 1,
	})

	running := true
	go func() {
		for running {
			time.Sleep(2 * time.Second)
			log.Println("begin get")
			_, err := c.Get("hello").Result()
			log.Println("end get")
			if err != nil {
				log.Println(err)
			}

		}
	}()

	time.Sleep(9 * time.Second)
}

func Test_LoadSave(t *testing.T) {
	c := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	info := &mymsg.PlayerInfo{}
	info.Base = &mymsg.BaseInfo{}
	info.Base.UId = 10
	info.Base.Level = 99

	SavePlayer(c, info)

	info1, _ := LoadPlayer(c, 10)
	log.Printf("%v", info1)

	c.Close()

	assert.Equal(t, int32(99), info1.Base.Level)
}

func Test_LoadBadInfo(t *testing.T) {
	c := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	c.Set("player.10", "xxxxx", 0)

	info1, _ := LoadPlayer(c, 10)
	log.Printf("%v", info1)

	c.Close()

	assert.Equal(t, true, info1 == nil)
}

func TestNonExsit(t *testing.T) {
	c := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	c.Del("_no_exsit_")
	_, err := c.Get("_no_exsit_").Result()

	// 服务器没链接上时
	// *net.OpError
	// proto.RedisError
	// proto.Nil
	if errors.Is(err, redis.Nil) {
		log.Println("nil")
	}

	log.Printf("%v", err)

}
