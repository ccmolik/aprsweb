package packethandler

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/ccmolik/aprsweb-proto/rpc/aprsweb"
	"github.com/ccmolik/aprsweb/packet"
)

// A PacketHandler handles AGWPE packets from the source and spits parsed frames to the store
type PacketHandler struct {
	Server          string
	Port            int32
	PacketStoreChan chan *pb.Checkin
}

// Tell Direwolf to speak to us in Raw Mode
const rawMode = "\x00" +
	"\x00\x00\x00" +
	"\x6b" +
	"\x00\x00\x00" +
	"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00" +
	"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00" +
	"\x00\x00\x00\x00" +
	"\x00\x00\x00\x00"

// ReadPackets reads packets from the AGWPE packet emitter and encodes them to be served by the web server
func (p PacketHandler) ReadPackets(jsonDtor *[]byte) {
	// unmarshal symbol descriptor
	symbolTable := SymbolTable{}
	err := json.Unmarshal(*jsonDtor, &symbolTable)
	if err != nil {
		panic("Couldn't unmarhsal JSON descriptor: " + err.Error())
	}
	if symbolTable.Symbols == nil || symbolTable.Symbols["/_"].Description != "Weather station" {
		panic("Error sanity-checking symbol table")
	}

	debug := false
	if os.Getenv("DEBUG") != "" {
		debug = true
	}
	for {
		serverAddr := fmt.Sprintf("%s:%v", p.Server, p.Port)
		log.Printf("[packethandler] Connecting to server...")
		conn, err := net.Dial("tcp", serverAddr)
		if err != nil {
			log.Printf("[packethandler] %v", err)
			time.Sleep(4 * time.Second)
			continue
		}
		log.Printf("[packethandler] Connected to %s\n", serverAddr)
		// Request raw TNC packets
		res, err := conn.Write([]byte(rawMode))
		if err != nil {
			log.Fatalf("[packethandler] Couldn't send 'Enable Raw Mode' bytes: %s", err)
		}
		if res != 36 {
			log.Fatalf("[packethandler] Failed connecting to TNC")
		}
		buf := make([]byte, 292)
		for {
			len, err := conn.Read(buf)
			if err != nil {
				log.Printf("[packethandler] error reading: %s\n", err)
				break
			}
			// fmt.Printf("Length of packet received is %v", len)
			// Now let's make a AGWPEPacket
			datLen := binary.LittleEndian.Uint32(buf[28:32])
			l := packet.AGWPEPacket{
				AGWPEPort: buf[1],
				DataKind:  buf[4],
				CallFrom:  string(buf[8:17]),
				CallTo:    string(buf[18:27]),
				DatLen:    datLen,
				Data:      buf[37 : len-1],
			}
			if debug {
				l.RawPacket = base64.StdEncoding.EncodeToString(buf[:len-1])
				fmt.Printf("AGWPE Packet: %+v\n", l)
				fmt.Printf("AGWPE Packet Data: %s", hex.Dump(l.Data))
			}
			checkin, err := packet.ParseAX25Frame(l.Data)
			if err != nil {
				log.Printf("[packethandler] Failed to parse this packet: %s %s\n", hex.Dump(buf), err)
				continue
			}
			if sym, ok := symbolTable.Symbols[checkin.SymbolTable+checkin.MapSymbol]; ok {
				checkin.SymbolDescription = sym.Description
			} else {
				checkin.SymbolDescription = fmt.Sprintf("Unknown, raw symbol %s", checkin.SymbolTable+checkin.MapSymbol)
			}
			// fmt.Printf("AGWPE Packet Data: %s\n", hex.Dump(buf[:len]))
			// fmt.Printf("AX25 Data: %s\n", hex.Dump(l.Data))
			// fmt.Printf("AX25 Frame: %+v\n", frame)
			//fmt.Printf("APRS Data: %s\n", hex.Dump(frame.APRSData))
			// // fmt.Printf("%x\n", buf[:len])
			// fmt.Println("################## END OF MESSAGE")
			// Now emit a frame
			// log.Printf("[packethandler] Sending checkin to Packet Store FROM %s\n", l.CallFrom)
			p.PacketStoreChan <- checkin
			// log.Printf("[packethandler] Sent packet to channel FROM %s\n", l.CallFrom)
		}
	}
}

type SymbolTable struct {
	Symbols map[string]Descriptor
}
type Descriptor struct {
	Tocall      string
	Description string
}
