package main

import (
	"fmt"
	"log"

	"github.com/its-me-sv/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Config: %+v\n", cfg)

	if err = config.SetUser("suraj"); err != nil {
		log.Fatalln(err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Config: %+v\n", cfg)
}
