package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Levels struct {
	Heavy    float64 `json:"Heavy"`
	Light    float64 `json:"Light"`
	Moderate float64 `json:"moderate"`
}

type Buien struct {
	Delta       float64   `json:"delta"`
	HumanStart  string    `json:"start_human"`
	Levels      Levels    `json:"levels"`
	Rain        []float64 `json:"precip"`
	Start       float64   `json:"Start"`
	Temperature float64   `json:"temp"`
	/*
	   "bounds" : {
	      "E" : 10.856429,
	      "N" : 55.973602,
	      "W" : 0,
	      "S" : 48.895302
	   },
	   "grid" : {
	      "x" : 419,
	      "y" : 456
	   },
	   "success" : true,
	   "source" : "nl",
	*/
}

func BuiNiveau(l float64, levelSpec Levels) string {
	switch {
	case l > levelSpec.Heavy:
		return "*"
	case l > levelSpec.Moderate:
		return "-"
	case l > levelSpec.Light:
		return "."
	default:
		return "_"
	}
}

func BuienForecast(b Buien) string {
	intensity := ""
	for _, l := range b.Rain {
		intensity = fmt.Sprintf("%s%s", intensity, BuiNiveau(l, b.Levels))
	}
	d, _ := time.ParseDuration(fmt.Sprintf("%.0fs", b.Delta))
	t, _ := time.Parse("15:04", b.HumanStart)
	hs, _ := time.ParseDuration(fmt.Sprintf("%.0fs", d.Seconds()*float64(len(b.Rain))))
	humanStop := t.Add(hs).Format("15:04")
	return fmt.Sprintf("Vanaf %s, per %.0fmin: %s %s {_ <= %.1f, . > %.1f, - > %.1f, * > %.1f}mm/u",
		b.HumanStart, d.Minutes(), intensity, humanStop, b.Levels.Light, b.Levels.Light, b.Levels.Moderate, b.Levels.Heavy)
}

func GetBuien() (Buien, error) {
	var b Buien
	url := "https://cdn-secure.buienalarm.nl/api/3.4/forecast.php?lat=51.8125626&lon=5.837226399999963&region=nl&unit=mm/u"
	resp, err := http.Get(url)
	if err != nil {
		return b, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return b, err
	}

	if err := json.Unmarshal(body, &b); err != nil {
		return b, err
	}

	return b, nil
}
