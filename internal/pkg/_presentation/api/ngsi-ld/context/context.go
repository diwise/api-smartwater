package context

import (
	"errors"
	"strings"
	"time"

	"github.com/diwise/api-smartwater/internal/pkg/application"
	"github.com/diwise/api-smartwater/internal/pkg/infrastructure/repositories/models"
	"github.com/diwise/ngsi-ld-golang/pkg/datamodels/fiware"
	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld"
	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type contextSource struct {
	app application.Application
	log zerolog.Logger
}

//CreateSource instantiates and returns a Fiware ContextSource that wraps the provided db interface
func CreateSource(app application.Application) ngsi.ContextSource {
	log := log.Logger

	return &contextSource{
		app: app,
		log: log,
	}
}

func (cs contextSource) CreateEntity(typeName, entityID string, req ngsi.Request) error {

	cs.log.With().Str("entityID", entityID).Str("entityType", typeName).Logger()

	if typeName != "WaterConsumptionObserved" {
		errorMessage := "entity type not supported"
		cs.log.Error().Msg(errorMessage)
		return errors.New(errorMessage)
	}

	wco := &fiware.WaterConsumptionObserved{}
	err := req.DecodeBodyInto(wco)
	if err != nil {
		return err // TODO: fix later
	}

	observedAt, err := time.Parse(time.RFC3339, wco.WaterConsumption.ObservedAt)
	if err != nil {
		return err // TODO: fix later
	}

	err = cs.app.UpdateWaterConsumption(wco.ID, wco.WaterConsumption.Value, observedAt)
	if err != nil {
		return err // TODO: fix later
	}

	return err
}

func (cs contextSource) GetEntities(query ngsi.Query, callback ngsi.QueryEntitiesCallback) error {
	var err error

	if query == nil {
		return errors.New("GetEntities: query may not be nil")
	}

	waterconsumptions, err := getWaterConsumptions(cs.app, query)
	if err != nil {
		return err // TODO: fix later
	}

	for _, w := range waterconsumptions {
		entity := fiware.NewWaterConsumptionObserved(w.Device)
		entity.WaterConsumption = &fiware.WCONumberProperty{
			NumberProperty: *types.NewNumberProperty(w.Consumption),
			ObservedAt:     w.Timestamp.Format(time.RFC3339),
		}
		err = callback(entity)
		if err != nil {
			break
		}
	}

	return err
}

func (cs contextSource) GetProvidedTypeFromID(entityID string) (string, error) {
	return "", errors.New("not implemented")
}

func (cs contextSource) ProvidesAttribute(attributeName string) bool {
	return attributeName == "waterconsumption"
}

func (cs contextSource) ProvidesEntitiesWithMatchingID(entityID string) bool {
	return strings.HasPrefix(entityID, "urn:ngsi-ld:WaterConsumptionObserved:")
}

func (cs contextSource) ProvidesType(typeName string) bool {
	return typeName == "WaterConsumptionObserved"
}

func (cs contextSource) RetrieveEntity(entityID string, request ngsi.Request) (ngsi.Entity, error) {
	return nil, errors.New("retrieve entity not implemented")
}

func (cs contextSource) UpdateEntityAttributes(entityID string, req ngsi.Request) error {
	return errors.New("UpdateEntityAttributes is not supported by this service")
}

func getWaterConsumptions(app application.Application, query ngsi.Query) ([]models.WaterConsumption, error) {
	deviceId := ""
	if query.Device() != "" {
		deviceId = strings.TrimPrefix(query.Device(), fiware.WaterConsumptionObservedIDPrefix)
	}

	from := time.Time{}
	to := time.Time{}
	if query.IsTemporalQuery() {
		from, to = query.Temporal().TimeSpan()
	}

	limit := query.PaginationLimit()

	wcos, err := app.RetrieveWaterConsumptions(deviceId, from, to, limit)
	if err != nil {
		return nil, err
	}

	return wcos, nil
}
