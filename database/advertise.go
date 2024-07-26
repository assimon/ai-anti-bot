package database

import "github.com/golang-module/carbon/v2"

type Advertise struct {
	BaseModel
	Title          string          `gorm:"column:title" json:"title"`
	Url            string          `gorm:"column:url" json:"url"`
	Sort           int             `gorm:"column:sort" json:"sort"`
	ValidityPeriod carbon.DateTime `gorm:"column:validity_period" json:"validity_period"`
}

func (Advertise) TableName() string {
	return "advertise"
}

func AddAdvertise(advertise Advertise) (err error) {
	return Gdb.Model(&advertise).Create(&advertise).Error
}

func AllAdvertise() (advertises []Advertise, err error) {
	err = Gdb.Model(&advertises).Find(&advertises).Error
	return
}

func GetEfficientAdvertise() (advertises []Advertise, err error) {
	err = Gdb.Model(&advertises).Where("validity_period > ?", carbon.Now().ToDateTimeString()).Order("sort desc").Find(&advertises).Error
	return
}

func DeleteAdvertise(id int64) (err error) {
	return Gdb.Where("id = ?", id).Delete(&Advertise{}).Error
}
