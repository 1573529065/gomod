package fishing

type UserStreamTrack struct {
	Type      string
	Ts        int64
	UserId    uint64
	RoomKind  int32
	Bet       int64
	Award     int64
	Credit    int64
	PayRate   int32
	FishIds   []uint64
	Kill      bool
	BonusRate float32
}

type StrategySys struct {
	Type       string
	Ts         int64
	UserId     uint64
	Water      float64
	Interval   float64
	Chair      int32
	StrategyId int
	Strategy   string
	RoomKind   int32
	Hit        bool
}
