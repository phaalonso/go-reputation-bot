package models

import (
	"github.com/jinzhu/gorm"
	"go-reputation-bot/pkg/config"
)

var db *gorm.DB

type UserReputation struct {
	gorm.Model
	ChatId     int64 `gorm:"primaryKey;autoIncrement:false"`
	UserId     int64 `gorm:"primaryKey;autoIncrement:false"`
	Reputation int32
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&UserReputation{})
}

func (u *UserReputation) CreateUserReputation() *UserReputation {
	db.NewRecord(u)
	db.Create(&u)

	return u
}

func GetUserReputationInChat(userId int64, chatId int64) *UserReputation {
	var ur UserReputation

	db.Where("ChatId = ? AND UserId = ?", chatId, userId).First(&ur)

	return &ur
}

func GetTotalUserReputation(userId int64) int32 {
	var reputation int32

	db.Where("UserId=?", userId).Group("UserId").Select("sum(Reputation) as total").Scan(&reputation)

	return reputation
}

func (u *UserReputation) UpdateUserReputation(userId int64, chatId int64, reputation int32) {
	db.Where("ChatId = ? AND UserId = ?", chatId, userId).Set("Reputation = ?", reputation)
}
