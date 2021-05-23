package main

import (
	"time"

	"github.com/Shikugawa/pingoo/pkg/consumer"
	"github.com/Shikugawa/pingoo/pkg/provider"
)

func main() {
	c := provider.NewTrelloProvider()
	s := consumer.NewSlackClient()

	notificationTimer := time.NewTicker(1 * time.Hour)

	go func() {
		for {
			select {
			case <-notificationTimer.C:
				tasks, err := c.Get()
				if err != nil {
					continue
				}

				if err := s.Post(tasks); err != nil {
					continue
				}
			}
		}
	}()

	h := consumer.NewSlackHandler(c)
	h.Start(3000)
}
