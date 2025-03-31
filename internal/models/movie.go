package models

type Movie struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null;varchar(255)"`
	Description string `gorm:"not null;varchar(10000)"`
	VideoID     uint  `gorm:"index"`
	Video       Video `gorm:"foreignKey:VideoID"`
}