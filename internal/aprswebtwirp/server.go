package aprswebtwirp

import (
	"context"
	"time"

	"github.com/ccmolik/aprsweb-proto/rpc/aprsweb"
	pb "github.com/ccmolik/aprsweb-proto/rpc/aprsweb"
	"github.com/ccmolik/aprsweb/internal/packetstore"
)

// Server implements the APRSWeb service
type Server struct {
	PacketStore *packetstore.PacketStore
	Version     string
}

// GetCheckins gets checkins.
func (s *Server) GetCheckins(ctx context.Context, req *pb.GetCheckinsRequest) (*pb.GetCheckinsResponse, error) {
	checkins := s.PacketStore.GetAllCheckins()
	// Filter based on requested datetime (if we got one)
	// Couldn't think of a clever way in Go of doing so, so we allocate a new array and copy over if they match >_>
	if req.Since != 0 {
		var retCheckins = make([]*pb.Checkin, 0)
		sinceTime := time.Unix(0, req.Since)
		for _, c := range checkins {
			if time.Unix(0, c.ReceivedDatetime).After(sinceTime) {
				retCheckins = append(retCheckins, c)
			}
		}
		checkins = retCheckins
	}
	resp := &aprsweb.GetCheckinsResponse{
		Checkins: checkins,
	}

	return resp, nil
}
