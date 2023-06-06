package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
func strToTimeObject(newLayout string, start, end, check string) (time.Time, time.Time, time.Time) {
	tcheck, _ := time.Parse(newLayout, check)
	tstart, _ := time.Parse(newLayout, start)
	tend, _ := time.Parse(newLayout, end)
	return tstart, tend, tcheck
}
func priceCalculator(watts float64, price float64) float64 {
	KwH := watts * 0.5 / 1000
	KwH = roundFloat(KwH, 4)
	KwH = KwH * price
	return roundFloat(KwH, 4)
	//return KwH

}
func inTimeSpan(start, end, check time.Time) bool {
	if start.Before(end) {
		return !check.Before(start) && !check.After(end)
	}
	if start.Equal(end) {
		return check.Equal(start)
	}
	return !start.After(check) || !end.Before(check)
}

var values []float64
var basepricesvalues []float64
var red_days_map = make(map[string]bool)
var white_days_map = make(map[string]bool)
var blue_days_map = make(map[string]bool)

type Test struct {
	start string
	end   string
	check string
}

func textfileread(textfilepath string) []string {
	textfile := []string{}
	file, err := os.Open(textfilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		textfile = append(textfile, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return textfile
}

func main() {

	//var off_peak_cost float32
	//var peak_cost float32
	var t Test
	var truc string
	var watts, hphcprice, baseprices float64
	textfilepath := [2]string{"red_days", "white_days"}
	reds_days_tempo := textfileread(textfilepath[0])
	white_days_tempo := textfileread(textfilepath[1])
	fmt.Println(reds_days_tempo)
	fmt.Println(white_days_tempo)
	//hcreuses array prices order are off-peak hours then peak hours
	//tempo array prices order are : off-peak blue, white, red prices then peak blue, white, red
	//hcreuses_hours have two intervals. order is start of the interval, then end, outside this interval, is peak hours
	//tempo_hours have only one interval, wich is the offpeak hours, outside this interval is peak hours

	/*
		TODO : i'm still not sure why, but despite EDF telling the prices are Tax included, if i don't add tak again,
		prices are wrong. still need to find a way to add flexible tax rate
	*/
	prices_hcreuses := [2]float64{0.2006, 0.2742}
	//prices_tempo := [6]float32{0.0970, 0.1140, 0.1216, 0.1249, 0.1508, 0.6712}
	price_baseprice := 0.2062
	hcreuses_hours := [4]string{"01:00", "07:00", "12:30", "13:30"}
	//tempo_hours := [2]string{"22:00", "06:00"}
	for b := range reds_days_tempo {
		red_days_map[reds_days_tempo[b]] = true
	}
	for b := range white_days_tempo {
		white_days_map[white_days_tempo[b]] = true
	}
	fmt.Println(red_days_map)
	//Todo : array of dates of peak red tempo prices, as they are on specific days in addition of specific times
	filepath := "Enedis_Conso_Heure_20220503-20230530_19119536803609.csv"
	result := readCsvFile(filepath)
	newLayout := "15:04"
	for i := 3; i < len(result); i++ {
		datetime := result[i][0]
		datetime, truc, _ = strings.Cut(datetime, ";")
		watts, _ = strconv.ParseFloat(truc, 64)
		tm, _ := time.Parse(time.RFC3339, datetime)
		mytime := (tm.Format("15:04"))
		start, end, check := strToTimeObject(newLayout, hcreuses_hours[0], hcreuses_hours[1], mytime)
		if inTimeSpan(start, end, check) {
			KwH := priceCalculator(watts, prices_hcreuses[0])
			values = append(values, KwH)
		} else {
			t.start, t.end, t.check = hcreuses_hours[2], hcreuses_hours[3], mytime
			start, end, check := strToTimeObject(newLayout, hcreuses_hours[0], hcreuses_hours[1], mytime)
			if inTimeSpan(start, end, check) {
				KwH := priceCalculator(watts, prices_hcreuses[0])
				values = append(values, KwH)
			} else {
				KwH := priceCalculator(watts, prices_hcreuses[1])
				values = append(values, KwH)
			}
			t.start, t.end = hcreuses_hours[2], hcreuses_hours[3]
		}
		KwH := priceCalculator(watts, price_baseprice)
		//fmt.Println(KwH)
		basepricesvalues = append(basepricesvalues, KwH)
	}
	//fmt.Println(basepricesvalues)
	for i := 0; i < len(values); i++ {
		hphcprice = hphcprice + values[i]
	}
	for i := 0; i < len(basepricesvalues); i++ {
		baseprices = baseprices + basepricesvalues[i]
	}

	//Todo : add a debug function

	hphcprice = roundFloat(hphcprice, 2)
	fmt.Println(hphcprice)
	baseprices = roundFloat(baseprices, 2)
	fmt.Println(baseprices)

}
