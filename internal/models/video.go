package models

type Video struct {
	ID          uint   `gorm:"primaryKey"`
	Path       string  `gorm:"not null;varchar(255)"`
}