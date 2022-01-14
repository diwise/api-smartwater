package context

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/diwise/api-smartwater/internal/pkg/application"
	"github.com/diwise/api-smartwater/internal/pkg/infrastructure/repositories/database"
	"github.com/diwise/api-smartwater/internal/pkg/infrastructure/repositories/models"
	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld"
	"github.com/matryer/is"
	"github.com/rs/zerolog/log"
)

func TestCreateEntityWorks(t *testing.T) {

	req, _ := http.NewRequest("POST", "http://localhost:8090/ngsi-ld/v1/entities", bytes.NewBuffer([]byte(wcoJson)))
	w := httptest.NewRecorder()

	is, ctxReg := testSetup(t)

	ngsi.NewCreateEntityHandler(ctxReg).ServeHTTP(w, req)

	fmt.Println(w.Body.String())

	is.Equal(w.Code, http.StatusCreated)
}

func TestRetrieveEntities(t *testing.T) {
	is, ctxReg := testSetup(t)

	req, _ := http.NewRequest("GET", "http://localhost:8090/ngsi-ld/v1/entities?type=WaterConsumptionObserved", nil)
	w := httptest.NewRecorder()

	ngsi.NewQueryEntitiesHandler(ctxReg).ServeHTTP(w, req)

	fmt.Println(w.Body.String())

	is.Equal(w.Code, http.StatusOK)
}

func testSetup(t *testing.T) (*is.I, ngsi.ContextRegistry) {
	is := is.New(t)

	wcos := []models.WaterConsumption{
		{
			Device:      "device",
			Consumption: 172.0,
			Timestamp:   time.Now().UTC(),
		},
	}

	db := &database.DatastoreMock{
		GetWaterConsumptionsFunc: func(deviceId string, from time.Time, to time.Time, limit uint64) ([]models.WaterConsumption, error) {
			return wcos, nil
		},
		StoreWaterConsumptionFunc: func(device string, consumption float64, timestamp time.Time) (*models.WaterConsumption, error) {
			return nil, nil
		},
	}
	log := log.Logger
	app := application.NewApplication(db, log, "api-smartwater")

	ctxReg := ngsi.NewContextRegistry()
	ctxSource := CreateSource(app)
	ctxReg.Register(ctxSource)

	return is, ctxReg
}

const wcoJson string = `{
    "id": "urn:ngsi-ld:WaterConsumptionObserved:Consumer01",
    "type": "WaterConsumptionObserved",
    "waterConsumption": {
        "type": "Property",
        "observedBy": {
            "type": "Relationship",
            "object": "urn:ngsi-ld:Device:01"
        },
        "value": 191051,
        "observedAt": "2021-05-23T23:14:16.000Z",
        "unitCode": "LTR"
    },
    "@context": [
        "https://raw.githubusercontent.com/easy-global-market/ngsild-api-data-models/master/WaterSmartMeter/jsonld-contexts/waterSmartMeter-compound.jsonld"
    ]
}`
