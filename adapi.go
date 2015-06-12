package adapi

import (
	"errors"
	"log"
	"sort"
	"time"

	"github.com/gernest/period"
	"github.com/jinzhu/now"
)

type Store interface {
	Set(key interface{}, value interface{}) error
	Get(key interface{}) (interface{}, error)
}

type MemoryStore struct {
	d map[string]interface{}
}
type ShowTime struct {
	Period period.Period
	Airs   AirTimes
}

type Air struct {
	Period period.Period
	Data   interface{}
}

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

type AirTimes []*Air

func (s AirTimes) Len() int { return len(s) }
func (s AirTimes) Less(i, j int) bool {
	b := &s[i].Period
	return b.IsBefore(s[j].Period)
}
func (s AirTimes) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func NewMemSTore() *MemoryStore {
	return &MemoryStore{make(map[string]interface{})}
}

func (m *MemoryStore) Set(key interface{}, value interface{}) error {
	var k string
	switch key.(type) {
	case string:
		k = key.(string)
	default:
		return errors.New("unsupported key type")
	}
	m.d[k] = value
	return nil
}

func (m *MemoryStore) Get(key interface{}) (interface{}, error) {
	var k string
	switch key.(type) {
	case string:
		k = key.(string)
	default:
		return nil, errors.New("unsupported key type")
	}
	if v, ok := m.d[k]; ok {
		return v, nil
	}
	return nil, errors.New("key not found")
}

func NewShowTime(start time.Time, duration time.Duration) *ShowTime {
	p, err := period.CreateFromDuration(start, duration)
	if err != nil {
		// TODO: Log this?
	}
	return &ShowTime{Period: p}
}
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

func NewAir(start time.Time, duration time.Duration, Data interface{}) *Air {
	p, err := period.CreateFromDuration(start, duration)
	if err != nil {
		//TODO: log this?
	}
	return &Air{p, Data}

}

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
