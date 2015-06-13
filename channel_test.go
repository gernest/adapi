package adapi

import (
	"testing"
	"time"
)

func TestCreateDaySchedule(t *testing.T) {
	c := &Channel{}
	CreateDaySchedule(c)
	if len(c.Shows) != 24 {
		t.Errorf("expected %d got %d", 24, len(c.Shows))
	}
}

func TestAddAirTime(t *testing.T) {
	c := &Channel{}
	CreateDaySchedule(c)
	a := NewAir(time.Now(), time.Minute, "bogus")
	err := AddAirTime(c, a)
	if err != nil {
		t.Errorf("adding airtime %v", err)
	}
	a1 := NewAir(time.Now().AddDate(1, 0, 0), time.Hour, "out of range")
	err = AddAirTime(c, a1)
	if err == nil {
		t.Error("expected an error")
	}
}
