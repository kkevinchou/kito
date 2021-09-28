package inputbuffer

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/lib/input"
	"github.com/kkevinchou/kito/lib/network"
)

// InputBuffer is a buffer of inputs. Inputs are sent from clients and stored in the buffer
// until the server is ready to consume them. The internet is a wild and scary place - inputs
// from clients can arrive in bursts or with huge delays between each input. What the server
// would like to see is a nice steady stream of inputs from clients. To accomplish this, we
// add in a configurable artificial latency that we buffer inputs behind. The job of InputBuffer
// is to abstract this and provide that steady stream of client inputs to the server.

type BufferedInput struct {
	TargetGlobalCommandFrame int
	PlayerCommandFrame       int
	PlayerID                 int
	Input                    input.Input
	ReceivedTimestamp        time.Time
}

type InputBuffer struct {
	playerInputs map[int][]BufferedInput
	bufferSize   int
}

func NewInputBuffer(bufferSize int) *InputBuffer {
	return &InputBuffer{
		playerInputs: map[int][]BufferedInput{},
		bufferSize:   bufferSize,
	}
}

func (i *InputBuffer) PushInput(globalCommandFrame int, playerCommandFrame int, playerID int, receivedTime time.Time, networkInput *network.InputMessage) {
	var commandFrameDelta int
	if len(i.playerInputs[playerID]) > 0 {
		lastCommandFrame := i.playerInputs[playerID][len(i.playerInputs)-2]
		commandFrameDelta = playerCommandFrame - lastCommandFrame.PlayerCommandFrame

		// assuming a properly behaving client they should only send one input message per
		// command frame. in the event that they send more than one, we naievly set it for
		// the next command frame
		if commandFrameDelta <= 0 {
			commandFrameDelta = 1
			fmt.Println("warning: received more than one input for a given command frame")
		}
	} else {
		commandFrameDelta = i.bufferSize
	}

	i.playerInputs[playerID] = append(
		i.playerInputs[playerID],
		BufferedInput{
			PlayerCommandFrame:       playerCommandFrame,
			PlayerID:                 playerID,
			Input:                    networkInput.Input,
			ReceivedTimestamp:        receivedTime,
			TargetGlobalCommandFrame: globalCommandFrame + commandFrameDelta,
		},
	)
}

// PullInput pulls messages for the current command frame
func (i *InputBuffer) PullInput(globalCommandFrame int, playerID int) *BufferedInput {
	if len(i.playerInputs[playerID]) == 0 {
		return nil
	}

	input := i.playerInputs[playerID][0]
	if input.TargetGlobalCommandFrame <= globalCommandFrame {
		i.playerInputs[playerID] = i.playerInputs[playerID][1:]
		return &input
	}

	return nil
}
