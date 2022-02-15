package context

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/diwise/api-smartwater/internal/pkg/application"
	"github.com/diwise/api-smartwater/internal/pkg/infrastructure/repositories/models"
	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld"
	"github.com/matryer/is"
	"github.com/rs/zerolog/log"
)

func TestUpdateWaterConsumption(t *testing.T) {

	req, _ := http.NewRequest("POST", "/ngsi-ld/v1/entities", bytes.NewBuffer([]byte(wcoJson)))
	w := httptest.NewRecorder()

	is, app, ctxReg := testSetup(t)

	ngsi.NewCreateEntityHandler(ctxReg).ServeHTTP(w, req)

	is.Equal(w.Code, http.StatusCreated)
	is.Equal(len(app.UpdateWaterConsumptionCalls()), 1)
	is.Equal(app.UpdateWaterConsumptionCalls()[0].Consumption, float64(191051))
	is.Equal(app.UpdateWaterConsumptionCalls()[0].Device, "testDevice01")
}

func TestRetrieveEntities(t *testing.T) {
	is, app, ctxReg := testSetup(t)

	req, _ := http.NewRequest("GET", "/ngsi-ld/v1/entities?type=WaterConsumptionObserved", nil)
	w := httptest.NewRecorder()

	ngsi.NewQueryEntitiesHandler(ctxReg).ServeHTTP(w, req)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(len(app.RetrieveWaterConsumptionsCalls()), 1)
}

func testSetup(t *testing.T) (*is.I, *application.ApplicationMock, ngsi.ContextRegistry) {
	is := is.New(t)

	log := log.Logger
	app := &application.ApplicationMock{
		RetrieveWaterConsumptionsFunc: func(deviceId string, from, to time.Time, limit uint64) ([]models.WaterConsumption, error) {
			return nil, nil
		},
		UpdateWaterConsumptionFunc: func(device string, consumption float64, timestamp time.Time) error {
			return nil
		},
	}

	ctxReg := ngsi.NewContextRegistry()
	ctxSource := CreateSource(app, log)
	ctxReg.Register(ctxSource)

	return is, app, ctxReg
}

const wcoJson string = `{
    "id": "urn:ngsi-ld:WaterConsumptionObserved:Consumer01",
    "type": "WaterConsumptionObserved",
    "waterConsumption": {
        "type": "Property",
        "observedBy": {
            "type": "Relationship",
            "object": "urn:ngsi-ld:Device:testDevice01"
        },
        "value": 191051,
        "observedAt": "2021-05-23T23:14:16.000Z",
        "unitCode": "LTR"
    },
    "@context": [
        "https://raw.githubusercontent.com/easy-global-market/ngsild-api-data-models/master/WaterSmartMeter/jsonld-contexts/waterSmartMeter-compound.jsonld"
    ]
}`
