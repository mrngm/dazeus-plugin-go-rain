package main

import (
	"fmt"
	"strings"
)

func ShortDate(m Meteo) string {
	return strings.TrimSuffix(m.Astro.Datum, "T00:00:00")
}

func ShortDateTime(m Meteo) string {
	return strings.TrimSuffix(strings.Replace(m.Observation.DateTime, "T", " ", 1), ":00")
}

func WeatherShortTerm(m Meteo) string {
	return fmt.Sprintf("[%s] %s", ShortDate(m), m.Texts.ShortTerm)
}

func WeatherLongTerm(m Meteo) string {
	return fmt.Sprintf("[%s] %s", ShortDate(m), m.Texts.LongTerm)
}

func Moon(m Meteo) string {
	return fmt.Sprintf("[%s] Maan op/neer: %s/%s, eerste kwartier: %s, vol: %s, laatste kwartier: %s, nieuw: %s",
		ShortDate(m),
		strings.TrimSuffix(m.Astro.Moonrise, ":00"),
		strings.TrimSuffix(m.Astro.Moonset, ":00"),
		m.Astro.Moonphase.FirstQuarter,
		m.Astro.Moonphase.FullMoon,
		m.Astro.Moonphase.LastQuarter,
		m.Astro.Moonphase.NewMoon)
}

func HereComesTheSun(m Meteo) string {
	return fmt.Sprintf("[%s] %s", ShortDate(m), m.Astro)
}

func CurrentWeather(m Meteo) string {
	return fmt.Sprintf("[%s] %s", ShortDate(m), m.Observation)
}
