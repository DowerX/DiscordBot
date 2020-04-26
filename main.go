package main

import (
	"fmt"
	"./config"
	"./bot"
)

func main() {
	fmt.Println("Starting...")

	bot.Start(config.GetConfig())
}
