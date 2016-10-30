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

func printLefts(nowMinute int, hour Hour, filter string) {
	for i := 0; i < len(hour.Minutes); i++ {
		result := strings.Split(hour.Minutes[i], " ")
		dests := strings.Split(filter, " ")
		minute, _ := strconv.Atoi(result[0])
		var dest string
		if len(result) > 1 {
			dest = result[1]
		}
		if dests != nil {
			for d := 0; d < len(dests); d++ {
				if dest != dests[d] {
					continue
				}
				left := minute - nowMinute
				if left > 0 && left < 60 {
					left := strconv.Itoa(left)
					take := fmt.Sprintf("%02d", minute)
					str := left + " minutes left to " + hour.Hour + ":" + take
					if dest != "" {
						str = str + " " + dest
					}
					fmt.Println(str)
				}
			}
		}
	}
}

func main() {
	var filter = flag.String("f", "筑 西 唐", "space separated destination filter")
	flag.Parse()
	decoder := json.NewDecoder(os.Stdin)
	timetable := make([]Hour, 0)
	decoder.Decode(&timetable)
	now := time.Now()
	var currentHour, nextHour Hour
	for i := 0; i < len(timetable); i++ {
		hour := timetable[i]
		oclock, _ := strconv.Atoi(hour.Hour)
		if oclock == now.Hour() {
			currentHour = hour
			nextHour = timetable[i+1]
			break
		}
	}
	printLefts(now.Minute(), currentHour, *filter)
	if now.Minute() > 44 {
		printLefts(now.Minute()-60, nextHour, *filter)
	}
}
