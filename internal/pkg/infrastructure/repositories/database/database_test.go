package database

import (
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/rs/zerolog/log"
)

func TestThatStoreWaterConsumptionDoesNotReturnError(t *testing.T) {
	is, db := setupTest(t)

	_, err := db.StoreWaterConsumption("entityId567", "deviceId", 376.0, time.Now().UTC())
	is.NoErr(err) // error when storing new water consumption
}

func TestThatStoreWaterConsumptionFailsOnIdenticalWCOIDs(t *testing.T) {
	is, db := setupTest(t)

	db.StoreWaterConsumption("entityId567", "deviceId", 376.0, time.Now().UTC())
	_, err := db.StoreWaterConsumption("entityId567", "deviceId", 376.0, time.Now().UTC())
	is.Equal(err.Error(), "UNIQUE constraint failed: water_consumptions.wco_id") // error when storing new water consumption
}

func TestRetrievingAllStoredWaterConsumptionReadings(t *testing.T) {
	is, db := setupTest(t)

	time1 := time.Now().UTC()
	time2 := time.Now().UTC().Add(2 * time.Hour)
	time3 := time.Now().UTC().Add(3 * time.Hour)

	db.StoreWaterConsumption("entityId123", "deviceId", 176.0, time1)
	db.StoreWaterConsumption("entityId456", "deviceId", 1799.0, time2)
	db.StoreWaterConsumption("entityId789", "deviceId", 17.0, time3)

	result, err := db.GetWaterConsumptions("", time.Time{}, time.Time{}, 0)

	is.NoErr(err)
	is.Equal(len(result), 3) // should equal 3
}

func TestRetrievingWaterConsumptionsWithinTimespan(t *testing.T) {
	is, db := setupTest(t)

	time1 := time.Now().UTC().Add(-3 * time.Hour)
	time2 := time.Now().UTC().Add(-2 * time.Hour)
	time3 := time.Now().UTC()

	db.StoreWaterConsumption("entityId123", "deviceId", 176.0, time1)
	db.StoreWaterConsumption("entityId456", "deviceId", 1799.0, time2)
	db.StoreWaterConsumption("entityId789", "deviceId", 17.0, time3)

	result, err := db.GetWaterConsumptions("", time1, time3, 0)

	is.NoErr(err)
	is.Equal(len(result), 2)             // should equal 2
	is.Equal(result[0].Timestamp, time2) // should be most recent entry by timestamp
}

func TestThatItIsNotPossibleToSaveWaterConsumptionFromTheSameDeviceAtTheSameTimestamp(t *testing.T) {
	is, db := setupTest(t)

	time := time.Now().UTC()

	db.StoreWaterConsumption("entityId123", "deviceId", 176.0, time)
	_, err := db.StoreWaterConsumption("entityId456", "deviceId", 1799.0, time)

	is.True(err != nil)                                                                                        // error should not be nil
	is.Equal(err.Error(), "UNIQUE constraint failed: water_consumptions.device, water_consumptions.timestamp") // error messages should match
}

func setupTest(t *testing.T) (*is.I, Datastore) {
	is := is.New(t)
	log := log.Logger
	db, err := NewDatabaseConnection(NewSQLiteConnector(log))
	is.NoErr(err) // error when creating new database connection

	return is, db
}
