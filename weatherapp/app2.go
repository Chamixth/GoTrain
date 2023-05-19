package main

import (
	"encoding/json"               // Package for JSON encoding and decoding
	"fmt"                         //package for formatted IO
	"github.com/labstack/echo/v4" //Package for building web application
	"io/ioutil"                   //Package for IO functions
	"net/http"                    //Package for HTTP client functionalities
)

func welcomePage(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the weather app") //return a string
}

func weatherDetails(c echo.Context) error {
	city := c.QueryParam("city")                 // Retrieve the value of "city" query parameter
	apiKey := "f49917d741bcdea1373729fe389afa13" //Api key
	//Then created the API URL
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)

	response, err := http.Get(url) //sent the http get request to the api url
	if err != nil {
		return err //return an error if there is any problem
	}
	defer response.Body.Close() //ensure the response body is closed

	body, err := ioutil.ReadAll(response.Body) //read the response
	if err != nil {
		return err //return an error if there are any
	}

	type WeatherData struct {
		Main struct {
			Temperature float64 `json:"temp"`     //struct field for the temp
			Humidity    int     `json:"humidity"` //struct field for the humidity
			Pressure    int     `json:"pressure"` //struct field for the pressure
		} `json:"main"`           //struct the field for the main weatherData with json tag
		Name string `json:"name"` //field for the city name
	}

	var weatherData WeatherData                                //created the variable to get the unmarshaled weather data
	if err := json.Unmarshal(body, &weatherData); err != nil { //unmarshal the json into the weatherdata variable
		return err
	}

	temperatureCelcius := weatherData.Main.Temperature - 273.15 // convert the kelvin to celcius

	// Build a response map with the required fields
	responseData := map[string]interface{}{
		"City":        weatherData.Name,
		"Humidity":    weatherData.Main.Humidity,
		"Pressure":    weatherData.Main.Pressure,
		"Temperature": temperatureCelcius,
	}

	return c.JSON(http.StatusOK, responseData) //return json response with weather details
}

func main() {
	e := echo.New()
	e.GET("/welcome", welcomePage)    //route for welcome page
	e.GET("/weather", weatherDetails) //route for weather page
	e.Logger.Fatal(e.Start(":8080"))  //start the server on 8080 port
}
