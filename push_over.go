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
	_, err := push.app.SendMessage(message, push.recipient)
	if err != nil {
		log.Panic(err)
	}
}

func PushOver(pushoverToken string, pushoverRecipient string) Push {
	return Push{
		pushover.New(pushoverToken),
		pushover.NewRecipient(pushoverRecipient),
	}
}
