package models

import "gorm.io/gorm"

type TodoAnalyticsDto struct {
	gorm.Model
	Count int32
}