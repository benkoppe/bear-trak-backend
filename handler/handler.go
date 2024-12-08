package handler

import (
	"context"

	backend "github.com/benkoppe/bear-trak-backend/backend"
	"github.com/benkoppe/bear-trak-backend/dining"
	"github.com/benkoppe/bear-trak-backend/gyms"
)

type BackendService struct{}

func (bs *BackendService) GetV1Dining(ctx context.Context) ([]backend.Eatery, error) {
	return dining.Get("https://now.dining.cornell.edu/api/1.0/dining/eateries.json")
}

func (bs *BackendService) GetV1Gyms(ctx context.Context) ([]backend.Gym, error) {
	return gyms.Get("https://connect2concepts.com/connect2/?type=bar&key=355de24d-d0e4-4262-ae97-bc0c78b92839&loc_status=false")
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
