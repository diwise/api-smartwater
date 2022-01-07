package application

import (
	"testing"
	"time"

	"github.com/diwise/api-smartwater/internal/pkg/infrastructure/repositories/database"
	"github.com/diwise/api-smartwater/internal/pkg/infrastructure/repositories/models"
	"github.com/matryer/is"
	"github.com/rs/zerolog/log"
)

func newAppForTesting() (*database.DatastoreMock, Application) {
	db := &database.DatastoreMock{
		StoreWaterConsumptionFunc: func(device string, consumption float64, timestamp time.Time) (*models.WaterConsumption, error) {
			return nil, nil
		},
	}
	log := log.Logger

	return db, newWaterConsumptionApp(db, log, "serviceName")
}
func TestUpdateWaterConsumption(t *testing.T) {
	is := is.New(t)
	db, app := newAppForTesting()

	err := app.UpdateWaterConsumption("device", 172.0, time.Now().UTC())
	is.NoErr(err)                                     // Check error
	is.Equal(len(db.StoreWaterConsumptionCalls()), 1) // StoreWaterConsumption should have been called once
}
