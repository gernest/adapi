package adapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func requestPost(urlPath, channelName string, data interface{}, start time.Time, duration time.Duration) (*http.Request, error) {
	req := &AdAPIRequest{
		ChannelName: channelName,
		Air:         NewAir(start, duration, data),
	}
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest("POST", urlPath, strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}
	r.Header.Set(contentType, contentJson)
	r.Header.Set("X-Requested-With", "XMLHttpRequest")
	return r, nil
}

func requestGet(channel, directive string) (*http.Request, error) {
	vars := url.Values{
		"chn": {channel},
		"dir": {directive},
	}
	urlPath := fmt.Sprintf("http://adapi.io?%s", vars.Encode())
	return http.NewRequest("GET", urlPath, nil)
}

func TestAdAPI_Post(t *testing.T) {
	urlPath := "http://adapi.io"
	data := "hello"
	now := time.Now()
	req, err := requestPost(urlPath, "charm", data, now, time.Hour)
	if err != nil {
		t.Errorf("creating request %v", err)
	}
	w := httptest.NewRecorder()
	api := NewAdAPI(NewMemSTore())
	api.Post(w, req)
	if !strings.Contains(w.Body.String(), data) {
		t.Errorf("expected %s to contain %s", w.Body.String(), data)
	}

	req1, err := requestPost(urlPath, "charm", data, now, time.Hour)
	if err != nil {
		t.Errorf("creating request %v", err)
	}
	w1 := httptest.NewRecorder()
	api.Post(w1, req1)
	if !strings.Contains(w1.Body.String(), "error") {
		t.Errorf("expected %s to contain %s", w1.Body.String(), "error")
	}

	req2, err := requestPost(urlPath, "charm", data, now.Add(time.Hour), time.Hour)
	if err != nil {
		t.Errorf("creating request %v", err)
	}
	w2 := httptest.NewRecorder()
	api.Post(w2, req2)
	if !strings.Contains(w2.Body.String(), data) {
		t.Errorf("expected %s to contain %s", w2.Body.String(), data)
	}
}

func TestAdAPI_Get(t *testing.T) {
	channelName := "adapi"
	data := "my time"

	// no channel
	req, err := requestGet(channelName, "")
	if err != nil {
		t.Errorf("creating request %v", err)
	}
	store := NewMemSTore()
	w := httptest.NewRecorder()
	api := NewAdAPI(store)
	api.Get(w, req)

	// create a channel
	c := &Channel{Name: channelName}
	CreateDaySchedule(c)
	a := NewAir(time.Now(), time.Minute, data)
	err = AddAirTime(c, a)
	if err != nil {
		t.Errorf("adding airtime %v", err)
	}

	// defaults
	store.Set(channelName, c)
	w1 := httptest.NewRecorder()
	api.Get(w1, req)

	// get  whats on air
	req2, err := requestGet(channelName, "air")
	if err != nil {
		t.Errorf("creating request %v", err)
	}
	w2 := httptest.NewRecorder()
	api.Get(w2, req2)

	// get the showtime
	req3, err := requestGet(channelName, "show")
	if err != nil {
		t.Errorf("creating request %v", err)
	}
	w3 := httptest.NewRecorder()
	api.Get(w3, req3)

	req4, err := requestGet(channelName, "schedule")
	if err != nil {
		t.Errorf("creating request %v", err)
	}
	w4 := httptest.NewRecorder()
	api.Get(w4, req4)
}
