package main

import (
	"fmt"
	"github.com/gregdel/pushover"
	"log"
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
	} else {
		fmt.Println(msg)
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
