package main

import (
	"github.com/gregdel/pushover"
	"log"
)

type Push struct {
	app       *pushover.Pushover
	recipient *pushover.Recipient
}

func (push Push) Message(msg string) {
	message := pushover.NewMessage(msg)
	response, err := push.app.SendMessage(message, push.recipient)
	if err != nil {
		log.Panic(err)

		// Print the response if you want
		log.Println(response)
	}
}

func PushOver(pushoverToken string, pushoverRecipient string) Push {
	return Push{
		pushover.New(pushoverToken),
		pushover.NewRecipient(pushoverRecipient),
	}
}
