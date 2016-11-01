package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Hour struct {
	Hour    string
	Minutes []string
}

type Train struct {
	Departure   time.Time
	Destination string
}

func printCommingTrains(trains []Train, filter string) {
	now := time.Now()
	destinations := strings.Split(filter, " ")
	for i := 0; i < len(trains); i++ {
		train := trains[i]
		for d := 0; d < len(destinations); d++ {
			if destinations[d] != "" && destinations[d] != train.Destination {
				continue
			}
			if train.Departure.Before(now) {
				continue
			}
			duration := train.Departure.Sub(now)
			if int(duration.Hours()) < 1 {
				str := strconv.Itoa(int(duration.Minutes())) +
					" minutes left to " + strconv.Itoa(train.Departure.Hour()) +
					":" + fmt.Sprintf("%02d", train.Departure.Minute())
				if train.Destination != "" {
					str = str + " " + train.Destination
				}
				fmt.Println(str)
			}
		}
	}
}

func timetableToTrains(timetable []Hour) (trains []Train) {
	trains = []Train{}
	for i := 0; i < len(timetable); i++ {
		hour := timetable[i]
		for j := 0; j < len(hour.Minutes); j++ {
			var destination string
			t := time.Now()
			split := strings.Split(hour.Minutes[j], " ")
			hourInt, _ := strconv.Atoi(hour.Hour)
			minuteInt, _ := strconv.Atoi(split[0])
			if len(split) > 1 {
				destination = split[1]
			}

			departure := time.Date(t.Year(), t.Month(), t.Day(), hourInt, minuteInt, 0, 0, time.Local)
			train := Train{Departure: departure, Destination: destination}
			trains = append(trains, train)
		}
	}
	return trains
}

func main() {
	var filter = flag.String("f", "", "space separated destination filter")
	flag.Parse()
	decoder := json.NewDecoder(os.Stdin)
	timetable := make([]Hour, 0)
	decoder.Decode(&timetable)
	trains := timetableToTrains(timetable)
	printCommingTrains(trains, *filter)
}
