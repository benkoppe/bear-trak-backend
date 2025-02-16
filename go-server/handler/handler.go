package handler

import (
	"context"
	"fmt"

	"github.com/benkoppe/bear-trak-backend/go-server/alerts"
	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/dining"
	dining_users "github.com/benkoppe/bear-trak-backend/go-server/dining-users"
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

func (bs *BackendService) PostV1DiningUser(ctx context.Context, req api.OptPostV1DiningUserReq) (api.PostV1DiningUserRes, error) {
	if !req.IsSet() {
		return nil, fmt.Errorf("missing required fields")
	}

	return dining_users.CreateUser(cbordBaseUrl, req.Value)
}

func (bs *BackendService) DeleteV1DiningUser(ctx context.Context, req api.OptDiningUserSession) (api.DeleteV1DiningUserRes, error) {
	if !req.IsSet() {
		return nil, fmt.Errorf("missing required fields")
	}

	return dining_users.DeleteUser(cbordBaseUrl, req.Value)
}

func (bs *BackendService) GetV1DiningUserSession(ctx context.Context, req api.OptDiningUserDevice) (api.GetV1DiningUserSessionRes, error) {
	if !req.IsSet() {
		return nil, fmt.Errorf("missing required fields")
	}

	return dining_users.RefreshUserToken(cbordBaseUrl, req.Value)
}

func (bs *BackendService) GetV1DiningUserAccounts(ctx context.Context, req api.OptDiningUserSession) (api.GetV1DiningUserAccountsRes, error) {
	if !req.IsSet() {
		return nil, fmt.Errorf("missing required fields")
	}

	return dining_users.GetUserAccounts(cbordBaseUrl, req.Value)
}

func (bs *BackendService) GetV1DiningUserBarcode(ctx context.Context, req api.OptDiningUserSession) (api.GetV1DiningUserBarcodeRes, error) {
	if !req.IsSet() {
		return nil, fmt.Errorf("missing required fields")
	}

	return dining_users.GetUserBarcode(cbordBaseUrl, req.Value)
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
