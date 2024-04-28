package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func main() {
	s1, err := net.Dial("tcp", "192.168.0.18:8088")
	if err != nil {
		log.Println("Error connecting:", err)
		return
	}
	defer s1.Close()

	for {
		data := make([]byte, 1024)
		_, err := s1.Read(data)
		if err != nil {
			log.Println("Error reading:", err)
			return
		}

		s2, err := net.Dial("tcp", "192.168.0.18:8088")
		if err != nil {
			log.Println("Error connecting:", err)
			return
		}
		defer s2.Close()

		response, err := http.Get("http://localhost:8080")
		if err != nil {
			log.Println("Error getting response:", err)
			return
		}
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println("Error reading response:", err)
			return
		}

		_, err = s2.Write(body)
		if err != nil {
			log.Println("Error sending response:", err)
			return
		}
	}
}
