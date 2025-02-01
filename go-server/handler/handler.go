package handler

import (
	"context"

	"github.com/benkoppe/bear-trak-backend/go-server/alerts"
	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/dining"
	"github.com/benkoppe/bear-trak-backend/go-server/gyms"
	"github.com/benkoppe/bear-trak-backend/go-server/transit"
)

type BackendService struct{}

func (bs *BackendService) GetV1Alerts(ctx context.Context) ([]api.Alert, error) {
	return alerts.Get()
}

func (bs *BackendService) GetV1Dining(ctx context.Context) ([]api.Eatery, error) {
	return dining.Get(eateriesUrl)
}

func (bs *BackendService) GetV1Gyms(ctx context.Context) ([]api.Gym, error) {
	return gyms.Get(gymCapacitiesUrl)
}

func (bs *BackendService) GetV1TransitRoutes(ctx context.Context) ([]api.BusRoute, error) {
	return transit.GetRoutes(availtecUrl, gtfsStaticUrl)
}

func (bs *BackendService) GetV1TransitVehicles(ctx context.Context) ([]api.Vehicle, error) {
	return transit.GetVehicles(availtecUrl)
}

func (bs *BackendService) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	return &api.ErrorStatusCode{
		StatusCode: 400,
		Response: api.Error{
			Code:    400,
			Message: err.Error(),
		},
	}
}
