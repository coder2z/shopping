package models

import "github.com/jinzhu/gorm"

type Commodity struct {
	gorm.Model
	Name      string `json:"name" gorm:"type:varchar(100);not null"`
	Link      string `json:"link" gorm:"type:varchar(100);not null"`
	Price     string `json:"price" gorm:"not null"`
	Stock     int    `json:"stock" gorm:"not null"`
	StartTime int64  `json:"startTime" gorm:"not null"`
}
