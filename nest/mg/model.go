package mg

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

//mongo官方库暂不稳定，使用mgo

type GoldOP struct {
	Count   int32     `bson:"count"`
	Total   int64     `bson:"total"`
	LastAt  time.Time `bson:"last_at"`
	FirstAt time.Time `bson:"first_at"`
}

type AmountOP struct {
	Amount int64 `bson:"amount"`
}

type UserStats struct {
	ID         uint64    `bson:"_id"`
	RegisterAt time.Time `bson:"register_at"`
	Channel    string    `bson:"chn"`
	IP         string    `bson:"ip"`
	AID        string    `bson:"aid"`
	TotalGain  int64     `bson:"total_gain"`
	Recharge   GoldOP    `bson:"recharge"`
	TransferIn AmountOP  `bson:"transfer_in"`
	Withdraw   GoldOP    `bson:"withdraw"`
	TeenPatti  GameOP    `bson:"teen_patti"`
	RummyPlay  GameOP    `bson:"rummy_play"`
	RummyPool  GameOP    `bson:"rummy_pool"`
	Slots      GameOP    `bson:"slots"`
	QiuQiu     GameOP    `bson:"qiuqiu"`
	LuDo       GameOP    `bson:"ludo"`
}

type UserScore struct {
	ID    uint64 `bson:"_id"`
	Score int64  `bson:"score"`
}

type GameOP struct {
	Total Total `bson:"total"`
}

type Total struct {
	UserCount int64 `bson:"user_count"`
	Gain      int64 `bson:"gain"`
}

type DailyStats struct {
	ID        string `bson:"_id"`
	Channel   string `bson:"chn"`
	IP        string `bson:"ip"`
	AID       string `bson:"aid"`
	TotalGain int64  `bson:"total_gain"`
	Recharge  GoldOP `bson:"recharge"`
	Withdraw  GoldOP `bson:"withdraw"`
}

type IMTemplate struct {
	ID         bson.ObjectId `bson:"_id"`
	Template   string        `bson:"template"`
	Type       []string      `bson:"type"`
	MinAmount  int64         `bson:"min_amount"`
	ExcludeChn []string      `bson:"exclude_chn"`
	ExcludePkg []string      `bson:"exclude_pkg"`
	LastMod    int64         `bson:"last_modified"`
}
