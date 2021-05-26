package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Weather struct {
	Current struct {
		TempC     float64 `json:"temp_c"`
		IsDay     int     `json:"is_day"`
		Condition struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`

		WindKph    float64 `json:"wind_kph"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PrecipMm   float64 `json:"precip_mm"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		VisKm      float64 `json:"vis_km"`
	} `json:"current"`
}

func main() {
	var wh Weather
	starttime := time.Now()
	recived := true

	// Web server
	e := echo.New()
	e.Use(middleware.CORS())

	e.GET("/apiv1/:key/", func(c echo.Context) error {
		parameter := strings.ToLower(c.Param("key"))

		// Check if time is greater than a certain value
		nowtime := time.Now()
		elapsed := nowtime.Sub(starttime)
		elapsed2 := math.Round(elapsed.Seconds())
		fmt.Print("\n", elapsed2)

		if elapsed2 >= 1800 {
			starttime = time.Now()
			recived = true
		}

		switch parameter {
		case ("weather"):
			if recived {
				fmt.Print("Looking for new weather data...")
				json.Unmarshal(getweather(), &wh)
				recived = false
			}

			d, _ := json.Marshal(wh)
			return c.String(http.StatusAccepted, string(d))

		default:
			return c.String(http.StatusAccepted, "Please provide a value")
		}
	})

	e.Start("0.0.0.0:5000")
}

func getweather() []byte {
	resp, _ := http.Get("http://api.weatherapi.com/v1/current.json?key=REMOVEDforPROTECTION&q=Drammen&aqi=no")
	read, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return read
}
