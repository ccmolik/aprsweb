package packet

import (
	"strconv"

	"github.com/shopspring/decimal"
)

// DMHToDecimal converts APRS Degree Minutes HundredthsOfMinutes spec to a decimal representation
func DMHToDecimal(degrees float64, minutes float64, seconds float64) decimal.Decimal {
	sixtyDec := decimal.NewFromFloat(float64(60))
	// hundredDec := decimal.NewFromFloat(float64(100))
	thirtysixhundredDec := decimal.NewFromFloat(float64(3600))

	decDegrees := decimal.NewFromFloat(degrees)
	decMinutes := decimal.NewFromFloat(minutes)
	decHundredths := decimal.NewFromFloat(seconds)
	return decDegrees.
		Add(decMinutes.Div(sixtyDec)).
		Add(decHundredths.Div(thirtysixhundredDec))
}

// LocationDataToLatLng converts an APRS location data string without the header into lat, lng in decimal
func LocationDataToLatLng(locationString []byte) (decimal.Decimal, decimal.Decimal) {
	// Replace every space (0x20) with 0
	for i := range locationString {
		if locationString[i] == 0x20 {
			locationString[i] = 0x30
		}
	}
	// 4903.50N/07201.75W-Test 001234
	// :!3744.95N/12226.47W
	latDegrees, _ := strconv.ParseFloat(string(locationString[0:2]), 64)
	latMinutes, _ := strconv.ParseFloat(string(locationString[2:4]), 64)
	latHundredths, _ := strconv.ParseFloat(string(locationString[5:7]), 64)
	latitude := DMHToDecimal(latDegrees, latMinutes+(latHundredths/100), 0)

	if locationString[7] == 0x53 {
		latitude = latitude.Neg()
	}
	// fmt.Println(hex.Dump(locationString))

	// Skip the slash so we start at 9
	lngDegrees, _ := strconv.ParseFloat(string(locationString[9:12]), 64)
	lngMinutes, _ := strconv.ParseFloat(string(locationString[12:14]), 64)
	lngHundredths, _ := strconv.ParseFloat(string(locationString[15:17]), 64)

	longitude := DMHToDecimal(lngDegrees, lngMinutes+(lngHundredths/100), 0)
	if locationString[17] == 0x57 {
		longitude = longitude.Neg()
	}
	return latitude, longitude
}
