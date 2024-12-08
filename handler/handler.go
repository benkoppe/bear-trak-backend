package handler

import (
	"context"
	"fmt"

	backend "github.com/benkoppe/bear-trak-backend/backend"
	"github.com/benkoppe/bear-trak-backend/dining"
)

type BackendService struct{}

func (bs *BackendService) GetV1Dining(ctx context.Context) ([]backend.Eatery, error) {
	return dining.Get("https://now.dining.cornell.edu/api/1.0/dining/eateries.json")
}

func (bs *BackendService) GetV1Gyms(ctx context.Context) ([]backend.Gym, error) {
	return nil, fmt.Errorf("Not implemented.")
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
