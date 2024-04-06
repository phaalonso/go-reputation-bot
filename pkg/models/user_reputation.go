package models

import (
	"github.com/jinzhu/gorm"
	"github.com/phaalonso/go-reputation-bot/pkg/config"
)

var db *gorm.DB

type UserReputation struct {
	ChatID     int64 `gorm:"primaryKey;autoIncrement:false"`
	UserID     int64 `gorm:"primaryKey;autoIncrement:false"`
	Reputation int32
}

func init() {
	config.Connect()
	db = config.GetDB().Debug()
	db.AutoMigrate(&UserReputation{})
}

func UpdateOrCreateReputation(chatId int64, userId int64) *UserReputation {
	rep, err := GetUserReputationInChat(chatId, userId)

	if err != nil {
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

func GetUserReputationInChat(chatId int64, userId int64) (*UserReputation, error) {
	var ur UserReputation

	d := db.Where("chat_id=?", chatId).Where("user_id=?", userId).First(&ur)

	return &ur, d.Error
}

func GetTotalUserReputation(userId int64) (int32, error) {
	var reputation int32

	d := db.Model(&UserReputation{}).Where("user_id=?", userId).Select("sum(reputation) as total")

	if d.Error != nil {
		return 0, d.Error
	} else {
		d.Row().Scan(&reputation)

		return reputation, nil
	}
}

func (u *UserReputation) UpdateUserReputation(reputation int32) {
	db.Model(&u).Update("reputation", reputation)
}
