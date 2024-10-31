package messenger

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/aaronland/gomail-sender"
	"github.com/aaronland/gomail/v2"
)

type EmailDeliveryAgent struct {
	DeliveryAgent
	sender gomail.Sender
}

// In principle this could also be done with a sync.OnceFunc call but that will
// require that everyone uses Go 1.21 (whose package import changes broke everything)
// which is literally days old as I write this. So maybe a few releases after 1.21

var register_mu = new(sync.RWMutex)
var register_map = map[string]bool{}

func init() {

	ctx := context.Background()
	err := RegisterEmailSchemes(ctx)

	if err != nil {
		panic(err)
	}
}

func RegisterEmailSchemes(ctx context.Context) error {

	for _, uri := range sender.Schemes() {

		s := strings.Replace(uri, "://", "", 1)
		s = fmt.Sprintf("email-%s", s)

		_, exists := register_map[s]

		if exists {
			continue
		}

		err := RegisterDeliveryAgent(ctx, s, NewEmailDeliveryAgent)

		if err != nil {
			return err
		}

		register_map[s] = true
	}

	return nil
}

func NewEmailDeliveryAgent(ctx context.Context, uri string) (DeliveryAgent, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	u.Scheme = strings.TrimLeft(u.Scheme, "email-")

	s, err := sender.NewSender(ctx, u.String())

	if err != nil {
		return nil, fmt.Errorf("Failed to create new email sender, %w", err)
	}

	agent := &EmailDeliveryAgent{
		sender: s,
	}

	return agent, nil
}

func (agent *EmailDeliveryAgent) DeliverMessage(ctx context.Context, msg *Message) error {

	gomail_msg := gomail.NewMessage()

	gomail_msg.SetHeader("Subject", msg.Subject)
	gomail_msg.SetHeader("From", msg.From)
	gomail_msg.SetHeader("To", msg.To)

	gomail_msg.SetBody("text/html", msg.Body)

	err := gomail.Send(agent.sender, gomail_msg)

	if err != nil {
		return fmt.Errorf("Failed to deliver message, %w", err)
	}

	return nil
}
