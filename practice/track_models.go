package practice

type AnnounceTrack struct {
	Type        string // announce 结算
	Ts          int64
	UserId      uint64 `json:"user_id"`
	UserType    int32  `json:"user_type"`
	Room        int32  // room kind
	Table       uint64
	GameId      string `json:"game_id"`
	Seat        int32
	Base        int32
	Cards       string
	WildCard    int32
	AwardAmount int64 `json:"award_amount"` // 赢家-纯盈利，减去了抽水, 输家-0
	TaxAmount   int64 `json:"tax_amount"`
}
