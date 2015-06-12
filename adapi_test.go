package adapi

import (
	"testing"
	"time"
)

func TestMemoryStore_Set(t *testing.T) {
	result := []struct {
		key, value string
	}{
		{"adapi", "a personal ad server"},
		{"gernest", "Learning the hard way"},
	}
	m := NewMemSTore()
	for _, v := range result {
		err := m.Set(v.key, v.value)
		if err != nil {
			t.Errorf("setting key %v", err)
		}
	}
	err := m.Set(1, "one")
	if err == nil {
		t.Error("expected an error")
	}

}
func TestMemoryStore_Get(t *testing.T) {
	result := []struct {
		key, value string
	}{
		{"adapi", "a personal ad server"},
		{"gernest", "Learning the hard way"},
	}
	m := NewMemSTore()
	for _, v := range result {
		err := m.Set(v.key, v.value)
		if err != nil {
			t.Errorf("setting key %v", err)
		}
	}
	for _, v := range result {
		d, err := m.Get(v.key)
		if err != nil {
			t.Errorf("setting key %v", err)
		}
		if d.(string) != v.value {
			t.Errorf("expected %s got %s", v.value, d)
		}
	}
	d, err := m.Get(123)
	if err == nil {
		t.Error("expected an error")
	}
	if d != nil {
		t.Errorf("expected nil got %v", d)
	}

	d, err = m.Get("nothing")
	if err == nil {
		t.Error("expected an error")
	}
	if d != nil {
		t.Errorf("expected nil got %v", d)
	}
}

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

func TestShowTime_Show(t *testing.T) {
	s := NewShowTime(time.Now().Add(-time.Hour), 2*time.Hour)
	n := time.Now()
	d := time.Minute
	airs := []*Air{
		NewAir(n.Add(-d), d, nil),
		NewAir(n, d, nil),
		NewAir(n.Add(d), d, nil),
	}
	for k, v := range airs {
		v.Data = k
		err := s.Add(v)
		if err != nil {
			t.Errorf("adding air %v", err)
		}
	}
	a := s.Show()
	if a == nil {
		t.Errorf("expected air time")
	}
	if a.Data.(int) != 1 {
		t.Errorf("expected 1 got %v", a.Data)
	}
}
