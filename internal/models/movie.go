package models

type Movie struct {
	ID          uint   `gorm:"primaryKey" json:"id"` 
	Title       string `gorm:"not null;varchar(255)" json:"title"`
	Description string `gorm:"not null;varchar(10000)" json:"description"`
	VideoID     uint  `gorm:"index;unique" json:"video_id"`
	Video       Video `gorm:"foreignKey:VideoID" json:"video"`
}