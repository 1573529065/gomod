package ludo

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
	Cards       string
	AwardAmount int64  `json:"award_amount"` // 赢家-纯盈利，减去了抽水, 输家-0
	TaxAmount   int64  `json:"tax_amount"`
	BetAmount   int64  `json:"bet_amount"` // 赢家-纯盈利，减去了抽水, 输家-0
	Reason      int32  // 输赢原因
	ChairInfo   string // 座位信息
	ScoreInfo   string // 分数信息
}

type StrategyUser struct {
	Type       string
	Ts         int64
	UserId     string `json:"user_id"`
	Room       int32
	Table      uint64
	Chair      int32
	GameId     string
	ChesssStr  string
	StrategyId string
	Before     int32
	After      int32
	Hit        bool
	Match      bool
	IsEat      int32 // 默认是1，如果吃子，为2，没有吃子为3
	IsMid      int32 // 默认是1，如果进中心道，为2，没有进中心道为3
}

type BetTrack struct {
	Type      string // bet
	Ts        int64
	UserId    uint64 `json:"user_id"`
	UserType  int32  `json:"user_type"`
	Room      int32  // room kind
	Table     uint64
	GameId    string `json:"game_id"`
	BetAmount int64  `json:"bet_amount"`
}

type GameStreamTrack struct {
	Type       string //game stream
	Infos      map[string]*GameStream
	ChairInfos string
}

//牌局日志
type GameStream struct {
	Ts         int64
	RoomKind   int32
	Pool       int64
	RoomId     uint64
	GameId     string
	Tax        int32
	TableChess string
	Award      string
}

//丢牌日志
type DiscardStream struct {
	Type       string
	Ts         int64
	RoomKind   int32
	RoomId     uint64
	GameId     string
	Turn       int32
	TableChess string
}

type ChairChess struct {
	UserId     uint64
	RollNum    int32
	Step       int32
	ChessIndex int32
	HandCard   string
	Before     string
}

type ChessInfo struct {
	Chess  string
	Chair  int32
	UserId uint64
	Score  int32
	Award  int64
}
