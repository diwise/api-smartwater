package database

import (
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/rs/zerolog/log"
)

func TestThatStoreWaterConsumptionDoesNotReturnError(t *testing.T) {
	is, db := setupTest(t)

	_, err := db.StoreWaterConsumption("deviceId", 176.0, time.Now().UTC())

	is.NoErr(err) // error when storing new water consumption
}

func setupTest(t *testing.T) (*is.I, Datastore) {
	is := is.New(t)
	log := log.Logger
	db, err := NewDatabaseConnection(NewSQLiteConnector(log))
	is.NoErr(err) // error when creating new database connection

	return is, db
}
