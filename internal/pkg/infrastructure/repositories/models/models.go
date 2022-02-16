package models

import (
	"time"

	"gorm.io/gorm"
)

type WaterConsumption struct {
	gorm.Model
	WCOID       string `gorm:"index:idx_wco_id,unique"`
	Device      string
	Consumption float64
	Timestamp   time.Time
}
