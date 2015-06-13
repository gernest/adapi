package adapi

import (
	"testing"
	"time"
)

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
