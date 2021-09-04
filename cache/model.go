package cache

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gogf/gf/container/glist"
	"github.com/gomodule/redigo/redis"
	"lol.com/server/nest.git/log"
	"strconv"
	"time"
)

//所有游戏的公用redis
type PublicCache struct {
	pool        *redis.Pool
	gameID      int32
	serverURL   string
	serverKey   string
	rpcCallback map[int32]func([]byte)
	listened    bool
}

type PublicStream struct {
	pool       *redis.Pool
	streamName string
	streamChan chan string
	publicChan glist.List
}

type SinglePurseBet struct {
	Ts            int64
	TransactionId string
	UserId        uint64
	GameId        int32
	Bet           int64
	Term          string
}

type SinglePurseBetResult struct {
	TransactionId string
	UserId        uint64
	GameId        int32
	Bet           int64
	Status        int32
	Balance       int64
	Code          int32
}

type SinglePurseAward struct {
	Type          string
	Ts            int64
	TransactionId string
	UserId        uint64
	GameId        int32
	Award         int64
	Term          string
}

type SinglePurseBetTimeout struct {
	Ts            int64
	Type          string
	TransactionId string
	UserId        uint64
	GameId        int32
	Bet           int64
}

type GameStreamTrack struct {
	Type       string  `bson:"type"`
	Ts         int64   `bson:"ts"`
	RoomKind   int32   `bson:"room_kind"`
	RoomId     uint64  `bson:"room_id"`
	GameId     string  `bson:"game_id"`
	Tax        float64 `bson:"tax"`
	TableCards string  `bson:"table_cards"`
}

type TableCard struct {
	Cards  string
	Chair  int32
	UserId uint64
	Score  int32
	Award  int32
}

const (
	ImChannel = "lucky:im:channel"
	statusKey = "lucky:user:status"
	//新用户则存放，否则啥也不做
	saveUserScript = `local cur=redis.call('HSETNX', KEYS[1], ARGV[1], ARGV[2]);if(cur==1) then redis.call('SADD', KEYS[2], ARGV[1]); return 1; end return 0;`
	StreamField    = "data"
)

func prefixKey(key string, prefix string) string {
	return fmt.Sprintf("%s:%s", prefix, key)
}

//这里不完全封装，不负责pool的生命周期，以便各项目自由使用公用cache
func NewPublicCache(pool *redis.Pool, gameID int32, prefix string, serverURL string) *PublicCache {
	serverKey := prefixKey(fmt.Sprintf("server:%v", serverURL), prefix)
	return &PublicCache{
		pool:        pool,
		gameID:      gameID,
		serverURL:   serverURL,
		serverKey:   serverKey,
		rpcCallback: make(map[int32]func([]byte)),
	}
}

func NewPublicStream(pool *redis.Pool, streamName string, streamChan chan string, publicChan *glist.List) *PublicStream {
	return &PublicStream{
		pool:       pool,
		streamName: streamName,
		streamChan: streamChan,
		publicChan: *publicChan,
	}
}

//not goroutine safe
func (cache *PublicCache) RegisterRPCCallback(eventID int32, fn func([]byte)) {
	cache.rpcCallback[eventID] = fn
}

//not goroutine safe
func (cache *PublicCache) UnRegisterRPCCallback(eventID int32) {
	delete(cache.rpcCallback, eventID)
}

func InitPublicStreams(publicStream *PublicStream) {
	for i := 0; i < 8; i++ {
		go func() {
			for {
				select {
				case s := <-publicStream.streamChan:
					func() {
						publicStream.publicChan.PushBack(s)
					}()
				}
			}
		}()
	}
	for i := 0; i < 10; i++ {
		var n *list.Element
		go func() {
			for {
				var con = publicStream.pool.Get()
				publicStream.publicChan.LockFunc(func(list *list.List) {
					length := list.Len()
					if length > 0 {
						for e := list.Front(); e != nil; e = n {
							con.Send("XADD", publicStream.streamName, "MAXLEN", "~", "80000000", "*", StreamField, e.Value)
							n = e.Next()
							list.Remove(e)
						}
					}
				})
				con.Flush()
				con.Close()
				time.Sleep(1 * time.Second)
			}
		}()
	}
}

func StreamTrack(streamChan chan string, info interface{}) {
	s, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}
	streamChan <- string(s)
}

func listenLockUser(gameId string, gameUrl string, keyPrefix string, pool *redis.Pool, mg *mgo.Session) {
	rd := pool.Get()
	defer rd.Close()
	session := mg.Clone()
	defer session.Close()
	//先找出所有redis内的用户锁
	statusList, err := redis.StringMap(rd.Do("HGETALL", statusKey))
	if err != redis.ErrNil {
		//找出所有时间超过10分钟，并且20分钟内没有game_stream的用户
	statusLoop:
		for uid, status := range statusList {
			user := &userStatus{}
			intUid, _ := strconv.Atoi(uid)
			uid := uint64(intUid)
			if err = json.Unmarshal([]byte(status), user); err != nil {
				log.Error("json format error :%v,%v", err.Error(), uid)
				continue
			}
			//针对rummy,rummy_pool,teen这三个游戏
			if user.Game == 110 || user.Game == 111 || user.Game == 102 {
				var (
					tableCard map[string]TableCard
					result    []GameStreamTrack
				)
				now := time.Now().Unix()
				if now-user.Ts > 10*60 {
					//根据uid去查询game_stream
					c := session.DB("game").C(gameId)
					start := time.Now().Add(time.Duration(-20) * time.Minute).Unix()
					err := c.Find(bson.M{"ts": bson.M{"$gt": start, "$lt": now}}).All(&result)
					if err != nil && err != mgo.ErrNotFound {
						log.Info("mgo error :%v", err.Error())
						continue
					}
					for _, res := range result {
						if err = json.Unmarshal([]byte(res.TableCards), &tableCard); err != nil {
							log.Info("json format error :%v,%v", err.Error(), uid)
							continue
						}
						for _, card := range tableCard {
							if uid == card.UserId {
								continue statusLoop
							}
						}
					}
					//符合条件的用户进行删除金币锁,并且删除牌局
					serverKey := prefixKey(fmt.Sprintf("server:%v", gameUrl), keyPrefix)
					rd.Send("MULTI")
					rd.Send("HDEL", statusKey, uid)
					rd.Send("SREM", serverKey, uid)
					if _, err := rd.Do("EXEC"); err != nil {
						log.Info("remove user error :%v", err.Error())
						continue
					}
					log.Info("clean uid:%v", uid)
				}
			}
		}
	}
}
