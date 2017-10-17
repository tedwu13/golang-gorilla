package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", hello) // handler function at root path of web server also known as the Serve Mux

	http.HandleFunc("/weather/", weatherHandler)
	http.ListenAndServe(":8080", nil) //port, handler
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!")) //Response Writers write responses to the client
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	mw := multiWeatherProvider{
		openWeatherMap{},
		weatherUnderground{apiKey: "0123456789abcdef"},
	}
	begin := time.Now()
	city := strings.SplitN(r.URL.Path, "/", 3)[2]

	// data, err := query(city)
	temp, err := mw.temperature(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8") //set Headers
	// json.NewEncoder(w).Encode(data)
	// d := map[string]int{"apple": 5, "lettuce": 7}
	// d := map[string]float {"kelvin": 21.11}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"city": city,
		"temp": temp,
		"took": time.Since(begin).String(),
	})
}

//To do that, we leverage Go’s concurrency primitives: goroutines and channels.
//We’ll spawn each API query in its own goroutine, which will run concurrently.
//We’ll collect the responses in a single channel, and perform the average calculation when everything is finished.

//{
//     "name": "Tokyo",
//     "coord": {
//         "lon": 139.69,
//         "lat": 35.69
//     },
//     "weather": [
//         {
//             "id": 803,
//             "main": "Clouds",
//             "description": "broken clouds",
//             "icon": "04n"
//         }
//     ],
//     "main": {
//         "temp": 296.69,
//         "pressure": 1014,
//         "humidity": 83,
//         "temp_min": 295.37,
//         "temp_max": 298.15
//     }
// }
type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

type weatherProvider interface {
	temperature(city string) (float64, error) // in Kelvin, naturally
}

// the difference between an interface and a struct cintauns naned fields
// Interface has a method set/list of methods that a type must have in order to implement the interface

// Use the encoding/json packageto directly unmarshal API's response to our struct
// UnMarshal parses JSON encoded data and stores the result into an interface
// Marshal convers json encoding of a struct/interface
// Good for type safety vs dynamic languages like Python and Ruby
//

//If the http.Get succeeds, we defer a call to close the response body, which will execute when we leave the function scope (when we return from the query function) and is an elegant form of resource management.
//Meanwhile, we allocate a weatherData struct, and use a json.Decoder to unmarshal from the response body directly into our struct.
func query(city string) (weatherData, error) {
	url := "http://api.openweathermap.org/data/2.5/weather?APPID=5eb262c4c99f01b45b71029154344115&q=" + city
	fmt.Printf("URL: %s\n", url)
	resp, err := http.Get(url)

	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var d weatherData

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}

	return d, nil

}

type openWeatherMap struct{}
type weatherUnderground struct {
	apiKey string
}

// func (w weatherUnderground) temperature(city string) (float64, error) {
// 	resp, err := http.Get("http://api.wunderground.com/api/" + w.apiKey + "/conditions/q/" + city + ".json")
// 	if err != nil {
// 		return 0, err
// 	}

// 	defer resp.Body.Close()

// 	var d struct {
// 		Observation struct {
// 			Celsius float64 `json:"temp_c"`
// 		} `json:"current_observation"`
// 	}

// 	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
// 		return 0, err
// 	}

// 	kelvin := d.Observation.Celsius + 273.15
// 	log.Printf("weatherUnderground: %s: %.2f", city, kelvin)
// 	return kelvin, nil
// }

// func (w openWeatherMap) temperature(city string) (float64, error) {
// 	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=5eb262c4c99f01b45b71029154344115&q=" + city)
// 	if err != nil {
// 		return 0, err
// 	}

// 	defer resp.Body.Close()

// 	var d struct {
// 		Main struct {
// 			Kelvin float64 `json:"temp"`
// 		} `json:"main"`
// 	}

// 	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
// 		return 0, err
// 	}

// 	log.Printf("openWeatherMap: %s: %.2f", city, d.Main.Kelvin)
// 	return d.Main.Kelvin, nil
// }

type multiWeatherProvider []weatherProvider

func (w multiWeatherProvider) temperature(city string) (float64, error) {
	sum := 0.0

	for _, provider := range w {
		k, err := provider.temperature(city)
		if err != nil {
			return 0, err
		}

		sum += k
	}

	return sum / float64(len(w)), nil
}

func (w multiWeatherProvider) temperature(city string) (float64, error) {
	// Make a channel for temperatures, and a channel for errors.
	// Each provider will push a value into only one.
	temps := make(chan float64, len(w))
	errs := make(chan error, len(w))

	//Channels provide a way for two goroutines to communicate with one another and synchronize their execution. Here is an example program using channels

	// For each provider, spawn a goroutine with an anonymous function.
	// That function will invoke the temperature method, and forward the response.
	for _, provider := range w {
		go func(p weatherProvider) {
			k, err := p.temperature(city)
			if err != nil {
				errs <- err
				return
			}
			temps <- k
		}(provider)
	}

	sum := 0.0

	// Collect a temperature or an error from each provider.
	for i := 0; i < len(w); i++ {
		select {
		case temp := <-temps:
			sum += temp
		case err := <-errs:
			return 0, err
		}
	}

	// Return the average, same as before.
	return sum / float64(len(w)), nil
}
