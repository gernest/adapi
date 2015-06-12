package adapi

import (
	"errors"
	"sort"
	"time"

	"github.com/jinzhu/now"
)

type Channel struct {
	Name  string
	Shows Shows
}

type Shows []*ShowTime

func (s Shows) Len() int { return len(s) }
func (s Shows) Less(i, j int) bool {
	b := &s[i].Period
	return b.IsBefore(s[j].Period)
}
func (s Shows) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (c *Channel) AddShow(s *ShowTime) {
	c.Shows = append(c.Shows, s)
}

func CreateDaySchedule(c *Channel) *Channel {
	begin := now.BeginningOfDay()
	start := begin
	interval := time.Hour
	for _ = range make([]struct{}, 24) {
		s := NewShowTime(start, interval)
		start = s.Period.End
		c.AddShow(s)
	}
	sort.Sort(c.Shows)
	return c
}

func AddAirTime(c *Channel, air *Air) error {
	for _, v := range c.Shows {
		if v.Period.Contains(air.Period.Start) {
			return v.Add(air)
		}
	}
	// TODO: meaningful message?
	return errors.New("failed to add airtime")
}
