package models

import (
	"time"

	"gorm.io/gorm"
)

type WaterConsumption struct {
	gorm.Model
	WCOID       string `gorm:"index:idx_wco_id,unique"`
	Device      string `gorm:"index;index:device_at_time,unique"`
	Consumption float64
	Timestamp   time.Time `gorm:"index:device_at_time,unique"`
}
