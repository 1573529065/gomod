package cache

import (
	"bytes"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"lol.com/server/nest.git/log"
	"lol.com/server/nest.git/tools"
	"lol.com/server/nest.git/tools/jsonutils"
	"time"
)

type MsgReceive struct {
	Type int32                `json:"type"`
	Data jsonutils.JsonObject `json:"data"`
}

type userStatus struct {
	Game   int32  `json:"game"`
	Room   int32  `json:"room"`
	Server string `json:"server"`
	Ts     int64  `json:"ts"`
}

const (
	SystemError = 5000
	UnComplate  = 3000
	BetTimeOut  = 4000
	Success     = 1000
)

func (cache *PublicCache) SaveUser(userID uint64, roomKind int32) (bool, error) {
	rd := cache.pool.Get()
	defer rd.Close()
	script := redis.NewScript(2, saveUserScript)
	status := userStatus{
		Game:   cache.gameID,
		Room:   roomKind,
		Server: cache.serverURL,
		Ts:     time.Now().Unix(),
	}
	toSave, _ := json.Marshal(status)
	isNew, err := redis.Bool(script.Do(rd, statusKey, cache.serverKey, userID,
		string(toSave)))
	if err != nil {
		return false, err
	}
	if !isNew {
		return false, nil
	}
	return true, nil
}

func (cache *PublicCache) RemoveUser(userID uint64) error {
	rd := cache.pool.Get()
	defer rd.Close()
	rd.Send("MULTI")
	rd.Send("HDEL", statusKey, userID)
	rd.Send("SREM", cache.serverKey, userID)
	if _, err := rd.Do("EXEC"); err != nil {
		return err
	}
	return nil
}

func (cache *PublicCache) LoadUserServer(id uint64) (*userStatus, error) {
	rd := cache.pool.Get()
	defer rd.Close()
	value, err := redis.String(rd.Do("HGET", statusKey, id))
	if err == redis.ErrNil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	status := &userStatus{}
	if err = json.Unmarshal([]byte(value), status); err != nil {
		return nil, err
	}
	return status, nil
}

func (cache *PublicCache) CleanLoginInfo() {
	rd := cache.pool.Get()
	defer rd.Close()
	users, _ := redis.Int64s(rd.Do("SMEMBERS", cache.serverKey))
	if users != nil && len(users) > 0 {
		log.Warn("clean user in public redis of key:%s", cache.serverKey)
	}
	for _, user := range users {
		rd.Do("HDEL", statusKey, user)
		log.Debug("cleaning user %v from redis", user)
	}
	rd.Do("DEL", cache.serverKey)
}

//超时时间
func (cache *PublicCache) SyncAmount(timeOut int32, channels []string, uuid string, toPub string, result SinglePurseBetTimeout, streamChan chan string) int32 {
	var (
		err error
	)
	//Get a connection from a pool
	c := cache.pool.Get()
	//把下注（奖励）通知塞入channel，然后监听channel看是否有回调
	if _, err := c.Do("PUBLISH", channels[0], toPub); err != nil {
		log.Error("pub error:%v", err.Error())
		return SystemError
	}
	var initData = SinglePurseBetResult{}
	defer tools.RecoverFromPanic(nil) //panic 注释，调试
	closeChan := make(chan bool)
	//开始计时，假如超时还未处理，那么把超时标记传到channel内
	go timer(result, closeChan, streamChan)
	for {
		select {
		case <-closeChan:
			log.Info("break")
			return BetTimeOut
		default:
			psc := redis.PubSubConn{Conn: c}
			// Set up subscriptions
			err = psc.Subscribe(channels[1])
			if err != nil {
				log.Error("can't subscribe channel from im redis!!!!")
				continue
			}
			switch v := psc.Receive().(type) {
			case redis.Message:
				var data = initData
				//避免把长整形打印成科学计数法
				log.Info("receive message:%v", v.Data)
				d := json.NewDecoder(bytes.NewReader(v.Data))
				d.UseNumber()
				if err = d.Decode(&data); err == nil {
					log.Info("receive message:%v", data)
					rUuid := data.TransactionId
					log.Info("receive uuid:%v", rUuid)
					status := data.Status
					code := data.Code
					log.Info("receive status:%v", status)
					//假如收到失败回调，那么一样停止计时器
					if status == 4 {
						closeChn(psc, channels[1], c, closeChan)
						if code == UnComplate {
							return UnComplate
						}
						return code
					}
					if uuid == rUuid {
						//处理开始时，停掉计时器，停止订阅，塞入下注订阅频道，然后结束
						closeChn(psc, channels[1], c, closeChan)
						return Success
					}
				} else {
					log.Warn("received im json format error:%v", err.Error())
					closeChn(psc, channels[1], c, closeChan)
					return SystemError
				}
			case error:
				log.Error(v.Error())
				closeChn(psc, channels[1], c, closeChan)
				return SystemError
			}
		}
	}
}

func timer(result SinglePurseBetTimeout, closeChan chan bool, streamChan chan string) {
	timer1 := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-closeChan:
			log.Info("stop channel close")
			timer1.Stop()
		case <-timer1.C:
			//超时塞入下注失败stream，并且停止监听
			result.Ts = time.Now().Unix()
			StreamTrack(streamChan, result)
			log.Info("stop channel close time out")
			timer1.Stop()
			closeChan <- true
			return
		}
	}
}

func closeChn(psc redis.PubSubConn, channel string, c redis.Conn, closeChan chan bool) {
	_ = psc.Unsubscribe(channel)
	_ = psc.Close()
	_ = c.Close()
	closeChan <- true
}
