package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	city := strings.SplitN(r.URL.Path, "/", 3)[2]

	data, err := query(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8") //set Headers
	json.NewEncoder(w).Encode(data)
}

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
