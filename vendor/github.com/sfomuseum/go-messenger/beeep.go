package messenger

import (
	"context"

	"github.com/gen2brain/beeep"
)

type BeeepDeliveryAgent struct {
	DeliveryAgent
}

func init() {
	ctx := context.Background()
	RegisterDeliveryAgent(ctx, "beeep", NewBeeepDeliveryAgent)
}

func NewBeeepDeliveryAgent(ctx context.Context, uri string) (DeliveryAgent, error) {

	a := &BeeepDeliveryAgent{}
	return a, nil
}

func (a *BeeepDeliveryAgent) DeliverMessage(ctx context.Context, msg *Message) error {
	return beeep.Alert(msg.Subject, msg.Body, "assets/information.png")
}
