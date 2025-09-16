package cornell

import (
	"context"
	"fmt"

	alerts "github.com/benkoppe/bear-trak-backend/go-server/alerts/cornell"
	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/db"
	dining "github.com/benkoppe/bear-trak-backend/go-server/dining/cornell"
	"github.com/benkoppe/bear-trak-backend/go-server/diningusers"
	gyms "github.com/benkoppe/bear-trak-backend/go-server/gyms/cornell"
	"github.com/benkoppe/bear-trak-backend/go-server/schools/cornell/externalmap"
	"github.com/benkoppe/bear-trak-backend/go-server/schools/shared"
	study "github.com/benkoppe/bear-trak-backend/go-server/study/cornell"
	transit "github.com/benkoppe/bear-trak-backend/go-server/transit/cornell"
)

type Handler struct {
	DB *db.Queries

	diningCache   dining.Cache
	gymsCaches    gyms.Caches
	transitCaches transit.Caches
	studyCache    study.Cache
	mapCache      externalmap.Cache
}

func NewHandler(db *db.Queries) *Handler {
	h := &Handler{
		DB: db,
	}
	h.initCaches(db)

	return h
}

func (h *Handler) initCaches(db *db.Queries) {
	h.diningCache = dining.InitCache(eateriesURL)
	h.gymsCaches = gyms.InitCaches(gymCapacitiesURL, gymHoursURL, gymPredictionsURL, db)
	h.transitCaches = transit.InitCaches(availtecURL, gtfsStaticURL)
	h.studyCache = study.InitCache(librariesURL)
	h.mapCache = externalmap.InitCache(mapOverlaysURL)
}

func (h *Handler) GetV1Alerts(ctx context.Context) ([]api.Alert, error) {
	return alerts.Get()
}

func (h *Handler) GetV1Dining(ctx context.Context) ([]api.Eatery, error) {
	return dining.Get(h.diningCache)
}

func (h *Handler) GetV1Gyms(ctx context.Context) ([]api.Gym, error) {
	return gyms.Get(h.gymsCaches)
}

func (h *Handler) GetV1GymCapacities(ctx context.Context) ([]api.GymCapacityData, error) {
	return gyms.GetCapacityPoints(h.gymsCaches)
}

func (h *Handler) GetV1GymCapacityPredictions(ctx context.Context) ([]api.GymCapacityPredictions, error) {
	return gyms.GetCapacityPredictionPoints(h.gymsCaches)
}

func (h *Handler) GetV1TransitRoutes(ctx context.Context) ([]api.BusRoute, error) {
	return transit.GetRoutes(h.transitCaches)
}

func (h *Handler) GetV1TransitVehicles(ctx context.Context) ([]api.Vehicle, error) {
	return transit.GetVehicles(h.transitCaches)
}

func (h *Handler) GetV1Study(ctx context.Context) (*api.StudyData, error) {
	return study.Get(h.studyCache, h.mapCache)
}

func (h *Handler) GetV1Events(ctx context.Context) ([]api.Event, error) {
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
