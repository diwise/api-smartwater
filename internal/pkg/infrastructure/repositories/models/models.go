package models

import (
	"time"

	"gorm.io/gorm"
)

type WaterConsumption struct {
	gorm.Model
	WCOID       string
	Device      string
	Consumption float64
	Timestamp   time.Time
}
