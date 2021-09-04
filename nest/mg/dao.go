package mg

import (
	"fmt"
	"lol.com/server/nest.git/tools/mem"
	"lol.com/server/nest.git/tools/ternary"
	"math"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"lol.com/server/nest.git/log"
	"lol.com/server/nest.git/tools/tz"
)

type UserStatsDAO struct {
	Session *mgo.Session
}

func (dao *UserStatsDAO) GetUsersStats(uids []uint64) *[]UserStats {
	session := dao.Session.Clone()
	defer session.Close()
	c := dao.Session.DB("game").C("user_stats")
	var result []UserStats
	err := c.Find(bson.M{"_id": bson.M{"$in": uids}}).All(&result)
	if err != nil && err != mgo.ErrNotFound {
		log.Error("load from mongo error:%v", err.Error())
	}
	return &result
}

func (dao *UserStatsDAO) GetUserStats(uid uint64) *UserStats {
	session := dao.Session.Clone()
	defer session.Close()
	c := dao.Session.DB("game").C("user_stats")
	var result UserStats
	err := c.Find(bson.M{"_id": uid}).One(&result)
	if err != nil && err != mgo.ErrNotFound {
		log.Error("load from mongo error:%v", err.Error())
	}
	return &result
}

func (dao *UserStatsDAO) GetDailyStats(uid uint64) *DailyStats {
	session := dao.Session.Clone()
	defer session.Close()
	c := dao.Session.DB("game").C("daily_stats")
	var result DailyStats
	today := tz.GetTodayStr()
	err := c.Find(bson.M{"_id": fmt.Sprintf("%d-%s", uid, today)}).One(&result)
	if err != nil && err != mgo.ErrNotFound {
		log.Error("load from mongo error:%v", err.Error())
	}
	return &result
}

//统计设备上的账户数
func (dao *UserStatsDAO) CountDeviceAccount(aid string) int {
	session := dao.Session.Clone()
	defer session.Close()
	c := dao.Session.DB("game").C("user_stats")
	n, _ := c.Find(bson.M{"aid": aid}).Count()
	return n
}

func (dao *UserStatsDAO) ScoreTrack(totalScore int64, score int64, uid uint64, game int32, flowId string, roomType int32, cheatId int32) {
	dataLog := bson.M{"game": game, "created_at": bson.Now(), "room_type": roomType, "cheat_id": cheatId, "uid": uid,
		"game_id": flowId, "score": score, "total_score": totalScore}
	dao.Session.DB("game").C("user_score_log").Insert(dataLog)
}

func (dao *UserStatsDAO) CheatTrack(uid uint64, game int32, date string) {
	key := fmt.Sprintf("%d-%v", uid, date)
	data := bson.M{"$inc": bson.M{"count": 1}, "$set": bson.M{"_id": key, "uid": uid, "game": game, "date": date}}
	dao.Session.DB("game").C("user_cheat").UpsertId(key, data)
}

func (dao *UserStatsDAO) GetScore(uid uint64, game int32) *UserScore {
	c := dao.Session.DB("game").C("user_score")
	var result UserScore
	err := c.Find(bson.M{"_id": uid}).One(&result)
	if err != nil && err != mgo.ErrNotFound {
		log.Error("load from mongo error:%v", err.Error())
	}
	return &result
}

var templateCache = mem.NewTTLCache(60, 10)

//获取发送闪告的最小金额，单位是文
func (dao *UserStatsDAO) GetBroadcastMinAmount(gameId string) int64 {
	var minAmount int64
	if cached, err := templateCache.Get(gameId); err == nil {
		minAmount = cached.(int64)
		return minAmount
	} else {
		session := dao.Session.Clone()
		defer session.Close()
		c := dao.Session.DB("de").C("im_template")
		var template IMTemplate
		err := c.Find(bson.M{"type": gameId}).Sort("min_amount").Limit(1).One(&template)
		if err == nil {
			templateCache.Set(gameId, template.MinAmount)
			return template.MinAmount
		} else if err == mgo.ErrNotFound {
			templateCache.Set(gameId, int64(math.MaxInt64))
			return math.MaxInt64
		} else {
			log.Error("can't get im_template from mongo:%v", err.Error())
			return math.MaxInt64
		}
	}
}

type IpAids struct {
	ID   string   `bson:"_id"`
	IPS  []string `bson:"ips"`
	AIDS []string `bson:"aids"`
}

type Withdraw struct {
	ID                string   `bson:"_id"`
	AccountHolderName []string `bson:"accountHolderName"`
	AccountNumber     []string `bson:"accountNumber"`
	PayeeEmail        []string `bson:"payeeEmail"`
	PayeeMobile       []string `bson:"payeeMobile"`
}

type Ids struct {
	ID  string   `bson:"_id"`
	IDS []uint64 `bson:"ids"`
}

//匹配，确定桌子以后，那么跟入桌的玩家进行匹配信息核对，返回匹配结果，infos不为空。
func (dao *UserStatsDAO) MatchingPriority(uid uint64, uids []uint64) bool {
	session := dao.Session.Clone()
	defer session.Close()
	c := dao.Session.DB("game").C("user_date_ip_aid")
	c1 := dao.Session.DB("game").C("withdraw_info")
	userStats := dao.GetUserStats(uid)
	userStats1 := dao.GetUserStats(uids[0])
	isPay := userStats.Recharge.Total > 0
	isPay1 := userStats1.Recharge.Total > 0
	if isPay != isPay1 {
		return false
	}
	var results IpAids
	var results1, results2, results4, results5, results6, results7 Ids
	var results3 Withdraw
	o := bson.M{"$match": bson.M{"user_id": uid}}
	o1 := bson.M{"$group": bson.M{"_id": nil, "ips": bson.M{"$addToSet": "$ip"}, "aids": bson.M{"$addToSet": "$aid"}}}
	operations := []bson.M{o, o1}
	pipe := c.Pipe(operations)
	err := pipe.One(&results)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	aids := results.AIDS
	ips := results.IPS
	//通过aids 关联出结果
	if len(aids) > 0 {
		o2 := bson.M{"$match": bson.M{"aid": bson.M{"$in": aids}}}
		o3 := bson.M{"$group": bson.M{"_id": nil, "ids": bson.M{"$addToSet": "$user_id"}}}
		operations1 := []bson.M{o2, o3}
		pipe1 := c.Pipe(operations1)
		err = pipe1.One(&results1)
		if err != nil && err != mgo.ErrNotFound {
			panic(err)
		}
		ids := results1.IDS
		//查看所有的uids是否在result中，假如在那么说明有关联
		for _, id := range uids {
			for _, result := range ids {
				if id == result {
					return false
				}
			}
		}
	}
	//通过ips 关联出结果
	if len(ips) > 0 {
		o4 := bson.M{"$match": bson.M{"ip": bson.M{"$in": ips}}}
		o5 := bson.M{"$group": bson.M{"_id": nil, "ids": bson.M{"$addToSet": "$user_id"}}}
		operations2 := []bson.M{o4, o5}
		pipe2 := c.Pipe(operations2)
		err = pipe2.One(&results2)
		if err != nil {
			panic(err)
		}
		ids1 := results2.IDS
		//查看所有的uids是否在result2中，假如在那么说明有关联
		for _, id := range uids {
			for _, result := range ids1 {
				if id == result {
					return false
				}
			}
		}
	}
	//通过提现 关联出结果
	o6 := bson.M{"$match": bson.M{"user_id": uid}}
	o7 := bson.M{"$group": bson.M{"_id": nil, "accountHolderName": bson.M{"$addToSet": "$accountHolderName"},
		"accountNumber": bson.M{"$addToSet": "$accountNumber"},
		"payeeEmail":    bson.M{"$addToSet": "$payeeEmail"},
		"payeeMobile":   bson.M{"$addToSet": "$payeeMobile"},
	}}
	operations3 := []bson.M{o6, o7}
	pipe3 := c1.Pipe(operations3)
	err = pipe3.One(&results3)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	ahn := ternary.If(len(results3.AccountHolderName) > 0, results3.AccountHolderName, []string{}).([]string)
	an := ternary.If(len(results3.AccountNumber) > 0, results3.AccountNumber, []string{}).([]string)
	payEmails := ternary.If(len(results3.PayeeEmail) > 0, results3.PayeeEmail, []string{}).([]string)
	payMobiles := ternary.If(len(results3.PayeeMobile) > 0, results3.PayeeMobile, []string{}).([]string)
	o8 := bson.M{"$match": bson.M{"user_id": bson.M{"$ne": uid}, "accoundHolderName": bson.M{"$in": ahn}}}
	o9 := bson.M{"$group": bson.M{"_id": nil, "ids": bson.M{"$addToSet": "$user_id"}}}
	operations4 := []bson.M{o8, o9}
	pipe4 := c1.Pipe(operations4)
	err = pipe4.One(&results4)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	ids2 := results4.IDS
	//查看所有的uids是否在result中，假如在那么说明有关联
	for _, id := range uids {
		for _, result := range ids2 {
			if id == result {
				return false
			}
		}
	}
	o10 := bson.M{"$match": bson.M{"user_id": bson.M{"$ne": uid}, "accountNumber": bson.M{"$in": an}}}
	o11 := bson.M{"$group": bson.M{"_id": nil, "ids": bson.M{"$addToSet": "$user_id"}}}
	operations5 := []bson.M{o10, o11}
	pipe5 := c1.Pipe(operations5)
	err = pipe5.One(&results5)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	ids3 := results5.IDS
	//查看所有的uids是否在result中，假如在那么说明有关联
	for _, id := range uids {
		for _, result := range ids3 {
			if id == result {
				return false
			}
		}
	}
	o12 := bson.M{"$match": bson.M{"user_id": bson.M{"$ne": uid}, "payeeEmail": bson.M{"$in": payEmails}}}
	o13 := bson.M{"$group": bson.M{"_id": nil, "ids": bson.M{"$addToSet": "$user_id"}}}
	operations6 := []bson.M{o12, o13}
	pipe6 := c1.Pipe(operations6)
	err = pipe6.One(&results6)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	ids4 := results6.IDS
	//查看所有的uids是否在result中，假如在那么说明有关联
	for _, id := range uids {
		for _, result := range ids4 {
			if id == result {
				return false
			}
		}
	}
	o14 := bson.M{"$match": bson.M{"user_id": bson.M{"$ne": uid}, "payeeMobile": bson.M{"$in": payMobiles}}}
	o15 := bson.M{"$group": bson.M{"_id": nil, "ids": bson.M{"$addToSet": "$user_id"}}}
	operations7 := []bson.M{o14, o15}
	pipe7 := c1.Pipe(operations7)
	err = pipe7.One(&results7)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	ids5 := results7.IDS
	//查看所有的uids是否在result中，假如在那么说明有关联
	for _, id := range uids {
		for _, result := range ids5 {
			if id == result {
				return false
			}
		}
	}
	return true
}
