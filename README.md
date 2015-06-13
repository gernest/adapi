# adapi [![Build Status](https://drone.io/github.com/gernest/adapi/status.png)](https://drone.io/github.com/gernest/adapi/latest)
Personal advertisment management API service.

## Warning
This is experimental, for educational only. The problem at hand is more complex. But if you know what you are doing, its okay and you can help make this project better.

## Overview
Adapi is a library, providing handlers to construct a personalised advertisment service. It uses JSON as the main data exchange format.

An ad is persceived to be an `interface{}` that is inside a air time. Air time is a period in which the given ad should be broadcasted.

A channel is an object, which has a schedule(or a slice of showtimes). Each showtime is limited to a given period. showtimes contains the air times( which have ads). Now you can dig the source to see the definitions of these types, they are straight.

## Storage
Adapi can be used with whatever storage you prefer, it should only satisfy the following interface

```go
type Store interface {
	Set(key interface{}, value interface{}) error
	Get(key interface{}) (interface{}, error)
}
```

I have implemented a simple store `MemoryStore` which uses a map for key, value lookup. Read [here](store.go) to see how to implement a store backend.

## Confussion about Channels
I haven't yet found a good name to describe better  a sort of broadcasting entity, with a schedule and pumps different data depending on what time of the day other than a channel.

So, it might be confusing using this word channel, which might also be refering to golang channels. To make stuffs clear. The term "golang channel" will be used to refer to Go programming languaguage channels and "channel" or "AdAPI channel" will be used to refer to objects of type `Channel` throughout the project

## Does this use golang channels?
Nope, havent' found the need for them yet. But if you think they can be useful here, feel free to vent about it.'

## Usage
Initializing the handlers requires an `AdAPI` instance, which is created by the following method.

```go
func NewAdAPI(s Store) *AdAPI {
	return &AdAPI{s}
}
```

So, you only need to implement the `Store` interface to initialize the handler.


The following are examples of how the handlers can be used.

### Post
```go
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
```


### Get
`AdAPI.Get` uses url query paramenters to decide which part of the channel to render. For instance the url with  a query like this `?chn=adapi&dir=air"`, means get what is on air, for the channel with channel name of adapi.

Check this example
```go

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

```

Contributing
============

Please feel free to submit issues, fork the repository and send pull requests!

Contributions are welcome and will be fully credited. Please see [CONTRIBUTING](CONTRIBUTING.md) for details.

## Author
Geofrey Ernest

## License

This project is under the MIT License. See the [LICENSE](LICENCE) file for the full license text.
