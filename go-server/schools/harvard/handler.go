package harvard

import (
	"context"
	"fmt"
	"log"
	"os"

	alerts "github.com/benkoppe/bear-trak-backend/go-server/alerts/harvard"
	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/db"
	dining "github.com/benkoppe/bear-trak-backend/go-server/dining/harvard"
)

type Handler struct {
	DB *db.Queries

	diningCache dining.Cache
}

func NewHandler(db *db.Queries) *Handler {
	h := &Handler{
		DB: db,
	}
	h.initCaches()

	return h
}

const DINING_API_KEY_ENV_VAR = "HARVARD_DINING_API_KEY"

func (h *Handler) initCaches() {
	// get API key from environment variables
	diningApiKey := os.Getenv(DINING_API_KEY_ENV_VAR)
	if diningApiKey == "" {
		log.Fatalf("Dining API key not found in environment variables")
	}

	h.diningCache = dining.InitCache(eateriesBaseUrl, diningApiKey)
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
	return nil, fmt.Errorf("not implemented")
}

func (h *Handler) GetV1TransitVehicles(ctx context.Context) ([]api.Vehicle, error) {
	return nil, fmt.Errorf("not implemented")
}

func (h *Handler) GetV1Study(ctx context.Context) (*api.StudyData, error) {
	return nil, fmt.Errorf("not implemented")
}

func (h *Handler) GetV1DiningUser(ctx context.Context, params api.GetV1DiningUserParams) (api.GetV1DiningUserRes, error) {
	return nil, fmt.Errorf("not implemented")
}

func (h *Handler) PostV1DiningUser(ctx context.Context, params api.PostV1DiningUserParams) (api.PostV1DiningUserRes, error) {
	return nil, fmt.Errorf("not implemented")
}

func (h *Handler) DeleteV1DiningUser(ctx context.Context, params api.DeleteV1DiningUserParams) (api.DeleteV1DiningUserRes, error) {
	return nil, fmt.Errorf("not implemented")
}

func (h *Handler) GetV1DiningUserSession(ctx context.Context, params api.GetV1DiningUserSessionParams) (api.GetV1DiningUserSessionRes, error) {
	return nil, fmt.Errorf("not implemented")
}

func (h *Handler) GetV1DiningUserAccounts(ctx context.Context, params api.GetV1DiningUserAccountsParams) (api.GetV1DiningUserAccountsRes, error) {
	return nil, fmt.Errorf("not implemented")
}

func (h *Handler) GetV1DiningUserBarcode(ctx context.Context, params api.GetV1DiningUserBarcodeParams) (api.GetV1DiningUserBarcodeRes, error) {
	return nil, fmt.Errorf("not implemented")
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
