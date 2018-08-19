package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/djherbis/times"
)

type Hour struct {
	Hour    string
	Minutes []string
}

type Train struct {
	Departure   time.Time
	Destination string
}

var path string

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

func getPath(targetUrl string) string {
	u, err := url.Parse(targetUrl)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("/tmp/%s%s%s.html", u.Host, strings.Replace(u.Path, "/", "-", -1), url.QueryEscape(u.RawQuery))
}

func getPage(targetUrl string) bool {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	t, err := times.Stat(path)
	if err != nil {
		panic(err)
	}

	mtime := t.ModTime()
	anHourAgo := time.Now().Add(-1 * time.Hour)

	filestat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	if mtime.After(anHourAgo) && filestat.Size() > 0 {
		return true
	}

	response, err := http.Get(targetUrl)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		panic(err)
	}
	return true
}

func parseHtml(path string) []Hour {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	buf := make([]byte, 65000)
	for {
		n, err := file.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			break
		}
	}

	timetable := []Hour{}

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(buf)))
	doc.Find("table.tblDiaDetail").Each(func(_ int, table *goquery.Selection) {
		table.Find("tr").Each(func(_ int, tr *goquery.Selection) {
			hour := Hour{}
			var oclock string
			tr.Find("td.hour").Each(func(_ int, td *goquery.Selection) {
				oclock = td.Text()
			})
			hour.Hour = oclock
			var minutes []string
			tr.Find("ul li.timeNumb dl").Each(func(_ int, dl *goquery.Selection) {
				var minute string
				dl.Find("dt").Each(func(_ int, dt *goquery.Selection) {
					minute = dt.Text()
				})
				dl.Find("dd.trainFor").Each(func(_ int, dd *goquery.Selection) {
					destination := dd.Text()
					minute = fmt.Sprintf("%s %s", minute, destination)
				})
				minutes = append(minutes, minute)
			})
			hour.Minutes = minutes
			if oclock != "" {
				timetable = append(timetable, hour)
			}
		})
	})

	return timetable
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
	var targetUrl = flag.String("u", "", "timetable targetUrl")
	flag.Parse()
	path = getPath(*targetUrl)
	getPage(*targetUrl)
	timetable := parseHtml(path)
	trains := timetableToTrains(timetable)
	printCommingTrains(trains, *filter)
}
