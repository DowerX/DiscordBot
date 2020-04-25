package main

import (
	"fmt"

	"./bot"
)

func main() {
	fmt.Println("Parsing token")
	bot.Start(getToken())
}
