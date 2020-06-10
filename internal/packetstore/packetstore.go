package packetstore

import (
	pb "github.com/ccmolik/aprsweb-proto/rpc/aprsweb"
)

type PacketStore struct {
	InChan         chan *pb.Checkin
	receivedFrames []*pb.Checkin
}

func (s *PacketStore) Loop() {
	// Run this in a goroutine, srsly
	for {
		// Read from channel
		// log.Printf("[packetstore] Waiting for checkin...")

		f := <-s.InChan
		// log.Printf("[packetstore] Received from channel FROM %s\n", f.Source)
		s.receivedFrames = append(s.receivedFrames, f)
		// log.Printf("[packetstore] Appended to receivedFrames.")

	}

}

// GetCheckins gets a list of _all_ checkins
func (s PacketStore) GetAllCheckins() []*pb.Checkin {
	return s.receivedFrames
}
