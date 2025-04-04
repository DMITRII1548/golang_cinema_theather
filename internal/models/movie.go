package models

type Movie struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `gorm:"not null;size:255" json:"title" validate:"required,min=2,max=255"`
	Description string `gorm:"not null;size:10000" json:"description" validate:"required,min=100,max=10000"`
	VideoID     uint   `gorm:"index;unique" json:"video_id" validate:"required"`
	Video       Video  `gorm:"foreignKey:VideoID" json:"video"`
	Thumbnail   string `gorm:"not null;size:255" json:"thumbnail"`
	Preview     string `gorm:"not null;size:255" json:"preview"`
}
