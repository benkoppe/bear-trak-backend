package schools

import (
	"fmt"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/db"
	"github.com/benkoppe/bear-trak-backend/go-server/schools/cornell"
)

type SchoolCode string

const (
	Cornell SchoolCode = "cornell"
	UMich   SchoolCode = "umich"
)

func NewHandler(code SchoolCode, db *db.Queries) (api.Handler, error) {
	switch code {
	case Cornell:
		return cornell.NewHandler(db), nil
	default:
		return nil, fmt.Errorf("unsupported school: %s", code)
	}
}
