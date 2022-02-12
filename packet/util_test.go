package packet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlainLatLng(t *testing.T) {

	// Parse in string
	// [0.1] KD2AOY-8>APDR10,BKELEY,WIDE1*:!3744.94N/12226.47W[From Anytone D878UV
	// N 37 44.9400, W 122 26.4700
	input := []byte("3744.94N/12226.47W")
	expectedLat := 37.749
	expectedLng := -122.44116666666666
	lat, lng := LocationDataToLatLng(input)
	// sigh
	computedLat, _ := lat.Float64()
	computedLng, _ := lng.Float64()
	assert.Equal(t, computedLat, expectedLat, "The parsed latitude should match expected latitude")
	// Flaky toString() nonsense :(
	assert.Equal(t, computedLng, expectedLng, "The parsed longitude should match expected longitude")
}

func TestParseMicE(t *testing.T) {
	assert := assert.New(t)
	//2020-06-07 06:51:05 PDT: KF6AXA>S8TXPR,WIDE1-1,WIDE2-1,qAR,GAPGTE:'10 l <0x1c>-/]=
	// AGWPE Packet: {AGWPEPort:0 DataKind:75 PID:0 CallFrom:KF6AXA^ CallTo:S8TXPR DatLen:50 Data:[166 112 168 176 160 164 96 150 140 108 130 17
	// 	6 130 224 174 108 134 176 64 64 230 132 150 138 152 138 178 224 174 146 136 138 100 64 225 3 240 39 49 48 32 108 32 28 45 47 93 61] RawPacket:AAAAAE
	// 	sAAABLRjZBWEEAAAAAUzhUWFBSAAAAADIAAAAAAAAAAKZwqLCgpGCWjGyCsILgrmyGsEBA5oSWipiKsuCukoiKZEDhA/AnMTAgbCAcLS9dPQ==}
	pck, err := AGWPEPacketFromB64("AAAAAEsAAABLRjZBWEEAAAAAUzhUWFBSAAAAADIAAAAAAAAAAKZwqLCgpGCWjGyCsILgrmyGsEBA5oSWipiKsuCukoiKZEDhA/AnMTAgbCAcLS9dPQ==")
	assert.NoError(err)
	chkin, err := ParseAX25Frame(pck.Data)
	assert.NoError(err)
	assert.Equal("KF6AXA", chkin.Source)
	assert.Equal("S8TXPR", chkin.Dest)
	// Probably dubious
	//assert.Equal(chkin.GetLocationCheckin().Location.Latitude, 38.800333333333334)
	//assert.Equal(chkin.GetLocationCheckin().Location.Longitude, -121.334)
}

func TestParseValidMicE(t *testing.T) {
	// AGWPE Packet: {AGWPEPort:0 DataKind:75 PID:0 CallFrom:KE6CAC^@^@^@ CallTo:S8RTYT^@^@^@ DatLen:49 Data:[166 112 164 168 178 168 96 150 138 108 134 13
	// 0 134 224 174 108 134 176 64 64 230 132 150 138 152 138 178 224 156 96 144 170 160 64 225 3 240 39 49 55 88 108 32 28 45 47 93] RawPacket:AAAAAEsAAA
	// BLRTZDQUMAAAAAUzhSVFlUAAAAADEAAAAAAAAAAKZwpKiyqGCWimyGgobgrmyGsEBA5oSWipiKsuCcYJCqoEDhA/AnMTdYbCAcLS9d}
	assert := assert.New(t)
	pkt, err := AGWPEPacketFromB64("AAAAAEsAAABLRTZDQUMAAAAAUzhSVFlUAAAAADEAAAAAAAAAAKZwpKiyqGCWimyGgobgrmyGsEBA5oSWipiKsuCcYJCqoEDhA/AnMTdYbCAcLS9d")
	assert.NoError(err)
	checkin, err := ParseAX25Frame(pkt.Data)
	assert.NoError(err)
	assert.Equal("KE6CAC", checkin.GetSource())
	assert.Equal("S8RTYT", checkin.GetDest())
	// Potentially dubious
	//assert.Equal(checkin.GetLocationCheckin().Location.Latitude, 38.41566666666667)
	//assert.Equal(checkin.GetLocationCheckin().Location.Longitude, -121.46)
}
func TestDMH(t *testing.T) {
	assert := assert.New(t)
	capitolLng := DMHToDecimal(38.0, 53.0, 23.0)
	capLng, _ := capitolLng.Float64()
	assert.Equal(ToFixed(capLng, 4), 38.8897)
	capitolLat := DMHToDecimal(77, 00, 32)
	capLat, _ := capitolLat.Float64()
	assert.Equal(ToFixed(capLat, 4), 77.0089)

}
