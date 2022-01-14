package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/diwise/api-smartwater/internal/pkg/application"
	"github.com/diwise/api-smartwater/internal/pkg/infrastructure/repositories/database"
	"github.com/diwise/ngsi-ld-golang/pkg/datamodels/fiware"
	"github.com/go-chi/chi"
	"github.com/matryer/is"
	"github.com/rs/zerolog/log"
)

func newRouterForTesting() chi.Router {
	r := chi.NewRouter()
	log := log.Logger
	db, _ := database.NewDatabaseConnection(database.NewSQLiteConnector(log))
	app := application.NewApplication(db, log, "serviceName")

	RegisterHandlers(r, app, log)

	return r
}

func newTestRequest(is *is.I, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, _ := http.NewRequest(method, ts.URL+path, body)
	resp, _ := http.DefaultClient.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestWaterConsumptionHandler(t *testing.T) {
	is := is.New(t)

	r := newRouterForTesting()

	ts := httptest.NewServer(r)
	defer ts.Close()

	id := "waterConsumption01"
	observedAt := time.Now().UTC()

	wco := fiware.NewWaterConsumptionObserved(id)
	wco.WithConsumption(id, 806040.0, observedAt)

	wcoJson, err := json.Marshal(wco)
	is.NoErr(err) // could not marshal wco to json

	resp, _ := newTestRequest(is, ts, "POST", "/ngsi-ld/v1/entities", bytes.NewBuffer(wcoJson))
	is.Equal(resp.StatusCode, http.StatusCreated) // Check status code
}
