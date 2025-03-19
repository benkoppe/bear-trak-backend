package time_utils

import (
	"fmt"
	"time"
)

// resillient datetime parser (good for both dates and times)
// will try multiple layouts in order before failing
// requires a list of templates
func ParseDateTime(s string, layouts []string) (time.Time, error) {
	var err error
	for _, layout := range layouts {
		if t, err2 := time.Parse(layout, s); err2 == nil {
			return t, nil
		} else {
			err = err2
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse datetime %q: %v", s, err)
}
