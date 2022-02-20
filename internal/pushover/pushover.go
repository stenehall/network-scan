package pushover

import (
	"fmt"
	"os"

	"github.com/gregdel/pushover"
)

// Push model
type Push struct {
	app       *pushover.Pushover
	recipient *pushover.Recipient
}

// Message sends a new message over push over
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

// Validate the provided setup
func (push Push) validate() error {
	if push.app != nil && push.recipient != nil {
		_, err := push.app.GetRecipientDetails(push.recipient)
		return err
	}

	return nil
}

// PushOver creates a new instance of the pushover client
func PushOver(pushoverToken string, pushoverRecipient string) (Push, error) {
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

	err := push.validate()

	return push, err
}
