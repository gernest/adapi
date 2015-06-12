/*
Package adapi is a personal advertisment management API service.
*/
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
	errChanOffline = errors.New("sorry the channel is offline")
	contentType    = "Content-Type"
	contentJSON    = "application/json"
	charset        = "charset=utf-8"
)

// AdAPI is the main type with the core handlers
type AdAPI struct {
	s Store
}

// Request is the request object, used in posting new airtime data.
type Request struct {
	ChannelName string `json:"name"`
	Air         *Air   `json:"air"`
}

// NewAdAPI creates a new AdAPI object
func NewAdAPI(s Store) *AdAPI {
	return &AdAPI{s}
}

type jsonError struct {
	Error string `json:"error"`
}

// Get retrieves data from the AdAPI server. This method relies on url queries to understand
// the request. To avoid long names, "chn" is used to represent the channel name, and "dir" is used
// to represent directive.
//
// Directives are strings, which have different meaning in the method. The following are the base directives.
//    schedule : represent a 24 hour schedule with 24 showtimes. Each showtime takes a duration of one hour.
//               this  is the value contained in the Shows property of a channel.
//    show     : represent the current show, it is a showtime in the channel whose period contains time.Now()
//    air      : whats on air right now, it is the airtime whose period contains time.Now()
//
// So, if you want to get what is currently on air in adapi channel, you can write a query like
//   ?chn=adapi&dir=air
func (api *AdAPI) Get(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	w.Header().Set(contentType, fmt.Sprintf("%s;%s", contentJSON, charset))
	if r.Method == "GET" {
		channelName := r.URL.Query().Get("chn")
		directive := r.URL.Query().Get("dir")
		c, err := api.s.Get(channelName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
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
		default:
			encoder.Encode(ch)
			return
		}

	}
	w.WriteHeader(http.StatusNotFound)
	encoder.Encode(&jsonError{errChanOffline.Error()})
	return
}

// Post adds an Air to the channel, the request body should be of type Request
// if the channel is not found, a new channel is created.
func (api *AdAPI) Post(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	w.Header().Set(contentType, fmt.Sprintf("%s;%s", contentJSON, charset))
	if r.Method == "POST" {
		req := &Request{}
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

func (api *AdAPI) postToChannel(req *Request) (*Channel, error) {
	return postAChannel(api.s, req)
}

func postAChannel(s Store, req *Request) (*Channel, error) {
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
