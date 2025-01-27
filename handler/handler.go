package handler

import (
	"context"

	"github.com/benkoppe/bear-trak-backend/alerts"
	backend "github.com/benkoppe/bear-trak-backend/backend"
	"github.com/benkoppe/bear-trak-backend/dining"
	"github.com/benkoppe/bear-trak-backend/gyms"
	"github.com/benkoppe/bear-trak-backend/transit"
)

type BackendService struct{}

func (bs *BackendService) GetV1Alerts(ctx context.Context) ([]backend.Alert, error) {
	return alerts.Get()
}

func (bs *BackendService) GetV1Dining(ctx context.Context) ([]backend.Eatery, error) {
	return dining.Get(eateriesUrl)
}

func (bs *BackendService) GetV1Gyms(ctx context.Context) ([]backend.Gym, error) {
	return gyms.Get(gymCapacitiesUrl)
}

func (bs *BackendService) GetV1TransitRoutes(ctx context.Context) ([]backend.BusRoute, error) {
	return transit.GetRoutes(availtecUrl, gtfsStaticUrl)
}

func (bs *BackendService) GetV1TransitVehicles(ctx context.Context) ([]backend.Vehicle, error) {
	return transit.GetVehicles(availtecUrl)
}

func (bs *BackendService) NewError(ctx context.Context, err error) *backend.ErrorStatusCode {
	return &backend.ErrorStatusCode{
		StatusCode: 400,
		Response: backend.Error{
			Code:    400,
			Message: err.Error(),
		},
	}
}
