package dbop

import (
	"encoding/base64"
	"errors"
	"fmt"

	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/serialize/proto"
	"github.com/go-redis/redis/v7"

	"mmo/common/applog"
	mymsg "mmo/messages"
)

const (
	AutoIncrUIdKey = "card.AutoUId"
)

// .简单在redis中插入一条 [username].id来判定是否注册过
// .通过AutoIncrUIdKey来作为自增长key
// .通过player.[id]来获取角色信息

func TryInitDB(client *redis.Client) error {
	uidKey := MakeRedisKey(AutoIncrUIdKey)
	_, err := client.Get(uidKey).Int64()
	if err == nil {
		// init already
		return nil
	}

	_, err = client.Set(uidKey, 1, 0).Result()
	return err
}

// Auth 不验证密码
func Auth(client *redis.Client, username, password string) int64 {
	// check if exsit username.id

	uid, err := client.Get(MakeRedisKey(fmt.Sprintf("%v.id", username))).Int64()
	if err == nil {
		// 认证成功
		return uid
	}

	if IsNilError(err) {
		return Register(client, username, password)
	}

	// 表示认证失败
	return 0
}

// Register 基于自增长id
func Register(client *redis.Client, username, password string) int64 {
	// 不存在会被初始化
	uid, err := client.Incr(MakeRedisKey(AutoIncrUIdKey)).Result()
	if err != nil {
		l.L.Errorf("RegisterCmdHandler auto incr error: %v", err)
		return 0
	}

	_, err = client.Set(MakeRedisKey(fmt.Sprintf("%v.id", username)), uid, 0).Result()
	if err != nil {
		l.L.Errorf("RegisterCmdHandler set error: %v", err)
		return 0
	}
	return uid
}

func LoadPlayer(client *redis.Client, uid int64) (*mymsg.PlayerInfo, error) {
	serializer := proto.GetDefaultSerializer()
	raw, err := client.Get(MakeRedisKey(fmt.Sprintf("player.%v", uid))).Result()
	if err != nil {
		l.L.Errorf("LoadPlayer failed: %v %v", uid, err)
		return nil, err
	}
	data, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		l.L.Errorf("LoadPlayer base64 DecodeString: %v %v", uid, err)
		return nil, err
	}

	info := &mymsg.PlayerInfo{}
	err = serializer.Unmarshal(data, info)
	if err != nil {
		l.L.Errorf("LoadPlayer proto Unmarshal failed: %v %v", uid, err)
		return nil, err
	}
	return info, nil
}

func SavePlayer(client *redis.Client, info *mymsg.PlayerInfo) bool {
	if info == nil {
		l.L.Errorf("SavePlayer info is nil")
		return false
	}
	serializer := proto.GetDefaultSerializer()
	data, err := serializer.Marshal(info)
	if err != nil {
		l.L.Errorf("SavePlayer marshal error: %v", err)
		return false
	}

	strBase64 := base64.StdEncoding.EncodeToString(data)
	if applog.PlayerDB != nil {
		applog.PlayerDB.Infof("save player size:%v", len(strBase64))
		applog.PlayerDB.Infof(" :%v", strBase64)
	}

	_, err = client.Set(MakeRedisKey(fmt.Sprintf("player.%v", info.Base.UId)),
		strBase64, 0).Result()
	if err != nil {
		l.L.Errorf("SavePlayer redis.Set error: %v", err)
		return false
	}
	return true
}

func IsNilError(err error) bool {
	return errors.Is(err, redis.Nil)
}
