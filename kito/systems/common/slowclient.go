package common

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/lib/network"
)

type ArtificallySlowClient struct {
	client  *network.Client
	latency time.Duration
}

func NewArtificallySlowClient(client *network.Client, latency time.Duration) *ArtificallySlowClient {
	return &ArtificallySlowClient{client: client, latency: latency}
}

func (c *ArtificallySlowClient) SendMessage(messageType int, messageBody any) error {
	go func() {
		time.Sleep(c.latency)
		err := c.client.SendMessage(messageType, messageBody)
		if err != nil {
			fmt.Println("artificiallySlowClient send message failed with error", err)
		}
	}()
	return nil
}

func (c *ArtificallySlowClient) PullIncomingMessages() []*network.Message {
	return c.client.PullIncomingMessages()
}
