package adapi

import (
	"time"

	"github.com/gernest/period"
)

type Air struct {
	Period period.Period
	Data   interface{}
}

func NewAir(start time.Time, duration time.Duration, Data interface{}) *Air {
	p, err := period.CreateFromDuration(start, duration)
	if err != nil {
		//TODO: log this?
	}
	return &Air{p, Data}
}
