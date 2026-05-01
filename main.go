package main

import (
	"fmt"

	"github.com/cyberfly100/bootdev_gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}
	cfg.SetUser("lucas")
	new_cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}
	fmt.Println(new_cfg)
}
