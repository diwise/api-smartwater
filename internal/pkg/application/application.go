package application

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/sundsvall/api-smartwater/internal/pkg/infrastructure/repositories/database"
)

type Application interface {
	UpdateWaterConsumption(device string, consumption float64, timestamp time.Time) error
}

type waterConsumptionApp struct {
	db  database.Datastore
	log zerolog.Logger
}

func NewApplication(db database.Datastore, log zerolog.Logger, serviceName string) Application {
	return newWaterConsumptionApp(db, log, serviceName)
}

func newWaterConsumptionApp(db database.Datastore, log zerolog.Logger, serviceName string) *waterConsumptionApp {
	w := &waterConsumptionApp{
		db:  db,
		log: log,
	}

	return w
}

func (w *waterConsumptionApp) UpdateWaterConsumption(device string, consumption float64, timestamp time.Time) error {
	_, err := w.db.StoreWaterConsumption(device, consumption, timestamp)
	if err != nil {
		return err
	}

	return nil
}
