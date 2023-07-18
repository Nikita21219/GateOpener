package event_consumer

import (
	"log"
	"main/pkg/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("error fetch events in consumer: %s", err.Error())
			continue
		}
		if len(gotEvents) == 0 {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		c.handleEvents(gotEvents)
	}
}

func (c *Consumer) handleEvents(events []events.Event) {
	for _, e := range events {
		if err := c.processor.Process(e); err != nil {
			log.Printf("can't handle event: %s", err.Error())
			continue
		}
	}
}
