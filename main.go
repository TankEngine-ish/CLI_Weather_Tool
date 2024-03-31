package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the API key from the environment variable
	apiKey := os.Getenv("WEATHER_API_KEY")

	res, err := http.Get("http://api.weatherapi.com/v1/current.json?key=" + apiKey + "&q=Varna&aqi=no")
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

	fmt.Println(string(body))
}
