package adapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
}
