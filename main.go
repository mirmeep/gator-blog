package main

import (
	"github.com/mirmeep/gator-blog/internal/config"
	"errors"
	"fmt"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		errors.New("cannot read config file")
	}

	conf.SetUser("meep")
	conf, conf_err := config.Read()
	if conf_err != nil {
		errors.New("cannot change username")
	}
	fmt.Println(conf)
}