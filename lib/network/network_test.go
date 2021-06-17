package network_test

import (
	"testing"

	"github.com/kkevinchou/kito/lib/network"
)

func TestBasic(t *testing.T) {
	host := "localhost"
	port := "8080"
	connectionType := "tcp"

	server := network.NewServer(host, port, connectionType)
	err := server.Start()
	if err != nil {
		t.Errorf("failed to start server %s", err)
	}

	client := network.NewClient()
	acceptMessage, err := client.Connect(host, port, connectionType)
	if err != nil {
		t.Errorf("failed to connect %s", err)
	}

	if acceptMessage.PlayerID == 0 {
		t.Error("expected a non zero player ID from the accept message")
	}
}
