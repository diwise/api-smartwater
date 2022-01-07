package models

import (
	"time"

	"gorm.io/gorm"
)

type WaterConsumption struct {
	gorm.Model
	Device      string
	Consumption float64
	Timestamp   time.Time
}
