package messenger

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type DeliveryAgent interface {
	DeliverMessage(context.Context, *Message) error
}

var delivery_agent_roster roster.Roster

// DeliveryAgentInitializationFunc is a function defined by individual delivery_agent package and used to create
// an instance of that delivery_agent
type DeliveryAgentInitializationFunc func(ctx context.Context, uri string) (DeliveryAgent, error)

// RegisterDeliveryAgent registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `DeliveryAgent` instances by the `NewDeliveryAgent` method.
func RegisterDeliveryAgent(ctx context.Context, scheme string, init_func DeliveryAgentInitializationFunc) error {

	err := ensureDeliveryAgentRoster()

	if err != nil {
		return err
	}

	return delivery_agent_roster.Register(ctx, scheme, init_func)
}

func ensureDeliveryAgentRoster() error {

	if delivery_agent_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		delivery_agent_roster = r
	}

	return nil
}

// NewDeliveryAgent returns a new `DeliveryAgent` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `DeliveryAgentInitializationFunc`
// function used to instantiate the new `DeliveryAgent`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterDeliveryAgent` method.
func NewDeliveryAgent(ctx context.Context, uri string) (DeliveryAgent, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	fn, err := delivery_agent_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	if fn == nil {
		return nil, fmt.Errorf("Initialization function not defined")
	}

	init_func := fn.(DeliveryAgentInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func AgentSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureDeliveryAgentRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range delivery_agent_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
