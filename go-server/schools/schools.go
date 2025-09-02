// Package schools contains all general cross-school methods.
package schools

import (
	"fmt"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/db"
	dining_email "github.com/benkoppe/bear-trak-backend/go-server/dining/cornell/email"
	"github.com/benkoppe/bear-trak-backend/go-server/schools/cornell"
	"github.com/benkoppe/bear-trak-backend/go-server/schools/harvard"
	_ "github.com/benkoppe/bear-trak-backend/go-server/schools/shared"
	"github.com/benkoppe/bear-trak-backend/go-server/schools/umich"
)

type (
	SchoolCode string
)

const (
	Cornell SchoolCode = "cornell"
	UMich   SchoolCode = "umich"
	Harvard SchoolCode = "harvard"
)

func NewHandler(code SchoolCode, db *db.Queries, config *Config) (api.Handler, error) {
	switch code {
	case Cornell:
		return cornell.NewHandler(db, config.HouseDinnerCache), nil
	case UMich:
		return umich.NewHandler(db), nil
	case Harvard:
		return harvard.NewHandler(db), nil
	default:
		return nil, fmt.Errorf("unsupported school: %s", code)
	}
}

type Config struct {
	EnabledGymCapacities bool
	HouseDinnerCache     dining_email.Cache
}

func GetConfig(code SchoolCode) (*Config, error) {
	switch code {
	case Cornell:
		houseDinnerCache := cornell.InitHouseDinnerCache()
		return &Config{
			EnabledGymCapacities: true,
			HouseDinnerCache:     houseDinnerCache,
		}, nil
	case UMich:
		return &Config{
			EnabledGymCapacities: false,
		}, nil
	case Harvard:
		return &Config{
			EnabledGymCapacities: false,
		}, nil
	}

	return nil, fmt.Errorf("unsupported school: %s", code)
}
