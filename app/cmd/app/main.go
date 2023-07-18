package main

import (
	"fmt"
	"main/pkg/clients/telegram"
	"os"
)

func main() {
	//err := gatecontroller.ControlGate(true)
	//
	//if err != nil {
	//	fmt.Println("Error to do action with gate:", err)
	//}

	tgClient := telegram.New("api.telegram.org", os.Getenv("BOT_TOKEN"))
	fmt.Println(tgClient)
}
