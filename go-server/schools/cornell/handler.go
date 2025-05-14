package cornell

import (
	"context"
	"log"
	"os"

	alerts "github.com/benkoppe/bear-trak-backend/go-server/alerts/cornell"
	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/db"
	dining_users "github.com/benkoppe/bear-trak-backend/go-server/dining-users"
	dining "github.com/benkoppe/bear-trak-backend/go-server/dining/cornell"
	dining_email "github.com/benkoppe/bear-trak-backend/go-server/dining/cornell/email"
	gyms "github.com/benkoppe/bear-trak-backend/go-server/gyms/cornell"
	"github.com/benkoppe/bear-trak-backend/go-server/schools/cornell/external_map"
	"github.com/benkoppe/bear-trak-backend/go-server/schools/shared"
	study "github.com/benkoppe/bear-trak-backend/go-server/study/cornell"
	transit "github.com/benkoppe/bear-trak-backend/go-server/transit/cornell"
)

type Handler struct {
	DB *db.Queries

	diningCache      dining.Cache
	houseDinnerCache dining_email.Cache
	gymsCaches       gyms.Caches
	transitCaches    transit.Caches
	studyCache       study.Cache
	mapCache         external_map.Cache
}

func NewHandler(db *db.Queries, houseDinnerCache dining_email.Cache) *Handler {
	h := &Handler{
		DB:               db,
		houseDinnerCache: houseDinnerCache,
	}
	h.initCaches()

	return h
}

const (
	EMAIL_PASSWORD_ENV_VAR     = "EMAIL_PASSWORD"
	MISTRAL_API_KEY_ENV_VAR    = "MISTRAL_API_KEY"
	OPENROUTER_API_KEY_ENV_VAR = "OPENROUTER_API_KEY"
	OPENROUTER_MODEL_ENV_VAR   = "OPENROUTER_MODEL"
)

func InitHouseDinnerCache() dining_email.Cache {
	emailPassword := os.Getenv(EMAIL_PASSWORD_ENV_VAR)
	if emailPassword == "" {
		log.Fatalf("Email Password key " + EMAIL_PASSWORD_ENV_VAR + " not found in environment variables")
	}
	mistralApiKey := os.Getenv(MISTRAL_API_KEY_ENV_VAR)
	if mistralApiKey == "" {
		log.Fatalf("Mistral API key " + MISTRAL_API_KEY_ENV_VAR + " not found in environment variables")
	}
	openrouterApiKey := os.Getenv(OPENROUTER_API_KEY_ENV_VAR)
	if openrouterApiKey == "" {
		log.Fatalf("Openrouter API key " + OPENROUTER_API_KEY_ENV_VAR + " not found in environment variables")
	}
	openrouterModel := os.Getenv(OPENROUTER_MODEL_ENV_VAR)
	if openrouterModel == "" {
		openrouterModel = "google/gemini-2.0-flash-001"
	}

	return dining_email.InitCache(emailPassword, mistralApiKey, openrouterApiKey, openrouterModel)
}

func (h *Handler) initCaches() {
	h.diningCache = dining.InitCache(eateriesUrl)
	h.gymsCaches = gyms.InitCaches(gymCapacitiesUrl, gymHoursUrl)
	h.transitCaches = transit.InitCaches(availtecUrl, gtfsStaticUrl)
	h.studyCache = study.InitCache(librariesUrl)
	h.mapCache = external_map.InitCache(mapOverlaysUrl)
}

func (h *Handler) GetV1Alerts(ctx context.Context) ([]api.Alert, error) {
	return alerts.Get()
}

func (h *Handler) GetV1Dining(ctx context.Context) ([]api.Eatery, error) {
	return dining.Get(h.diningCache, h.houseDinnerCache)
}

func (h *Handler) GetV1Gyms(ctx context.Context) ([]api.Gym, error) {
	return gyms.Get(h.gymsCaches)
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

func (h *Handler) GetV1DiningUser(ctx context.Context, params api.GetV1DiningUserParams) (api.GetV1DiningUserRes, error) {
	return dining_users.GetUser(shared.CbordBaseUrl, cbordInstitutionId, params)
}

func (h *Handler) PostV1DiningUser(ctx context.Context, params api.PostV1DiningUserParams) (api.PostV1DiningUserRes, error) {
	return dining_users.CreateUser(ctx, shared.CbordBaseUrl, cbordInstitutionId, params, h.DB)
}

func (h *Handler) DeleteV1DiningUser(ctx context.Context, params api.DeleteV1DiningUserParams) (api.DeleteV1DiningUserRes, error) {
	return dining_users.DeleteUser(ctx, shared.CbordBaseUrl, params, h.DB)
}

func (h *Handler) GetV1DiningUserSession(ctx context.Context, params api.GetV1DiningUserSessionParams) (api.GetV1DiningUserSessionRes, error) {
	return dining_users.RefreshUserToken(ctx, shared.CbordBaseUrl, params, h.DB)
}

func (h *Handler) GetV1DiningUserAccounts(ctx context.Context, params api.GetV1DiningUserAccountsParams) (api.GetV1DiningUserAccountsRes, error) {
	return dining_users.GetUserAccounts(shared.CbordBaseUrl, cbordInstitutionId, params)
}

func (h *Handler) GetV1DiningUserBarcode(ctx context.Context, params api.GetV1DiningUserBarcodeParams) (api.GetV1DiningUserBarcodeRes, error) {
	return dining_users.GetUserBarcode(shared.CbordBaseUrl, params)
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
