package slots

type SpinTrack struct {
	Type           string //spin
	Ts             int64
	SeqNo          string
	UserId         uint64
	WinInfo        string
	IsFreeGame     bool
	HasGetFreeGame bool
	Bet            int64
	Award          int64
	Water          int64
	WaterLine      float64
	Auto           int32
	Fast           int32
	PolicyNo       string
	PolicyInfo     string
	IsHit          bool
}

type UserStreamTrack struct {
	Type    string
	Ts      int64
	Term    string
	UserId  uint64
	Bet     int64
	Award   int64
	Credit  int64
	Tax     int64
	PayRate float32
	Area    map[int][]string
}
