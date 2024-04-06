package models

import (
	"github.com/jinzhu/gorm"
	"github.com/phaalonso/go-reputation-bot/pkg/config"
)

var db *gorm.DB

type UserReputation struct {
	gorm.Model
	ChatID     int64 `gorm:"primaryKey;autoIncrement:false"`
	UserID     int64 `gorm:"primaryKey;autoIncrement:false"`
	Reputation int32
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.Debug().AutoMigrate(&UserReputation{})
}

func UpdateOrCreateReputation(chatId int64, userId int64) *UserReputation {
	rep := GetUserReputationInChat(chatId, userId)

	if rep == nil {
		rep = &UserReputation{
			ChatID:     chatId,
			UserID:     userId,
			Reputation: 1,
		}

		return rep.CreateUserReputation()
	}

	rep.Reputation += 1
	rep.UpdateUserReputation(rep.Reputation)

	return rep
}

func (u *UserReputation) CreateUserReputation() *UserReputation {
	db.NewRecord(u)
	db.Create(&u)

	return u
}

func GetUserReputationInChat(chatId int64, userId int64) *UserReputation {
	var ur UserReputation

	db.Where("ChatID=? AND UserID=?", chatId, userId).First(&ur)

	return &ur
}

func GetTotalUserReputation(userId int64) int32 {
	var reputation int32

	db.Where("UserID=?", userId).Group("UserID").Select("sum(Reputation) as total").Scan(&reputation)

	return reputation
}

func (u *UserReputation) UpdateUserReputation(reputation int32) {
	db.Where("ChatID=? AND UserID=?", u.ChatID, u.UserID).Set("Reputation=?", reputation)
}
