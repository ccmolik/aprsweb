package web

import (
	"time"
)

// RootHandler is the actual root handler
type RootHandler struct {
	StartTime time.Time
	// 	FrameChannel   chan *packet.AX25Frame
	// 	receivedFrames []*packet.AX25Frame
}

// // func (h RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// // 	l := log.Logger{}
// // 	log.Println("ServeHTTP got a request, passing to logging")
// // 	rootHandler := rootHandler{}
// // }

// // Logging handles request logging and then passes it along to the rootHandler

//  MapHandler handles plain HTTP requests to /
// func (h RootHandler) MapHandler() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		index, err := aprsweb.Asset("index.html")
// 		if err != nil {
// 			panic("Couldn't get index.html")
// 		}
// 		w.Write(index)
// 	})
// }

// // LeafletCSSHandler handles the HTTP get of leaflet.css

// func (h RootHandler) LeafletCSSHandler() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-type", "text/css")
// 		w.Write([]byte(leafletCSS))
// 	})

// }

// // LeafletJSHandler handles the HTTP get of leaflet.js

// func (h RootHandler) LeafletJSHandler() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-type", "application/javascript")
// 		w.Write([]byte(leafletJS))
// 	})
// }

// // ProtoHandler handles the HTTP get of aprweb.proto
// func (h RootHandler) ProtoHandler() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte(protoSchema))
// 	})
// }

// // MarkerHandler handles getting the marker img
// func (h RootHandler) MarkerHandler() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		f, err := ioutil.ReadFile("images/marker-icon.png")
// 		if err != nil {
// 			fmt.Printf("Error reading image file: %v\n", err)
// 			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		} else {
// 			w.Header().Set("Content-Type", "image/png")
// 			w.Write(f)
// 		}
// 	})
// }

// // MarkerShadowHandler handles serving the shadow
// func (h RootHandler) ShadowHandler() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		f, err := ioutil.ReadFile("images/marker-shadow.png")
// 		if err != nil {
// 			fmt.Printf("Error reading image file: %v\n", err)
// 			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		} else {
// 			w.Header().Set("Content-Type", "image/png")
// 			w.Write(f)
// 		}
// 	})
// }

// // Marker2XHandler handles serving the 2XMarker
// func (h RootHandler) Marker2XHandler() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		f, err := ioutil.ReadFile("images/marker-icon-2x.png")
// 		if err != nil {
// 			fmt.Printf("Error reading image file: %v\n", err)
// 			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		} else {
// 			w.Header().Set("Content-Type", "image/png")
// 			w.Write(f)
// 		}
// 	})
// }

// // LayersHandler handles serving the Layers
// func (h RootHandler) LayersHandler() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		f, err := ioutil.ReadFile("images/layers.png")
// 		if err != nil {
// 			fmt.Printf("Error reading image file: %v\n", err)
// 			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		} else {
// 			w.Header().Set("Content-Type", "image/png")
// 			w.Write(f)
// 		}
// 	})
// }
