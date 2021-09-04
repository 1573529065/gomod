package qiuqiu

type StartTrack struct {
	Type   string // start
	Ts     int64
	Room   int32 // room kind
	Table  uint64
	GameId string `json:"game_id"`
	Base   int64
}

type AnnounceTrack struct {
	Type        string // announce 结算
	Ts          int64
	UserId      uint64 `json:"user_id"`
	RoomKind    int32  // room kind
	GameId      string `json:"game_id"`
	BetAmount   int64  `json:"bet_amount"`   // 该玩家当局下注金额
	AwardAmount int64  `json:"award_amount"` // 赢家-纯盈利，减去了抽水, 输家-0
	TaxAmount   int64  `json:"tax_amount"`   // 税
}

type GameStreamTrack struct {
	Type       string   // announce 结算
	Ts         int64    //时间戳
	RoomKind   int32    //房间类型
	GameId     string   //GameId
	Tax        float64  //税
	Pool       int64    //下注池
	TableCards string   //tableCards json
	Uids       []uint64 //参与真实玩家列表
}

type Strategy struct {
	Type       string // strategy_user　　触发玩家策略
	Ts         int64
	Room       int32 // room kind
	Table      uint64
	GameId     string `json:"game_id"`
	UserId     uint64 `json:"user_id"`
	StrategyId string `json:"strategy_id"`
	Chair      int32  //
	Before     string //换牌前
	After      string //换牌后
	Match      bool   //是否符合
	Hit        bool   //是否命中
}
