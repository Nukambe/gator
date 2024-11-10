package main

import (
	"fmt"
	cfg "github.com/Nukambe/gator/internal/config"
)

func main() {
	config, err := cfg.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = config.SetUser("wesley")
	if err != nil {
		fmt.Println(err)
		return
	}
	config, err = cfg.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(config)
}
