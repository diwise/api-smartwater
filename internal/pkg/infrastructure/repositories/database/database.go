package database

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/sundsvall/api-smartwater/internal/pkg/infrastructure/repositories/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Datastore interface {
	StoreWaterConsumption(device string, consumption float64, timestamp time.Time) (*models.WaterConsumption, error)
}

type myDB struct {
	impl *gorm.DB
	log  zerolog.Logger
}

//ConnectorFunc is used to inject a database connection method into NewDatabaseConnection
type ConnectorFunc func() (*gorm.DB, zerolog.Logger, error)

//NewSQLiteConnector opens a connection to a local sqlite database
func NewSQLiteConnector(log zerolog.Logger) ConnectorFunc {
	return func() (*gorm.DB, zerolog.Logger, error) {
		db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})

		if err == nil {
			db.Exec("PRAGMA foreign_keys = ON")
		}

		return db, log, err
	}
}

//NewDatabaseConnection initializes a new connection to the database and wraps it in a Datastore
func NewDatabaseConnection(connect ConnectorFunc) (Datastore, error) {
	impl, log, err := connect()
	if err != nil {
		return nil, err
	}

	db := &myDB{
		impl: impl.Debug(),
		log:  log,
	}

	db.impl.AutoMigrate(
		&models.WaterConsumption{},
	)

	return db, nil
}

func (db *myDB) StoreWaterConsumption(device string, consumption float64, timestamp time.Time) (*models.WaterConsumption, error) {
	wco := models.WaterConsumption{
		Device:      device,
		Consumption: consumption,
		Timestamp:   timestamp,
	}

	result := db.impl.Create(&wco)
	if result.Error != nil {
		return nil, result.Error
	}

	return &wco, nil
}
