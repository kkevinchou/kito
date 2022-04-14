package inputbuffer

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/kito/knetwork"
	"github.com/kkevinchou/kito/lib/input"
)

// InputBuffer is a buffer of inputs. Inputs are sent from clients and stored in the buffer
// until the server is ready to consume them. The internet is a wild and scary place - inputs
// from clients can arrive in bursts or with huge delays between each input. What the server
// would like to see is a nice steady stream of inputs from clients. To accomplish this, we
// add in a configurable artificial latency that we buffer inputs behind. The job of InputBuffer
// is to abstract this and provide that steady stream of client inputs to the server.

// TODO: when we begin sending messages via UDP which can result in packet loss, we will want
// to buffer a set of inputs (up until the last acknowledge command frame) so that when a
// packet finally does arrive on the server, the server can fill the dropped inputs in the
// input buffer

type BufferedInput struct {
	TargetGlobalCommandFrame int
	LocalCommandFrame        int
	PlayerID                 int
	Input                    input.Input
	ReceivedTimestamp        time.Time
}

type InputBuffer struct {
	playerInputs     map[int][]BufferedInput
	maxCommandFrames int
	seenInputs       map[int]map[int]any
}

func NewInputBuffer(maxCommandFrames int) *InputBuffer {
	return &InputBuffer{
		playerInputs:     map[int][]BufferedInput{},
		maxCommandFrames: maxCommandFrames,
	}
}

func (i *InputBuffer) PushInput(globalCommandFrame int, localCommandFrame int, lastInputCommandFrame int, playerID int, receivedTime time.Time, networkInput *knetwork.InputMessage) {
	var targetGlobalCommandFrame int
	if len(i.playerInputs[playerID]) > 0 {
		lastPlayerInput := i.playerInputs[playerID][len(i.playerInputs[playerID])-1]
		commandFrameDelta := localCommandFrame - lastPlayerInput.LocalCommandFrame

		// assuming a properly behaving client they should only send one input message per
		// command frame. in the event that they send more than one, we naively set it for
		// the next command frame
		if commandFrameDelta <= 0 {
			commandFrameDelta = 1
			fmt.Println("warning: received more than one input for a given command frame")
		}

		targetGlobalCommandFrame = lastPlayerInput.TargetGlobalCommandFrame + commandFrameDelta
	} else {
		targetGlobalCommandFrame = globalCommandFrame + i.maxCommandFrames
	}

	// target exceeds the buffer size. drop the input to avoid spiral of death
	if targetGlobalCommandFrame > (globalCommandFrame + i.maxCommandFrames) {
		fmt.Printf("target gcf exceeded buffer size %d > %d\n", targetGlobalCommandFrame, (globalCommandFrame + i.maxCommandFrames))
		return
	}

	i.playerInputs[playerID] = append(
		i.playerInputs[playerID],
		BufferedInput{
			PlayerID:                 playerID,
			LocalCommandFrame:        localCommandFrame,
			TargetGlobalCommandFrame: targetGlobalCommandFrame,
			Input:                    networkInput.Input,
			ReceivedTimestamp:        receivedTime,
		},
	)
}

// PullInput pulls a buffered input for the current command frame
func (i *InputBuffer) PullInput(globalCommandFrame int, playerID int) *BufferedInput {
	if len(i.playerInputs[playerID]) == 0 {
		return nil
	}

	playerInput := i.playerInputs[playerID][0]
	if playerInput.TargetGlobalCommandFrame <= globalCommandFrame {
		i.playerInputs[playerID] = i.playerInputs[playerID][1:]
		return &playerInput
	}

	return nil
}
