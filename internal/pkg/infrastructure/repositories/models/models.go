package models

import (
	"time"

	"gorm.io/gorm"
)

type WaterConsumption struct {
	gorm.Model
	Device      string
	Consumption int
	Timestamp   time.Time
}
