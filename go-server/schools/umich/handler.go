// Package umich configures an api handler for umich.
package umich

import (
	"context"
	"fmt"
	"log"
	"os"

	alerts "github.com/benkoppe/bear-trak-backend/go-server/alerts/umich"
	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/db"
	dining "github.com/benkoppe/bear-trak-backend/go-server/dining/umich"
	"github.com/benkoppe/bear-trak-backend/go-server/diningusers"
	"github.com/benkoppe/bear-trak-backend/go-server/schools/shared"
	transit "github.com/benkoppe/bear-trak-backend/go-server/transit/umich"
)

type Handler struct {
	DB *db.Queries

	diningCache   dining.Cache
	transitCaches transit.Caches
}

func NewHandler(db *db.Queries) *Handler {
	h := &Handler{
		DB: db,
	}
	h.initCaches()

	return h
}

const APIKeyEnvVar = "BUSTIME_API_KEY"

func (h *Handler) initCaches() {
	// get API key from environment variables
	bustimeAPIKey := os.Getenv(APIKeyEnvVar)
	if bustimeAPIKey == "" {
		log.Fatalf("API key not found in environment variables")
	}
	h.transitCaches = transit.InitCaches(bustimeURL, bustimeAPIKey, gtfsStaticURL)

	h.diningCache = dining.InitCache(eateriesBaseURL)
}

func (h *Handler) GetV1Alerts(ctx context.Context) ([]api.Alert, error) {
	return alerts.Get()
}

func (h *Handler) GetV1Dining(ctx context.Context) ([]api.Eatery, error) {
	return dining.Get(h.diningCache)
}

func (h *Handler) GetV1Gyms(ctx context.Context) ([]api.Gym, error) {
	return nil, fmt.Errorf("umich doesn't implement the gyms feature")
}

func (h *Handler) GetV1TransitRoutes(ctx context.Context) ([]api.BusRoute, error) {
	return transit.GetRoutes(h.transitCaches)
}

func (h *Handler) GetV1TransitVehicles(ctx context.Context) ([]api.Vehicle, error) {
	return transit.GetVehicles(h.transitCaches)
}

func (h *Handler) GetV1Study(ctx context.Context) (*api.StudyData, error) {
	return nil, fmt.Errorf("not implemented")
}

func (h *Handler) GetV1DiningUser(ctx context.Context, params api.GetV1DiningUserParams) (api.GetV1DiningUserRes, error) {
	return diningusers.GetUser(shared.CbordBaseURL, cbordInstitutionID, params)
}

func (h *Handler) PostV1DiningUser(ctx context.Context, params api.PostV1DiningUserParams) (api.PostV1DiningUserRes, error) {
	return diningusers.CreateUser(ctx, shared.CbordBaseURL, cbordInstitutionID, params, h.DB)
}

func (h *Handler) DeleteV1DiningUser(ctx context.Context, params api.DeleteV1DiningUserParams) (api.DeleteV1DiningUserRes, error) {
	return diningusers.DeleteUser(ctx, shared.CbordBaseURL, params, h.DB)
}

func (h *Handler) GetV1DiningUserSession(ctx context.Context, params api.GetV1DiningUserSessionParams) (api.GetV1DiningUserSessionRes, error) {
	return diningusers.RefreshUserToken(ctx, shared.CbordBaseURL, params, h.DB)
}

func (h *Handler) GetV1DiningUserAccounts(ctx context.Context, params api.GetV1DiningUserAccountsParams) (api.GetV1DiningUserAccountsRes, error) {
	return diningusers.GetUserAccounts(shared.CbordBaseURL, cbordInstitutionID, params)
}

func (h *Handler) GetV1DiningUserBarcode(ctx context.Context, params api.GetV1DiningUserBarcodeParams) (api.GetV1DiningUserBarcodeRes, error) {
	return diningusers.GetUserBarcode(shared.CbordBaseURL, params)
}

func (h *Handler) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	return &api.ErrorStatusCode{
		StatusCode: 400,
		Response: api.Error{
			Code:    400,
			Message: err.Error(),
		},
	}
}
