package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Description string `json:"description" gorm:"text;not null;default:null"`
	Priority uint `json:"priority" gorm:"default:3"`
	Completed bool`json:"completed" gorm:"default:false"`
}