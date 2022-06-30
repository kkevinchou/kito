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
	playerInputs     map[int]map[int]BufferedInput
	lastPlayerInput  map[int]int
	maxCommandFrames int
	seenInputs       map[int]map[int]any
}

func NewInputBuffer(maxCommandFrames int) *InputBuffer {
	return &InputBuffer{
		playerInputs:     map[int]map[int]BufferedInput{},
		lastPlayerInput:  map[int]int{},
		maxCommandFrames: maxCommandFrames,
	}
}

func (inputBuffer *InputBuffer) PushInput(globalCommandFrame int, localCommandFrame int, playerID int, receivedTime time.Time, networkInput *knetwork.InputMessage) {
	if _, ok := inputBuffer.playerInputs[playerID]; !ok {
		inputBuffer.playerInputs[playerID] = map[int]BufferedInput{}
	}

	targetGlobalCommandFrame := globalCommandFrame + inputBuffer.maxCommandFrames
	if len(inputBuffer.playerInputs[playerID]) > 0 {
		lastPlayerInputCF := inputBuffer.lastPlayerInput[playerID]
		lastPlayerInput := inputBuffer.playerInputs[playerID][lastPlayerInputCF]
		commandFrameDelta := localCommandFrame - lastPlayerInput.LocalCommandFrame

		// assuming a properly behaving client they should only send one input message per
		// command frame. in the event that they send more than one, we naively set it for
		// the next command frame
		if commandFrameDelta <= 0 {
			commandFrameDelta = 1
			fmt.Println("warning: received more than one input for a given command frame")
		}

		relativeCF := lastPlayerInput.TargetGlobalCommandFrame + commandFrameDelta
		if relativeCF >= globalCommandFrame && relativeCF < targetGlobalCommandFrame {
			// fmt.Println(globalCommandFrame, "|", lastPlayerInput.TargetGlobalCommandFrame, relativeCF, targetGlobalCommandFrame)
			targetGlobalCommandFrame = relativeCF
		}
		// fmt.Println(globalCommandFrame, "|", targetGlobalCommandFrame)
	}

	// target exceeds the buffer size. drop the input to avoid spiral of death
	if targetGlobalCommandFrame > (globalCommandFrame + inputBuffer.maxCommandFrames) {
		fmt.Printf("target gcf exceeded buffer size %d > %d\n", targetGlobalCommandFrame, (globalCommandFrame + inputBuffer.maxCommandFrames))
		return
	}

	inputBuffer.playerInputs[playerID][targetGlobalCommandFrame] = BufferedInput{
		PlayerID:                 playerID,
		LocalCommandFrame:        localCommandFrame,
		TargetGlobalCommandFrame: targetGlobalCommandFrame,
		Input:                    networkInput.Input,
		ReceivedTimestamp:        receivedTime,
	}
	inputBuffer.lastPlayerInput[playerID] = targetGlobalCommandFrame
}

// PullInput pulls a buffered input for the current command frame
func (inputBuffer *InputBuffer) PullInput(globalCommandFrame int, playerID int) *BufferedInput {
	if len(inputBuffer.playerInputs[playerID]) == 0 {
		return nil
	}

	if in, ok := inputBuffer.playerInputs[playerID][globalCommandFrame]; ok {
		return &in
	}

	max := 5
	for i := 1; i < max; i++ {
		if in, ok := inputBuffer.playerInputs[playerID][globalCommandFrame-i]; ok {
			copiedInput := in
			copiedInput.LocalCommandFrame += i
			copiedInput.TargetGlobalCommandFrame = globalCommandFrame
			inputBuffer.playerInputs[playerID][globalCommandFrame] = copiedInput

			result := inputBuffer.playerInputs[playerID][globalCommandFrame]
			return &result
		}
	}

	fmt.Printf("failed to fetch input for player %d cf %d\n", playerID, globalCommandFrame)
	return nil
}
