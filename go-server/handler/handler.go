package handler

import (
	"context"

	"github.com/benkoppe/bear-trak-backend/go-server/alerts"
	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/db"
	"github.com/benkoppe/bear-trak-backend/go-server/dining"
	dining_users "github.com/benkoppe/bear-trak-backend/go-server/dining-users"
	"github.com/benkoppe/bear-trak-backend/go-server/gyms"
	"github.com/benkoppe/bear-trak-backend/go-server/transit"
)

type BackendService struct {
	DB *db.Queries
}

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

func (bs *BackendService) GetV1DiningUser(ctx context.Context, params api.GetV1DiningUserParams) (api.GetV1DiningUserRes, error) {
	return dining_users.GetUser(cbordBaseUrl, params)
}

func (bs *BackendService) PostV1DiningUser(ctx context.Context, params api.PostV1DiningUserParams) (api.PostV1DiningUserRes, error) {
	return dining_users.CreateUser(ctx, cbordBaseUrl, params, bs.DB)
}

func (bs *BackendService) DeleteV1DiningUser(ctx context.Context, params api.DeleteV1DiningUserParams) (api.DeleteV1DiningUserRes, error) {
	return dining_users.DeleteUser(ctx, cbordBaseUrl, params, bs.DB)
}

func (bs *BackendService) GetV1DiningUserSession(ctx context.Context, params api.GetV1DiningUserSessionParams) (api.GetV1DiningUserSessionRes, error) {
	return dining_users.RefreshUserToken(ctx, cbordBaseUrl, params, bs.DB)
}

func (bs *BackendService) GetV1DiningUserAccounts(ctx context.Context, params api.GetV1DiningUserAccountsParams) (api.GetV1DiningUserAccountsRes, error) {
	return dining_users.GetUserAccounts(cbordBaseUrl, params)
}

func (bs *BackendService) GetV1DiningUserBarcode(ctx context.Context, params api.GetV1DiningUserBarcodeParams) (api.GetV1DiningUserBarcodeRes, error) {
	return dining_users.GetUserBarcode(cbordBaseUrl, params)
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
