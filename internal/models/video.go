package models

type Video struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Path       string  `gorm:"not null;varchar(255)" json:"path"`
}