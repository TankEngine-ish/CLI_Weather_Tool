package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

// this struct maps to the json response of the API call
type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {

	q := "Varna"
	if len(os.Args) >= 2 {
		q = os.Args[1]

	}

	// Loads the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Gets the API key from the environment variable
	apiKey := os.Getenv("WEATHER_API_KEY")

	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=" + apiKey + "&q=" + q + "&days=1&aqi=no&alerts=no")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Failed to fetch data")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour
	fmt.Printf("Location: %s, %s: %.0fC, %s\n",
		location.Name,
		location.Country,
		current.TempC,
		current.Condition.Text,
	)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf("%s - %.0fC, %.0f%%, %s\n",
			date.Format("02-01-2006 15:04:05"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)
		// prints the message in red if the chance of rain is above 50%
		if hour.ChanceOfRain < 50 {
			fmt.Print(message)
		} else {
			color.Red(message)
		}
	}
}
