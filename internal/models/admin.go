package models

type Admin struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Login    string `gorm:"not null;unique;size:255"`
	Password string `gorm:"not null;size:100"`
}


