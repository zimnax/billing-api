package pub_sub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"log"
)

const (
	companyCreatedEvent = "CompanyCreatedEvent"
	companyUpdatedEvent = "CompanyUpdatedEvent"
)

type PubSub struct {
	subscription *pubsub.Subscription
}

func Init() *PubSub {
	ctx := context.Background()
	projectID := "internal"

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	s := client.Subscription("metadata-api-subscriber")

	return &PubSub{
		subscription: s,
	}
}

func (p *PubSub) Listen() {
	p.CompanyCreated()
	p.CompanyUpdated()
}

func (p *PubSub) CompanyCreated() {
	errorReceive := p.subscription.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		log.Printf("Got message in app handler: %s", m.Data)

		me := &MetadataEvent{}
		json.Unmarshal(m.Data, me)

		if me.EventType == companyCreatedEvent {
			//TODO: register company wallet
		}

		m.Ack()
	})
	if errorReceive != nil {
		fmt.Println("errorReceive", errorReceive)
	}
}

func (p *PubSub) CompanyUpdated() {
	errorReceive := p.subscription.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		log.Printf("Got message in app handler: %s", m.Data)

		me := &MetadataEvent{}
		json.Unmarshal(m.Data, me)

		if me.EventType == companyUpdatedEvent {
			//TODO
		}

		m.Ack()
	})
	if errorReceive != nil {
		fmt.Println("errorReceive", errorReceive)
	}
}
