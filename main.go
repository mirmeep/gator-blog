package main

import (
	"github.com/mirmeep/gator-blog/internal/config"
	"errors"
	"fmt"
	"log"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		errors.New("cannot read config file")
	}

	err = conf.SetUser("meepowo")
	if err != nil {
		log.Fatalf("couldn't set current user: %v", err)
	}
	conf, conf_err := config.Read()
	if conf_err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Url: %v\n", conf.Url)
	fmt.Printf("Username: %v\n", conf.Username)
}