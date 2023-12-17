package main

import (
	"log"
)

func main() {
	s := apiserver.New()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
