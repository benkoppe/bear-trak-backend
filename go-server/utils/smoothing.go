package utils

import (
	"math"
	"time"
)

// SmoothTime applies a rolling average over a time window.
// - points: input slice (must be sorted by timestamp ascending)
// - window: duration around each point (e.g. 10 * time.Minute)
// - get: extracts a float64 from T
// - set: rebuilds T with the smoothed value
// - ts: extracts the timestamp from T
func SmoothTime[T any](
	points []T,
	window time.Duration,
	get func(T) float64,
	set func(T, float64) T,
	ts func(T) time.Time,
) []T {
	if window <= 0 || len(points) == 0 {
		return points
	}
	smoothed := make([]T, len(points))

	for i, p := range points {
		center := ts(p)
		lower := center.Add(-window)
		upper := center.Add(window)

		var sum, weightSum float64
		for _, q := range points {
			t := ts(q)
			if t.Before(lower) || t.After(upper) {
				continue
			}
			// weight: closer in time gets more emphasis
			dt := math.Abs(center.Sub(t).Minutes())
			weight := 1.0 / (1.0 + dt) // exponential-ish decay
			sum += get(q) * weight
			weightSum += weight
		}

		avg := sum / weightSum
		smoothed[i] = set(p, avg)
	}

	return smoothed
}
