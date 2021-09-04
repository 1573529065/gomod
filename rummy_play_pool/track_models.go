package rummy_play_pool

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
	GameScore  int32
}

type LoginTrack struct {
	Type   string // login
	Ts     int64
	UserId uint64
	GameId string
	Addr   string
}

type LogoutTrack struct {
	Type   string // logout
	Ts     int64
	UserId uint64
	GameId string
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
	CardsStr   string
	StrategyId string
	Before     string
	After      string
	Hit        bool
	Match      bool
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
	Type  string //game stream
	Infos map[string]*GameStream
}

//牌局日志
type GameStream struct {
	Ts         int64
	RoomKind   int32
	Pool       int64
	RoomId     uint64
	GameId     string
	Tax        int32
	TableCards string
	Award      string
}

//丢牌日志
type DiscardStream struct {
	Type        string
	Ts          int64
	RoomKind    int32
	RoomId      uint64
	GameId      string
	GameIdBig   string
	GameIdSmall string
	Turn        int32
	TableCards  string
	WildCard    string
}

type ChairCard struct {
	UserId    uint64
	Visible   bool
	TouchCard string
	Discard   string
	FirstTurn bool
	Fold      bool
	HandCard  string
	Before    string
}

type CardsInfo struct {
	Cards  string
	Chair  int32
	UserId uint64
	Score  int32
	Award  int64
}
