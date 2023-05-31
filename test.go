package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
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

type Test struct {
	start string
	end   string
	check string
}

func main() {

	//var off_peak_cost float32
	//var peak_cost float32
	var t Test
	var truc string
	var watts, final float64

	//hcreuses array prices order are off-peak hours then peak hours
	//tempo array prices order are : off-peak blue, white, red prices then peak blue, white, red
	//hcreuses_hours have two intervals. order is start of the interval, then end, outside this interval, is peak hours
	//tempo_hours have only one interval, wich is the offpeak hours, outside this interval is peak hours

	/*
		TODO: prices from edf doesn't include electricity transport cost, wich is almost
		hidden in your annual bill, and make all cost calculation wrong if you don't add
		it
		So, the real todo is to add a flexible value with those transport electricity cost and TVA
	*/
	prices_hcreuses := [2]float64{0.2006, 0.2742}
	//prices_tempo := [6]float32{0.0970, 0.1140, 0.1216, 0.1249, 0.1508, 0.6712}
	//price_baseprice = 0.2062
	hcreuses_hours := [4]string{"01:00", "07:00", "12:30", "13:30"}
	//tempo_hours := [2]string{"22:00", "06:00"}

	//Todo : array of dates of peak red tempo prices, as they are on specific days in addition of specific times
	filepath := "Enedis_Conso_Heure_20230522-20230523_19119536803609.csv"
	result := readCsvFile(filepath)
	newLayout := "15:04"
	for i := 3; i < len(result); i++ {
		datetime := result[i][0]
		datetime, truc, _ = strings.Cut(datetime, ";")
		watts, _ = strconv.ParseFloat(truc, 64)
		tm, _ := time.Parse(time.RFC3339, datetime)
		mytime := (tm.Format("15:04"))
		t.start, t.end, t.check = hcreuses_hours[0], hcreuses_hours[1], mytime
		check, _ := time.Parse(newLayout, t.check)
		start, _ := time.Parse(newLayout, t.start)
		end, _ := time.Parse(newLayout, t.end)
		if inTimeSpan(start, end, check) {
			KwH := watts * 0.5 / 1000
			KwH = roundFloat(KwH, 4)

			KwH = KwH * prices_hcreuses[0]
			KwH = roundFloat(KwH, 4)
			values = append(values, KwH)
		} else {
			t.start, t.end, t.check = hcreuses_hours[2], hcreuses_hours[3], mytime
			check, _ := time.Parse(newLayout, t.check)
			start, _ := time.Parse(newLayout, t.start)
			end, _ := time.Parse(newLayout, t.end)
			if inTimeSpan(start, end, check) {
				KwH := watts * 0.5 / 1000
				KwH = roundFloat(KwH, 4)
				KwH = KwH * prices_hcreuses[0]
				KwH = roundFloat(KwH, 4)
				values = append(values, KwH)
			} else {
				KwH := watts * 0.5 / 1000
				KwH = roundFloat(KwH, 4)
				KwH = KwH * prices_hcreuses[1]
				KwH = roundFloat(KwH, 4)
				values = append(values, KwH)
			}
		}
		t.start, t.end = hcreuses_hours[2], hcreuses_hours[3]

	}
	for i := 0; i < len(values); i++ {
		final = final + values[i]
	}

	//Todo : add a debug function

	final = roundFloat(final, 2)
	fmt.Println(final)

}
