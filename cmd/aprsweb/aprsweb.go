package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	_ "net/http/pprof"

	pb "github.com/ccmolik/aprsweb-proto/rpc/aprsweb"
	"github.com/ccmolik/aprsweb/bindata"
	"github.com/ccmolik/aprsweb/internal/aprswebtwirp"
	"github.com/ccmolik/aprsweb/internal/packethandler"
	"github.com/ccmolik/aprsweb/internal/packetstore"
)

var null = []byte("\x00")[0]
var webListenAddr string
var agwpeServer string
var agwpePort int
var Version = "development"

// RequestLogger logs every request and sets the version
func RequestLogger(targetMux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.Header().Set("X-aprsweb-version", Version)
		targetMux.ServeHTTP(w, r)

		log.Printf(
			"%s - - [%v] \"%s %s %s\" %v",
			r.RemoteAddr,
			start.Format(time.RFC1123),
			r.Method,
			r.RequestURI,
			r.Proto,
			time.Since(start),
		)
	})
}

// serveIndex handles serving the index.html file out of the assets bindata
func serveIndex() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		index, err := bindata.Asset("assets/index.html")
		if err != nil {
			panic("Couldn't get index.html")
		}
		w.Write(index)
	})
}

func main() {
	flag.StringVar(&webListenAddr, "l", ":5000", "webserver listen address")
	flag.StringVar(&agwpeServer, "s", "kd2aoy-pi", "AGWPE/Direwolf server address")
	flag.IntVar(&agwpePort, "p", 8000, "AGWPE/Direwolf server port")
	flag.Parse()

	framez := make(chan *pb.Checkin)
	packetHandler := &packethandler.PacketHandler{
		PacketStoreChan: framez,
		Server:          agwpeServer,
		Port:            int32(agwpePort),
	}
	// Start Packet handler goroutine
	go func() {
		log.Println("Starting packet handler goroutine")
		packetHandler.ReadPackets()
	}()

	// Start Packet Store goroutine
	packetStore := packetstore.PacketStore{
		InChan: framez,
	}
	go func() {
		log.Println("Starting Packet Store Goroutine")
		packetStore.Loop()
	}()

	log.Printf("Starting http server at %s", webListenAddr)
	twirpServer := &aprswebtwirp.Server{
		PacketStore: &packetStore,
		Version:     Version,
	}
	twirpHandler := pb.NewAPRSWebServer(twirpServer, nil)
	mux := http.NewServeMux()

	log.Printf("TWIRP handler available at %s", twirpHandler.PathPrefix())
	mux.Handle(twirpHandler.PathPrefix(), twirpHandler)
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(bindata.AssetFile())))
	// // mux.Handle("/list", h.DumpList(&receivedFrames))
	mux.Handle("/", serveIndex())
	err := http.ListenAndServe(webListenAddr, RequestLogger(mux))
	log.Fatal(err)
}
