package main

import (
	"github.com/mirmeep/gator-blog/internal/config"
	"errors"
	"fmt"
	"log"
)

func main() {
	conf_contents, err := config.Read()
	if err != nil {
		errors.New("cannot read config file")
	}

	conf_state := state{
		conf: conf_contents
	}
	fmt.Printf("Url: %v\n", conf.Url)
	fmt.Printf("Username: %v\n", conf.Username)
}