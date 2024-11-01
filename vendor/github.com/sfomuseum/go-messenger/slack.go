package messenger

import (
	"context"
	"fmt"
	"net/url"

	"github.com/sfomuseum/go-slack"
)

type SlackDeliveryAgent struct {
	DeliveryAgent
	webhook *slack.Webhook
	channel string
	from    string
}

func init() {
	ctx := context.Background()
	RegisterDeliveryAgent(ctx, "slack", NewSlackDeliveryAgent)
}

func NewSlackDeliveryAgent(ctx context.Context, uri string) (DeliveryAgent, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	webhook_uri := q.Get("webhook")

	wh, err := slack.NewWebhook(ctx, webhook_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create webhook, %w", err)
	}

	channel := q.Get("channel")
	from := q.Get("from")

	a := &SlackDeliveryAgent{
		webhook: wh,
		channel: channel,
		from:    from,
	}

	return a, nil
}

func (a *SlackDeliveryAgent) DeliverMessage(ctx context.Context, msg *Message) error {

	channel := a.channel
	from := a.from

	if msg.To != "" {
		channel = msg.To
	}

	if msg.From != "" {
		from = msg.From
	}

	m := &slack.Message{
		Channel: channel,
		Text:    fmt.Sprintf("[%s] %s %s", from, msg.Subject, msg.Body),
	}

	a.webhook.Post(ctx, m)

	return nil
}
