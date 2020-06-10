package packet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestParseMicE(t *testing.T) {
// 	p := AX25Frame{
// 		DestAddr: "S32U6T",
// 		IsMICE:   true,
// 	}
// 	fmt.Printf("%s", parseMicE(p))
// }

/*
May 17 22:56:29 kd2aoy-pi sh[307]: Digipeater WIDE1 (probably BKELEY) audio level = 98(22/20)   [NONE]   |||||____
May 17 22:56:29 kd2aoy-pi sh[307]: [0.2] AG6JA>APOTU0,BKELEY,WIDE1*,WIDE2-2:!/;T>+/Z7ckxCG
May 17 22:56:29 kd2aoy-pi sh[307]: Position, truck, Open Track
May 17 22:56:29 kd2aoy-pi sh[307]: N 37 26.9741, W 122 07.1746, 15 MPH, course 348
May 17 22:56:29 kd2aoy-pi sh[307]: Error getting message header from AGW client application 0.
*/

/*
May 17 23:15:44 kd2aoy-pi sh[307]: [0.2] AG6JA>APOTU0,BKELEY,WIDE1*,WIDE2-2:!/;T/2/ZC7k1?G
May 17 23:15:44 kd2aoy-pi sh[307]: Position, truck, Open Track
May 17 23:15:44 kd2aoy-pi sh[307]: N 37 27.1880, W 122 06.8445, 10 MPH, course 64
panic: runtime error: index out of range [17] with length 13
*/

func TestBreakerWithRange(t *testing.T) {
	assert := assert.New(t)

	p, err := AGWPEPacketFromB64("AAAAAEsAAABBRzZKQQAAAAAAQVBPVFUwAAAAADQAAAAAAAAAAIKgnqiqYOCCjmyUgkDghJaKmIqy4K6SiIpiQOCukoiKZEBlA/AhLztULzIvWkM3azE/Rz4vJyI0IX18IT8lZihPfCF3dS0hfDMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==")
	assert.NoError(err)
	assert.NotNil(p)
	// Attempt to do something useful with it
	checkin, err := ParseAX25Frame(p.Data)
	assert.NoError(err)
	assert.NotNil(checkin)

}

// AGWPE Packet: {AGWPEPort:0 DataKind:75 PID:0 CallFrom:K6IXA-3 CallTo:APIN20 DatLen:108 Data:[130 160 146 156 100 96 96 150 108 146 176 130 64 102 174 108 134 176 64 64 230 132 150 138 152 138 178 224 174 146 136 138 100 64 225 3 240 64 48 55 49 51 51 57 122 51 54 52 52 46 49 51 78 47 49 49 56 53 49 46 53 53 87 95 50 50 57 47 48 48 52 103 48 48 57 116 48 51 52 114 48 48 48 112 48 48 49 80 48 48 49 104 52 53 98 49 48 49 54 57 46 68 115 86 80 13] RawPacket:AAAAAEsAAABLNklYQS0zAAAAQVBJTjIwAAAAAGwAAAAAAAAAAIKgkpxkYGCWbJKwgkBmrmyGsEBA5oSWipiKsuCukoiKZEDhA/BAMDcxMzM5ejM2NDQuMTNOLzExODUxLjU1V18yMjkvMDA0ZzAwOXQwMzRyMDAwcDAwMVAwMDFoNDViMTAxNjkuRHNWUA0=}
// Jun 07 14:52:30 kd2aoy-pi sh[309]: Digipeater WIDE2 (probably BKELEY) audio level = 79(17/22)   [NONE]   _||||____
// Jun 07 14:52:30 kd2aoy-pi sh[309]: [0.2] K6IXA-3>APIN20,W6CX-3,BKELEY,WIDE2*:@071339z3644.13N/11851.55W_229/004g009t034r000p001P001h45b10169.DsVP<0x0d><0x0a>
// Jun 07 14:52:30 kd2aoy-pi sh[309]: Weather Report, WEATHER Station (blue)
// Jun 07 14:52:30 kd2aoy-pi sh[309]: N 36 44.1300, W 118 51.5500
// Jun 07 14:52:30 kd2aoy-pi sh[309]: wind 4.6 mph, direction 229, gust 9, temperature 34, rain 0.00 in last hour, rain 0.01 in last 24 hours, rain 0.01 since midnight, humidity 45, barometer 30.03, ".DsVP"

func TestTimeStampLocation(t *testing.T) {
	assert := assert.New(t)
	p, err := AGWPEPacketFromB64("AAAAAEsAAABLNklYQS0zAAAAQVBJTjIwAAAAAGwAAAAAAAAAAIKgkpxkYGCWbJKwgkBmrmyGsEBA5oSWipiKsuCukoiKZEDhA/BAMDcxMzM5ejM2NDQuMTNOLzExODUxLjU1V18yMjkvMDA0ZzAwOXQwMzRyMDAwcDAwMVAwMDFoNDViMTAxNjkuRHNWUA0=")
	assert.NoError(err)
	assert.NotNil(p)
	// Attempt to do something useful with it
	checkin, err := ParseAX25Frame(p.Data)
	assert.Equal("K6IXA-3", checkin.GetSource())
	assert.Equal("APIN20", checkin.GetDest())
	assert.Equal(36.7355, checkin.GetLocationCheckin().Location.Latitude)
	assert.Equal(-118.85916666666667, checkin.GetLocationCheckin().Location.Longitude)
	assert.NoError(err)
	assert.NotNil(checkin)

}
