package pushover

import (
	"fmt"
	"os"

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
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func (push Push) Validate() error {
	if push.app != nil && push.recipient != nil {
		_, err := push.app.GetRecipientDetails(push.recipient)
		return err
	}

	return nil
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
