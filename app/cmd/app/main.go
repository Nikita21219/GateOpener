package main

import (
	"log"
	"main/pkg/clients/telegram"
	"main/pkg/consumer/event-consumer"
	telegram2 "main/pkg/events/telegram"
	"os"
)

func main() {
	//err := gate-controller.ControlGate(true)
	//
	//if err != nil {
	//	fmt.Println("Error to do action with gate:", err)
	//}

	tgClient := telegram.New("api.telegram.org", os.Getenv("BOT_TOKEN"))
	eventsProcessor := telegram2.New(tgClient)
	consumer := event_consumer.New(eventsProcessor, eventsProcessor, 100)
	if err := consumer.Start(); err != nil {
		log.Fatalln("error to start consumer:", err)
	}
}
