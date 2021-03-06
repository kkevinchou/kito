package networkdispatch

import (
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/lib/network"
)

type MessageFetcher func(world World) []*network.Message

func defaultMessageFetcher(world World) []*network.Message {
	return nil
}

func ConnectedPlayersMessageFetcher(world World) []*network.Message {
	playerManager := directory.GetDirectory().PlayerManager()
	var allMessages []*network.Message

	for _, player := range playerManager.GetPlayers() {
		messages := player.Client.PullIncomingMessages()
		allMessages = append(allMessages, messages...)
	}

	return allMessages
}

func ClientMessageFetcher(world World) []*network.Message {
	singleton := world.GetSingleton()
	playerManager := directory.GetDirectory().PlayerManager()

	player := playerManager.GetPlayer(singleton.PlayerID)
	return player.Client.PullIncomingMessages()
}
