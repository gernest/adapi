package adapi

import (
	"time"

	"github.com/gernest/period"
)

// Air is the airtime, whose property Data holds the advertisement
type Air struct {
	Period period.Period `json:"period"`

	// Data is the actual advertisement.
	Data interface{} `json:"data"`
}

// NewAir initializes a new air
func NewAir(start time.Time, duration time.Duration, Data interface{}) *Air {
	p, err := period.CreateFromDuration(start, duration)
	if err != nil {
		//TODO: log this?
	}
	return &Air{p, Data}
}
