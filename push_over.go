package main

import (
	"log"

	"github.com/gregdel/pushover"
)

type Push struct {
	app       *pushover.Pushover
	recipient *pushover.Recipient
}

func (push Push) Message(msg string) {
	if push.app != nil && push.recipient != nil {
		message := pushover.NewMessage(msg)
		_, err := push.app.SendMessage(message, push.recipient)
		if err != nil {
			log.Panic(err)
		}
	}
}

func PushOver(pushoverToken string, pushoverRecipient string) Push {
	push := Push{
		nil,
		nil,
	}
	if pushoverToken != "" && pushoverRecipient != "" {
		push = Push{
			pushover.New(pushoverToken),
			pushover.NewRecipient(pushoverRecipient),
		}
	}

	return push
}
