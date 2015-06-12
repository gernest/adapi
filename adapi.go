package adapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
)

var (
	ErrChanOffline = errors.New("sorry the channel is offline")
	contentType    = "Content-Type"
	contentJson    = "application/json"
	charset        = "charset=utf-8"
)

type AdAPI struct {
	s Store
}
type AdAPIRequest struct {
	ChannelName string
	Air         *Air
}

func NewAdAPI(s Store) *AdAPI {
	return &AdAPI{s}
}

type jsonError struct {
	Error string `json:"error"`
}

func (api *AdAPI) Get(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	w.Header().Set(contentType, fmt.Sprintf("%s;%s", contentJson, charset))
	if r.Method == "GET" {
		channelName := r.URL.Query().Get("chn")
		directive := r.URL.Query().Get("dir")
		c, err := api.s.Get(channelName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			encoder.Encode(&jsonError{err.Error()})
			return
		}
		ch := c.(*Channel)
		w.WriteHeader(http.StatusOK)
		switch directive {
		case "schedule":
			sort.Sort(ch.Shows)
			encoder.Encode(ch.Shows)
		case "show":
			show := ch.Show()
			if show != nil {
				encoder.Encode(show)
				return
			}
			encoder.Encode("nothing to show")
			return
		case "air":
			show := ch.Show()
			if show != nil {
				onAir := show.Show()
				if onAir != nil {
					encoder.Encode(onAir)
					return
				}
			}
			encoder.Encode("nothing to show")
			return
		}

	}
	w.WriteHeader(http.StatusNotFound)
	encoder.Encode(&jsonError{ErrChanOffline.Error()})
	return
}

func (api *AdAPI) Post(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	w.Header().Set(contentType, fmt.Sprintf("%s;%s", contentJson, charset))
	if r.Method == "POST" {
		req := &AdAPIRequest{}
		buf := &bytes.Buffer{}
		io.Copy(buf, r.Body)
		err := json.Unmarshal(buf.Bytes(), req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(&jsonError{err.Error()})
			return
		}
		ch, err := api.postToChannel(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			encoder.Encode(&jsonError{err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		encoder.Encode(ch)
		return
	}
}

func (api *AdAPI) postToChannel(req *AdAPIRequest) (*Channel, error) {
	return postAChannel(api.s, req)
}

func postAChannel(s Store, req *AdAPIRequest) (*Channel, error) {
	c, _ := s.Get(req.ChannelName)
	if c != nil {
		ch := c.(*Channel)
		err := AddAirTime(ch, req.Air)
		if err != nil {
			return nil, err
		}
		s.Set(ch.Name, ch)
		return ch, nil
	}
	ch := CreateDaySchedule(&Channel{Name: req.ChannelName})
	err := AddAirTime(ch, req.Air)
	if err != nil {
		return nil, err
	}
	s.Set(ch.Name, ch)
	return ch, nil
}
