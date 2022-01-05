package application

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	fiware "github.com/diwise/ngsi-ld-golang/pkg/datamodels/fiware"
	"github.com/go-chi/chi"
	"github.com/matryer/is"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sundsvall/api-smartwater/internal/pkg/infrastructure/repositories/database"
)

func newAppForTesting() (*database.Datastore, zerolog.Logger, *waterConsumptionApp) {
	r := chi.NewRouter()
	log := log.Logger
	return nil, log, newWaterConsumptionApp(r, nil, log, "serviceName")
}

func newTestRequest(is *is.I, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, _ := http.NewRequest(method, ts.URL+path, body)
	resp, _ := http.DefaultClient.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestReceiveWaterConsumption(t *testing.T) {
	is := is.New(t)

	_, _, app := newAppForTesting()
	ts := httptest.NewServer(app.router)
	defer ts.Close()

	id := "waterConsumption01"
	observedAt := time.Now().UTC()

	wco := fiware.NewWaterConsumptionObserved(id)
	wco.WithConsumption(id, 806040.0, observedAt)

	wcoJson, err := json.Marshal(wco)
	if err != nil {
		t.Error("could not marshal wco to json")
	}

	resp, _ := newTestRequest(is, ts, "POST", "/ngsi-ld/v1/entities", bytes.NewBuffer(wcoJson))

	is.Equal(resp.StatusCode, http.StatusCreated) // Check status code

}
