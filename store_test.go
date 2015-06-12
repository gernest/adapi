package adapi

import "testing"

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
