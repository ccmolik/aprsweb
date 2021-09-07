package packet

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"time"

	pb "github.com/ccmolik/aprsweb-proto/rpc/aprsweb"
	"github.com/shopspring/decimal"
)

// An AGWPEPacket contains low-level data about a packet.
type AGWPEPacket struct {
	AGWPEPort byte   // Port where the data frame had been received thru {0=Port1,1=Port2,…}
	DataKind  byte   // The kind of data
	PID       byte   // Usually not used, it's the frame byte
	CallFrom  string // Callsign FROM of packet ASCII 10 bytes, null terminated
	CallTo    string // Callsign TO of packet
	DatLen    uint32 // The length of the data
	Data      []byte
	RawPacket string // The actual packet itself in base64
}

// AGWPEPacketFromB64 turns a base64 string into an AGWPE Packet.
// Useful for debugging / tests.
func AGWPEPacketFromB64(b64 string) (*AGWPEPacket, error) {
	buf, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	datLen := binary.LittleEndian.Uint32(buf[28:32])
	l := AGWPEPacket{
		AGWPEPort: buf[1],
		DataKind:  buf[4],
		CallFrom:  string(buf[8:17]),
		CallTo:    string(buf[18:27]),
		DatLen:    datLen,
		Data:      buf[37:],
		RawPacket: b64,
	}
	return &l, nil
}

// LatLng represents Latitude and Longitude
type LatLng struct {
	Latitude  decimal.Decimal
	Longitude decimal.Decimal
}

// ParseAX25Frame converts raw bytes into a checkin, or returns an error
func ParseAX25Frame(inputBytes []byte) (*pb.Checkin, error) {
	retCheckin := pb.Checkin{
		ReceivedDatetime: time.Now().UnixNano(),
	}

	destCallsign, err := parseAX25Callsign(inputBytes[0:7])
	if err != nil {
		return &retCheckin, err
	}
	retCheckin.Dest = destCallsign
	// parse SSID
	s, err := parseAX25Callsign(inputBytes[7:14])
	if err != nil {
		return &retCheckin, err
	}
	retCheckin.Source = s
	// TODO Parse checkin repated or not
	// retCheckin.IsRepeated = inputBytes[14]&0x80 != 0

	// if inputBytes[13]&0x01 == 0 {
	// 	retCheckin.IsRepeated = true
	// }

	APRSData, err := getAPRSData(inputBytes)
	if len(APRSData) == 0 {
		return &retCheckin, errors.New("Got a zero-length packet back from getAPRSData.. this happens I guess")
	}
	// Determine if we have MIC-E madness
	switch APRSData[0] {
	// Location Checkins
	// Mic-E:
	// 0x60 = ` ; 0x1c is current rev 0- beta, 0x1d is old rev0 beta,  0x27 = '
	// Non-Mic-E:
	// 0x2f (/) is location data with timestamp (without message)
	// 0x3d (=) is location without timestamp (with message)
	// 0x40 (@) is location with timestamp (with message)
	// 0x21 (!) is location without timestamp (no messaging) or a weather station?????

	case 0x60, 0x1c, 0x1d, 0x27, 0x2f, 0x3d, 0x40, 0x21:
		var checkinMsg pb.LocationCheckin
		retCheckin.Type = pb.CheckinType_POSITION
		// MIC-E data, do things with it
		if APRSData[0] == 0x60 || APRSData[0] == 0x1c || APRSData[0] == 0x1d || APRSData[0] == 0x27 {
			checkinMsg.IsMicE = true
			loc := parseMicE(&retCheckin, APRSData) // D:
			locLat, _ := loc.Latitude.Float64()
			locLng, _ := loc.Longitude.Float64()

			checkinMsg.Location = &pb.LatLng{
				Latitude:  locLat,
				Longitude: locLng,
			}
			// yes, these are reversed
			// this is why i drink
			retCheckin.MapSymbol = string([]byte{APRSData[7]})
			retCheckin.SymbolTable = string([]byte{APRSData[8]})

		}
		if APRSData[0] == 0x2f || APRSData[0] == 0x40 || APRSData[0] == 0x3d || APRSData[0] == 0x21 {
			checkinMsg.IsMicE = false
			locPayloadStart := 0
			found := false
			if APRSData[0] == 0x2f || APRSData[0] == 0x40 {
				// we have a timestamp
				// TODO parse timestamp
				locPayloadStart = 7
			} else {
				// no timestamp 4 us
				for offset, data := range APRSData {
					if data == 0x21 || data == 0x3D {
						locPayloadStart = offset
						found = true
						break
					}
				}
				if !found {
					// Didn't find a location data i guess
					//log.Printf("Did not find Offset Data: %+v", hex.Dump(APRSData))
					return &retCheckin, errors.New("Did not find Location Data offset in frame")
				}
			}
			// We have an offset for the actual loc data now
			lat, lng := LocationDataToLatLng(APRSData[locPayloadStart+1:])
			// Set the table selector and symbol (fixed offset)
			retCheckin.SymbolTable = string([]byte{APRSData[locPayloadStart+9]})
			retCheckin.MapSymbol = string([]byte{APRSData[locPayloadStart+19]})
			// TODO this ignores precision exactness in float conversion
			locLat, _ := lat.Float64()
			locLng, _ := lng.Float64()
			// TODO implement reported time
			checkinMsg.Location = &pb.LatLng{
				Latitude:  locLat,
				Longitude: locLng,
			}
		}
		retCheckin.Checkin = &pb.Checkin_LocationCheckin{
			LocationCheckin: &checkinMsg,
		}
	default:
		retCheckin.Type = pb.CheckinType_OTHER
	}
	return &retCheckin, nil
}

// parseAX25Callsign parses callsign bytes in network order and returns a string representing the callsign with optional SSID
func parseAX25Callsign(bytes []byte) (string, error) {
	// Callsigns are 6 bytes long
	if len(bytes) != 7 {
		return "", fmt.Errorf("Length of packet is %v bytes; expected 6", len(bytes))
	}
	var tmp []byte
	inPadSpaces := false
	for _, n := range bytes[0:6] {
		a := (n & 0xFF) >> 1
		// log.Printf("n: %x-a: %x,", n, a)
		if a == 0x20 {
			inPadSpaces = true
		}
		if !inPadSpaces {
			tmp = append(tmp, a)
		}
	}
	ssid := (bytes[6] & 0x1E) >> 1
	if ssid != 0 {
		return fmt.Sprintf("%s-%s", string(tmp), strconv.Itoa(int(ssid))), nil
	}
	return string(tmp), nil

}

func parseMicE(checkin *pb.Checkin, frame []byte) LatLng {
	latString := ""
	isLngOffset := false
	isLngWest := false
	isLatNorth := true
	for i, LatByte := range []byte(checkin.Dest) {
		// First, get the decimal represented in the letter
		switch {
		// 0-9
		case 0x30 <= LatByte && LatByte <= 0x39:
			latString += strconv.Itoa(int(LatByte - 0x30))
		// A-J
		case 0x41 <= LatByte && LatByte <= 0x4a:
			latString += strconv.Itoa(int(LatByte - 0x41))
		// K L or Z = " "
		case LatByte == 0x4b || LatByte == 0x4c || LatByte == 0x5a:
			latString += " "
		// P-Y
		case 0x50 <= LatByte && LatByte <= 0x59:
			latString += strconv.Itoa(int(LatByte - 0x50))
		default:
			// Couldn't decode.
			// TODO be less opaque
			return LatLng{}
		}
		// Now handle the nonsense
		switch i {
		// TODO: implement ABC codes
		// Handle N/S
		case 3:
			if 0x50 > LatByte {
				isLatNorth = false
			}

		// Handle Long offset (if above "P") on the ascii table
		case 4:
			if 0x50 <= LatByte {
				isLngOffset = true
			}
			// Handle West vs East
		case 5:
			if 0x50 <= LatByte {
				isLngWest = true
			}
		}
	}
	// Insert a decimal point for formatting
	// latString = latString[:2] + "." + latString[2:]

	latitudeDegrees, _ := strconv.ParseFloat(latString[:2], 64)
	latitudeMinutes, _ := strconv.ParseFloat(latString[2:4], 64)
	latitudeHundredthsOfMinutes, _ := strconv.ParseFloat(latString[4:], 64)
	computedLat := DMHToDecimal(latitudeDegrees, latitudeMinutes, latitudeHundredthsOfMinutes)
	if !isLatNorth {
		computedLat = computedLat.Neg()
	}

	longitudeDegrees := float64(frame[1]) - 28
	if isLngOffset {
		longitudeDegrees += 100
	}
	if longitudeDegrees >= 180 {
		if longitudeDegrees <= 189 {
			longitudeDegrees -= 80
		} else {
			longitudeDegrees -= 190
		}
	}
	longtitudeMinutes := (float64)(frame[2]) - 28
	if longtitudeMinutes >= 60 {
		longtitudeMinutes -= 60
	}
	lngHundredthsOfMinutes := (float64)(frame[3]) - 28
	if lngHundredthsOfMinutes > 60 {
		lngHundredthsOfMinutes -= 60
	}
	// log.Printf("Degress %v, Minutes %v, Hundredths %v", longitudeDegrees, longtitudeMinutes, lngHundredthsOfMinutes)
	computedLongitude := DMHToDecimal(longitudeDegrees, longtitudeMinutes, lngHundredthsOfMinutes)
	if isLngWest {
		computedLongitude = computedLongitude.Neg()
	}
	return LatLng{
		Latitude:  computedLat,
		Longitude: computedLongitude,
	}
}

func getAPRSData(inputBytes []byte) ([]byte, error) {
	payloadStart := 0
	// Disregarding the list of repeaters, look for the data
	for offset, data := range inputBytes[15:] {
		if offset == len(inputBytes)-15 {
			return nil, errors.New("Couldn't find data in this packet")
		}
		if data == 0x03 && inputBytes[15+offset+1] == 0xF0 {
			payloadStart = 15 + offset + 2
			break
		}
	}
	return inputBytes[payloadStart:], nil
}
