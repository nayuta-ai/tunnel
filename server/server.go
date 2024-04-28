package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	controlServer, err := net.Listen("tcp", "192.168.0.18:8088")
	if err != nil {
		log.Println("Error connecting:", err)
		return
	}
	defer controlServer.Close()

	con1, err := controlServer.Accept()
	if err != nil {
		log.Println("Error accepting connection:", err)
		return
	}
	defer con1.Close()

	for {
		// Create handler for HTTP server
		handler := func(w http.ResponseWriter, r *http.Request) {
			con1.Write([]byte(fmt.Sprintf("%s %s %s", r.Method, r.URL.Path, r.Proto)))
			con2, err := controlServer.Accept()
			if err != nil {
				log.Println("Error accepting connection:", err)
				return
			}
			defer con2.Close()

			// Read data from second connection and write to response
			data := make([]byte, 2048)
			_, err = con2.Read(data)
			if err != nil {
				log.Println("Error reading data:", err)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(data)
		}

		// Create HTTP server
		httpServer := http.Server{
			Addr:    "192.168.0.18:8089",
			Handler: http.HandlerFunc(handler),
		}

		// Start serving requests
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Println("Error starting HTTP server:", err)
			return
		}
	}
}
