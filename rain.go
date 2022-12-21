package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/dazeus/dazeus-go"
)

var myCommand string

func main() {
	connStr := "unix:/tmp/dazeus.sock"
	if len(os.Args) > 1 {
		connStr = os.Args[1]
	}
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("Paniek! %v\n", p)
			debug.PrintStack()
		}
	}()

	dz, err := dazeus.ConnectWithLoggingToStdErr(connStr)
	if err != nil {
		panic(err)
	}

	if _, hlerr := dz.HighlightCharacter(); hlerr != nil {
		panic(hlerr)
	}

	_, err = dz.SubscribeCommand("maan", dazeus.NewUniversalScope(), func(ev dazeus.Event) {
		f, err := GetMeteo(0, 0)
		if err != nil {
			ev.Reply(fmt.Sprintf("E_MAAN: %v", err), true)
			return
		}
		ev.Reply(Moon(f), true)
	})
	if err != nil {
		panic(err)
	}
	_, err = dz.SubscribeCommand("zon", dazeus.NewUniversalScope(), func(ev dazeus.Event) {
		f, err := GetMeteo(0, 0)
		if err != nil {
			ev.Reply(fmt.Sprintf("E_ZON: %v", err), true)
			return
		}
		ev.Reply(HereComesTheSun(f), true)
	})
	if err != nil {
		panic(err)
	}
	_, err = dz.SubscribeCommand("kortweer", dazeus.NewUniversalScope(), func(ev dazeus.Event) {
		f, err := GetMeteo(0, 0)
		if err != nil {
			ev.Reply(fmt.Sprintf("E_KORTWEER: %v", err), true)
			return
		}
		ev.Reply(WeatherShortTerm(f), true)
	})
	if err != nil {
		panic(err)
	}
	_, err = dz.SubscribeCommand("langweer", dazeus.NewUniversalScope(), func(ev dazeus.Event) {
		f, err := GetMeteo(0, 0)
		if err != nil {
			ev.Reply(fmt.Sprintf("E_FRIESLAND: %v", err), true)
			return
		}
		ev.Reply(WeatherLongTerm(f), true)
	})
	if err != nil {
		panic(err)
	}
	_, err = dz.SubscribeCommand("weer", dazeus.NewUniversalScope(), func(ev dazeus.Event) {
		fmt.Printf("params: %+v\n", ev.Params)
		f, err := GetMeteo(0, 0)
		if err != nil {
			ev.Reply(fmt.Sprintf("E_OUD: %v", err), true)
			return
		}
		ev.Reply(CurrentWeather(f), true)
	})
	if err != nil {
		panic(err)
	}
	_, err = dz.SubscribeCommand("regen", dazeus.NewUniversalScope(), func(ev dazeus.Event) {
		b, err := GetBuien()
		if err != nil {
			ev.Reply(fmt.Sprintf("E_DRUPPEL: %v", err), true)
			return
		}
		ev.Reply(BuienForecast(b), true)
	})
	_, err = dz.SubscribeCommand("korterweer", dazeus.NewUniversalScope(), func(ev dazeus.Event) {
		f, err := GetMeteo(0, 0)
		if err != nil {
			ev.Reply(fmt.Sprintf("E_MAND: %v", err), true)
			return
		}
		if f.Observation.TempFeel > 30 {
			ev.Reply("HEET!", true)
			return
		}
		if f.Observation.TempFeel > 22 {
			ev.Reply("WARM!", true)
			return
		}
		if f.Observation.TempFeel > 15 {
			ev.Reply("METZONDERJAS!", true)
			return
		}
		if f.Observation.TempFeel > 8 {
			ev.Reply("FRIS!", true)
			return
		}
		if f.Observation.TempFeel > 0 {
			ev.Reply("MUTS!", true)
			return
		}
		ev.Reply("FRIES!", true)
	})

	listenerr := dz.Listen()
	panic(listenerr)
}
