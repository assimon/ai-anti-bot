package database

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserInfo struct {
	BaseModel
	TelegramUserId    int64           `gorm:"uniqueIndex:idx_chat_user,column:telegram_user_id" json:"telegram_user_id"`
	TelegramChatId    int64           `gorm:"uniqueIndex:idx_chat_user,column:telegram_chat_id" json:"telegram_chat_id"`
	JoinedTime        carbon.DateTime `gorm:"column:joined_time" json:"joined_time"`
	NumberOfSpeeches  int64           `gorm:"column:number_of_speeches" json:"number_of_speeches"`
	VerificationTimes int64           `gorm:"column:verification_times" json:"verification_times"`
}

func (UserInfo) TableName() string {
	return "user_info"
}

func SaveUserInfo(user *UserInfo) error {
	err := Gdb.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "telegram_user_id"}, {Name: "telegram_chat_id"}},
		UpdateAll: true,
	}).Save(user).Error
	return err
}

func IncrementNumberOfSpeeches(user *UserInfo) error {
	err := Gdb.Transaction(func(tx *gorm.DB) error {
		var u UserInfo
		if err := tx.Where("telegram_user_id = ?", user.TelegramUserId).
			Where("telegram_chat_id = ?", user.TelegramChatId).
			First(&u).Error; err != nil {
			return err
		}
		if err := tx.Model(&u).Update("number_of_speeches", u.NumberOfSpeeches+1).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func IncrementVerificationTimes(user *UserInfo) error {
	err := Gdb.Transaction(func(tx *gorm.DB) error {
		var u UserInfo
		if err := tx.Where("telegram_user_id = ?", user.TelegramUserId).
			Where("telegram_chat_id = ?", user.TelegramChatId).
			First(&u).Error; err != nil {
			return err
		}
		if err := tx.Model(&u).Update("verification_times", u.VerificationTimes+1).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func GetUserInfo(user *UserInfo) (*UserInfo, error) {
	var u UserInfo
	err := Gdb.Where("telegram_user_id = ?", user.TelegramUserId).Where("telegram_chat_id = ?", user.TelegramChatId).First(&u).Error
	return &u, err
}
