package rummy_play

type SyncTrack struct {
	Type   string // sync
	Ts     int64
	UserId uint64 `json:"user_id"`
	Credit int64
	After  int64
	GameId string
	Bet    bool
}

type StartTrack struct {
	Type       string // start
	Ts         int64
	Room       int32 // room kind
	Table      uint64
	GameId     string `json:"game_id"`
	Base       int64
	UserCount  int32
	RobotCount int32
}

type GameStreamTrack struct {
	Type       string // game stream
	Ts         int64
	RoomKind   int32
	RoomId     uint64
	GameId     string
	UserCount  int32
	RobotCount int32
	Tax        int32
	TableCards string
}

type AnnounceTrack struct {
	Type        string // announce 结算
	Ts          int64
	UserId      uint64 `json:"user_id"`
	UserType    int32  `json:"user_type"`
	Room        int32  // room kind
	Table       uint64
	GameId      string `json:"game_id"`
	Seat        int32
	UserCount   int32
	Base        int32
	Cards       string
	WildCard    int32
	AwardAmount int64 `json:"award_amount"` // 赢家-纯盈利，减去了抽水, 输家-0
	TaxAmount   int64 `json:"tax_amount"`
}

type StrategyUser struct {
	Type       string
	Ts         int64
	UserId     uint64
	Room       int32
	Table      uint64
	Chair      int32
	GameId     string
	StrategyId string
	CardsStr   string
	Before     string
	After      string
	Hit        bool
	Match      bool
}
