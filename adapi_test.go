package adapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func requestPost(urlPath, channelName string, data interface{}, start time.Time, duration time.Duration) (*http.Request, error) {
	req := &Request{
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
	r.Header.Set(contentType, contentJSON)
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
	api := NewAdAPI(NewMemStore())
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
	store := NewMemStore()
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

func ExampleAdAPI_Post() {
	var (
		postPath    = "/adapi/post"
		channelName = "adapi"
		start       = time.Now()
		duration    = time.Hour
		data        = "hello world"
	)

	store := NewMemStore()
	api := NewAdAPI(store)

	// create a handle
	mux := http.NewServeMux()
	mux.HandleFunc(postPath, api.Post)

	// we are using a test server for this example
	server := httptest.NewServer(mux)
	defer server.Close()

	// The test server runs on random port, the only way to hit the correct socket is by
	// reconstructing the url
	currentPostPath := fmt.Sprintf("%s%s", server.URL, postPath)

	// create  data to be sent with the request
	req := &Request{
		ChannelName: channelName,
		Air:         NewAir(start, duration, data),
	}
	reqData, err := json.Marshal(req)
	if err != nil {
		// do something
	}

	// create a client.
	client := &http.Client{}
	response, err := client.Post(currentPostPath, "application/json", strings.NewReader(string(reqData)))
	if err != nil {
		// do something
	}
	buf := &bytes.Buffer{}
	io.Copy(buf, response.Body)
	defer response.Body.Close()

	// The response should be a Channel, this Channel will contain the airtime we have posted
	ch := &Channel{}
	err = json.Unmarshal(buf.Bytes(), ch)
	if err != nil {
		// do something
	}
	currentShowTime := ch.Show()
	currentAiring := currentShowTime.Show()

	// Within an hour range, we will should get the same air data.
	fmt.Println(currentAiring.Data)

	//Output:
	//hello world
}

func ExampleAdAPI_Get() {
	var (
		getPath     = "/adapi/get"
		channelName = "adapi"
		data        = "hello world"
	)
	store := NewMemStore()

	// create a channel and an air time
	c := &Channel{Name: channelName}
	CreateDaySchedule(c)
	a := NewAir(time.Now(), time.Minute, data)
	err := AddAirTime(c, a)
	if err != nil {
		// do something
	}

	// store the channel
	store.Set(c.Name, c)

	api := NewAdAPI(store)

	mux := http.NewServeMux()
	mux.HandleFunc(getPath, api.Get)

	server := httptest.NewServer(mux)
	defer server.Close()

	client := &http.Client{}
	vars := url.Values{
		"chn": {channelName},
		"dir": {"air"},
	}
	currentGetPath := fmt.Sprintf("%s%s?%s", server.URL, getPath, vars.Encode())
	response, err := client.Get(currentGetPath)
	if err != nil {
		// do something
	}

	buf := &bytes.Buffer{}
	io.Copy(buf, response.Body)
	defer response.Body.Close()

	air := &Air{}
	err = json.Unmarshal(buf.Bytes(), air)
	if err != nil {
		// do something
	}
	fmt.Println(air.Data)

	// Output:
	// hello world

}
