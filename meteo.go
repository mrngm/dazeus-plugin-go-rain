package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Meteo struct {
	Astro       Astro `json:"Astro"`
	Observation Obs   `json:"Obs"`
	Texts       Texts `json:"Texts"`
}

type Astro struct {
	CurrentTimeISO string    `json:"CurrentTime_ISO"`
	CurrentTime    string    `json:"CurrentTime"`
	Datum          string    `json:"Datum"`
	DaylengthDiff  float64   `json:"DayLengthDiff"`
	Daylength      string    `json:"DayLength"`
	Moonphase      Moonphase `json:"MoonPhase"`
	Moonrise       string    `json:"MoonRise"`
	Moonset        string    `json:"MoonSet"`
	Sunrise        string    `json:"SunRise"`
	SunriseISO     string    `json:"SunRise_ISO"`
	SunsetISO      string    `json:"SunSet_ISO"`
	Sunset         string    `json:"SunSet"`
}

func (a Astro) String() string {
	return fmt.Sprintf("Zon op/neer: %s/%s (daglicht: %sm, %+0.0fm t.o.v. vorige week)",
		strings.TrimSuffix(a.Sunrise, ":00"),
		strings.TrimSuffix(a.Sunset, ":00"),
		strings.Replace(a.Daylength, ":", "h", 1),
		a.DaylengthDiff)
}

type Moonphase struct {
	FirstQuarter string `json:"EK"`
	FullMoon     string `json:"VM"`
	LastQuarter  string `json:"LK"`
	NewMoon      string `json:"NM"`
}

type Obs struct {
	WindDirection string  `json:"DDDD"`
	WindBft       float64 `json:"FFFF"`
	TempFeel      float64 `json:"FEELS_LIKE"`
	Pressure      float64 `json:"PPPP"`
	Temp          float64 `json:"TTTT"`
	RelHumidity   float64 `json:"RHRH"`
	DateTime      string  `json:"Validdt"`
	/*
	   "FFFF" : 1,
	   "FFFF_KMH" : 5,
	   "FFFF_KTS" : 3,
	   "FFFF_MPH" : 3,
	   "FFFF_MS" : 1.5,
	   "FFFX" : "4.290",
	   "FreezingLevel" : "1952.05",
	   "HAIL" : "0.00",
	   "LocationLevel" : "Dorp",
	   "NNNN" : "9.000"
	   "QHQH" : "105"
	   "RRRK" : 10,
	   "RRRR" : 0,
	   "SSSP" : null,
	   "ValidDt" : "2018-11-08T14:00:00",
	   "ValidUtc" : "2018-11-08 13:00:00",
	   "VVVV" : 10020,
	   "WeatherIconMask" : "000000000",
	   "WXCO" : "A",
	   "WXCO_EXTENDED" : "A001D",
	   "WXNUM" : "9",
	*/
}

func (o Obs) String() string {
	return fmt.Sprintf("(%s) %0.0f°C (feel: %0.0f°C), %0.0f %%RH, %.2fhPa, wind: %s kracht %0.0f",
		strings.TrimSuffix(o.DateTime[strings.Index(o.DateTime, "T")+1:], ":00"),
		o.Temp, o.TempFeel, o.RelHumidity, o.Pressure, o.WindDirection, o.WindBft)
}

type Texts struct {
	LongTerm     string `json:"LongTerm"`
	ShortTerm    string `json:"ShortTerm"`
	ForecastText string `json:"ForecastText"`
}

func GetMeteo(lat, lon int) (Meteo, error) {
	var f Meteo
	if lat == 0 {
		lat = 5181
	}
	if lon == 0 {
		lon = 584
	}
	url := fmt.Sprintf("https://api.meteoplaza.com/v2/meteo/completelocation/%d.%d?lang=nl", lat, lon)
	resp, err := http.Get(url)
	if err != nil {
		return f, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return f, err
	}

	if err := json.Unmarshal(body, &f); err != nil {
		return f, err
	}

	return f, nil
}
