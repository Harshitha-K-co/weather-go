package main

import (
	"encoding/json" //In this code, for the function "query", we are using the "encoding/json" package to parse the JSON response from the OpenWeatherMap API.
	"log"
	"net/http" //In this code, for the function "query", we are using the "net/http" package to make HTTP requests to the OpenWeatherMap API.
	"os"       //In this code, for the function "loadApiConfig", we are using the "os" package to read the API configuration file.
	"strings"  //In this code, for the function "query", we are using the "strings" package to split the URL path to extract the city name.
)

// apiConfig holds the API key for OpenWeatherMap. This struct is used to unmarshal the JSON configuration file.
// The backtick syntax is used to define struct tags in Go, which provide metadata about the struct fields.
// The `json:"OpenWeatherMapApiKey"` tag indicates that when this struct is marshaled to or unmarshaled from JSON, the field should be represented with the key "OpenWeatherMapApiKey".
// This is useful for ensuring that the JSON keys match the expected field names in the struct.
type apiConfig struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"` // means this field is expected to be in the JSON file with the key "OpenWeatherMapApiKey".
}

type weatherData struct { // weatherData struct represents the structure of the weather data returned by the OpenWeatherMap API.
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (apiConfig, error) { //
	bytes, err := os.ReadFile(filename) // This line reads the contents of the file specified by `filename` into a byte slice. If the file does not exist or cannot be read, it returns an error.
	// conerted to bytes because the json.Unmarshal function expects a byte slice as input. json.unmarshal takes a byte slice containing JSON data and unmarshals it into the provided struct type.
	if err != nil {
		return apiConfig{}, err
	}

	var c apiConfig
	err = json.Unmarshal(bytes, &c) //whats hapepening here is that the JSON data read from the file is being unmarshaled into the `apiConfig` struct. The `json.Unmarshal` function takes a byte slice (the contents of the file) and a pointer to a struct (in this case, `&c`) and populates the struct with the data from the JSON.
	// The `&c` is a pointer to the `apiConfig` struct, which allows `json.Unmarshal` to modify the struct directly with the data it reads from the JSON.If the JSON data does not match the struct fields, `json.Unmarshal` will return an error.
	if err != nil {
		return apiConfig{}, err
	}
	return c, nil
}

func hello(w http.ResponseWriter, r *http.Request) { // hello is a simple HTTP handler function that responds with "Hello, World!" when accessed.
	// http.ResponseWriter is used to send the response back to the client, and *http.Request contains information about the HTTP request.
	w.Write([]byte("Hello, World!\n")) // w.Write writes the byte slice containing "Hello, World!\n" to the response writer, which sends it back to the client.
}

func query(city string) (weatherData, error) { // query is a function takes a city name as input and queries the OpenWeatherMap API for the current weather data of that city.
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err // If there is an error loading the API configuration, it returns an empty weatherData struct and the error.
	}
	resp, err := http.Get("http://api.Openweathermap.org/data/2.5/weather?q=" + city + "&appid=" + apiConfig.OpenWeatherMapApiKey) // This line constructs the URL for the API request
	// using the city name and the API key loaded from the configuration file. It uses the http.Get function to send a GET request to the OpenWeatherMap API.
	// The URL is constructed by concatenating the base URL of the OpenWeatherMap API with the city name and the API key.
	// If there is an error making the request, it returns an empty weatherData struct and the error.
	if err != nil {
		return weatherData{}, err
	}
	defer resp.Body.Close() // This line ensures that the response body is closed after the function completes, preventing resource leaks.
	// It defers the closing of the response body until the surrounding function returns.
	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil { // This line uses the json.NewDecoder function to create a new JSON decoder that reads from the response body.
		// It then calls the Decode method on the decoder to unmarshal the JSON data into the `weatherData` struct `d`.
		// If there is an error during decoding (for example, if the JSON response does not match the expected structure), it returns an empty weatherData struct and the error.
		// The `Decode` method reads the JSON data from the response body and populates the fields of the `weatherData` struct with the corresponding values.
		// If the decoding is successful, the `weatherData` struct `d` will contain the weather information for the specified city.
		return weatherData{}, err
	}
	return d, nil // If the decoding is successful, it returns the populated `weatherData` struct and a nil error, indicating that the query was successful.
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.SplitN(r.URL.Path, "/", 3)
		if len(parts) < 3 || parts[2] == "" {
			http.Error(w, "City name not provided", http.StatusBadRequest)
			return
		}
		data, err := query(parts[2])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})

	log.Println("Server listening on http://localhost:8081") // Listening from 8081 host
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}
