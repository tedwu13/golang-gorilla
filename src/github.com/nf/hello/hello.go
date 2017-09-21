package main

import "net/http"

func main() {
    http.HandleFunc("/", hello)
    // http.HandleFunc('router', handler)
    http.ListenAndServe(":8080", nil)
}
//ListenAndServe starts an HTTP server with a given address and handler. The handler is usually nil, which means to use DefaultServeMux.

func hello(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("hello!"))
}

// Using JSON and go using a Marshal function
// JSON is used to communciate with web backends with Javascript programs running in the browse. Marshal makes it easier to read and write JSOn data from Go Programs
// To encode JSON data, we use Marshal function


type weatherData struct {
    Name string `json:"name"`
    Coordinate struct {
        Longitude float64 `json:"lon"`
        Latitude float64 `json:"lat"`
    }
    Main struct {
        Kelvin float64 `json:"temp"`
    } `json:"main"`
}

// Example json response looks like this
// 
// {
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
// 
func query(city string) (weatherData, error) {
    resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=YOUR_API_KEY&q=" + city)
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

