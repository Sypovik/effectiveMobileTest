package main

import (
	"fmt"

	"github.com/Sypovik/effectiveMobileTest/internal/config"
)

func main() {
	config := config.LoadConfig()
	fmt.Println(config)
}
