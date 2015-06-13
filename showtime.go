package adapi

import (
	"errors"
	"log"
	"sort"
	"time"

	"github.com/gernest/period"
)

// ShowTime holds air times in a given period
type ShowTime struct {
	Period period.Period `json:"period"`
	Airs   AirTimes      `json:"airs"`
}

// AirTimes a sortable list of air times
type AirTimes []*Air

func (s AirTimes) Len() int { return len(s) }
func (s AirTimes) Less(i, j int) bool {
	b := &s[i].Period
	return b.IsBefore(s[j].Period)
}
func (s AirTimes) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// NewShowTime initializes a new showtime
func NewShowTime(start time.Time, duration time.Duration) *ShowTime {
	p, err := period.CreateFromDuration(start, duration)
	if err != nil {
		// TODO: Log this?
	}
	return &ShowTime{Period: p}
}

// Add ads a air time to showtime
func (s *ShowTime) Add(a *Air) error {
	sort.Sort(s.Airs)
	if s.Period.Contains(a.Period.Start) {
		var isOk = true
		for _, v := range s.Airs {
			vp := &v.Period
			if vp.Contains(a.Period.Start) || vp.Start.Equal(a.Period.Start) {
				log.Println(vp.Start)
				isOk = false
				break
			}

		}
		if isOk {
			s.Airs = append(s.Airs, a)
			return nil
		}
		return errors.New("airtime already taken")
	}
	return errors.New("airtime out of range")
}

// Show returns what is in the air right now.
func (s *ShowTime) Show() *Air {
	return s.showAt(time.Now())
}
func (s *ShowTime) showAt(t time.Time) *Air {
	for _, v := range s.Airs {
		if v.Period.Contains(t) {
			return v
		}
	}
	return nil
}
