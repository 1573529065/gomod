package cheatscore

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/gomodule/redigo/redis"
	"lol.com/server/nest.git/tools/tz"
)

const (
	onlineExpire = 5
)

type UserScoreDAO struct {
	Session *mgo.Session
}

func MarkIn(uid uint64, pool *redis.Pool, gameId string) {
	//记录5秒缓存时间的用户标记
	key := fmt.Sprintf("mark-user-in:%d-%d-%s", uid, gameId)
	rd := pool.Get()
	defer rd.Close()
	_, _ = rd.Do("SETEX", key, onlineExpire, tz.GetNowTs())
}

func MarkTouch(uid uint64, pool *redis.Pool, gameId string) {
	//记录5秒缓存时间的用户标记
	key := fmt.Sprintf("mark-user-in:%d-%d-%s", uid, gameId)
	rd := pool.Get()
	defer rd.Close()
	_, _ = rd.Do("SETEX", key, onlineExpire, tz.GetNowTs())
}

func CheckMark(uid uint64, pool *redis.Pool, gameId string) bool {
	//记录5秒缓存时间的用户标记
	key := fmt.Sprintf("mark-user-in:%d-%s", uid, gameId)
	rd := pool.Get()
	defer rd.Close()
	exists, _ := redis.Bool(rd.Do("EXISTS", key))
	if exists {
		return false
	}
	return true
}
