package cornell

import (
	_ "embed"

	"github.com/benkoppe/bear-trak-backend/go-server/alerts"
	"github.com/benkoppe/bear-trak-backend/go-server/api"
)

//go:embed alerts.json
var alertBytes []byte

func Get() ([]api.Alert, error) {
	return alerts.Get(alertBytes)
}
