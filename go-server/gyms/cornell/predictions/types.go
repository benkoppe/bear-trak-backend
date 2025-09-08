package predictions

import "time"

type Prediction struct {
	GymName        string
	Timestamp      time.Time
	Predicted      int
	PredictionMade time.Time
}
