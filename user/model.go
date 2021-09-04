package user

import (
	"time"
)

type Account struct {
	ID           uint64 `gorm:"primary_key"`
	UserName     string `gorm:"type:varchar(128);"`
	PasswordHash string `gorm:"type:varchar(128);"`
	Email        string `gorm:"type:varchar(128);"`
	Phone        string `gorm:"type:varchar(20)"`
	Bio          string `gorm:"type:varchar(20)"`
	Gender       bool
	Avatar       string
	Credit       int64
	IsVirtual    int
	Chn          string
	Status       int
	CreatedAt    time.Time `gorm:"type:timestamp;"`
	UpdatedAt    time.Time `gorm:"type:timestamp;"`
}

type ChannelAdmin struct {
	ID               uint64 `gorm:"primary_key"`
	TrackName        string `gorm:"type:varchar(30) ;"`
	ChannelName      string `gorm:"type:varchar(30);"`
	Password         string `gorm:"type:varchar(30);"`
	Remark           string `gorm:"type:varchar(100)"`
	Token            string `gorm:"type:varchar(60)"`
	Perms            string `gorm:"type:text"`
	GpListenUrl      string
	ChannelType      int
	WebChannelToken  string
	WebChannelConfig string
	CreatedAt        time.Time `gorm:"type:timestamp;"`
	UpdatedAt        time.Time `gorm:"type:timestamp;"`
}

type AccountToken struct {
	UserId    uint64 `gorm:"primary_key"`
	Token     string
	Extend    string
	CreatedAt time.Time `gorm:"type:timestamp;"`
	UpdatedAt time.Time `gorm:"type:timestamp;"`
}

//流水记录
type CreditRecord struct {
	ID        int64 `gorm:"primary_key;AUTO_INCREMENT"`
	UserId    uint64
	Type      int32 `gorm:"default:1"`
	Title     string
	Amount    int64
	Balance   int64
	Extend    string
	CreatedAt time.Time `gorm:"type:timestamp;"`
	UpdatedAt time.Time `gorm:"type:timestamp;"`
}

//保险箱
type SafeBox struct {
	UserId   uint64 `gorm:"primary_key"`
	Amount   int64
	Password string
}

//客户端登录的临时信息
type ClientInfo struct {
	IP   string `json:"ip"`
	Chn  string `json:"chn"`
	Addr string `json:"addr"`
}

type MixedInfo struct {
	*Account
	*ClientInfo
}
