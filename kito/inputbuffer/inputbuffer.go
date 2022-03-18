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
		seenInputs:       map[int]map[int]any{},
	}
}

func (i *InputBuffer) StartFrame(gcf int) {
	// clear out the old seen inputs
	if _, ok := i.seenInputs[gcf-1]; ok {
		delete(i.seenInputs, gcf-1)
	}
	i.seenInputs[gcf] = map[int]any{}
}

func (i *InputBuffer) PushInput(globalCommandFrame int, localCommandFrame int, lastInputCommandFrame int, playerID int, receivedTime time.Time, networkInput *knetwork.InputMessage) {
	// this handles the case where a client can spam the server with inputs (e.g. if they hold down the title bar and stack up a bunch of network requests).
	// we avoid processing all the inputs so that we don't lag for the inputs of that user. This is because we only process one input for each command frame.
	// the result is that we only process the first input and drop the rest.
	if _, ok := i.seenInputs[globalCommandFrame][playerID]; ok {
		return
	}
	i.seenInputs[globalCommandFrame][playerID] = true

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
		if _, ok := networkInput.Input.KeyboardInput[input.KeyboardKeySpace]; ok {
		}
		targetGlobalCommandFrame = globalCommandFrame + i.maxCommandFrames
	}

	i.playerInputs[playerID] = append(
		i.playerInputs[playerID],
		BufferedInput{
			LocalCommandFrame:        localCommandFrame,
			PlayerID:                 playerID,
			Input:                    networkInput.Input,
			ReceivedTimestamp:        receivedTime,
			TargetGlobalCommandFrame: targetGlobalCommandFrame,
		},
	)
}

// PullInput pulls a buffered input for the current command frame
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
