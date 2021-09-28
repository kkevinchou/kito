package network_test

import (
	"testing"

	"github.com/kkevinchou/kito/lib/network"
)

func TestBasic(t *testing.T) {
	host := "localhost"
	port := "8080"
	connectionType := "tcp"

	server := network.NewServer(host, port, connectionType, 19)
	err := server.Start()
	if err != nil {
		t.Errorf("failed to start server %s", err)
	}

	client, id, err := network.Connect(host, port, connectionType, 0)
	if err != nil {
		t.Errorf("failed to connect %s", err)
	}

	if id == network.UnsetClientID {
		t.Error("expected a non zero player ID from the accept message")
	}

	err = client.SendMessage(network.MessageTypeInput, nil)
	if err != nil {
		t.Error(err)
	}
}
