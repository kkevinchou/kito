package app

import "github.com/kkevinchou/kito/lib/network"

type StateBundle struct {
	cf     int
	update *network.GameStateUpdateMessage
}

type StateInterpolator struct {
	stateBuffer map[int]*StateBundle
	cfBuffer    []int

	bufferSize  int
	bufferIndex int
}

func NewStateInterpolator(bufferSize int) *StateInterpolator {
	return &StateInterpolator{
		stateBuffer: map[int]*StateBundle{},
		bufferSize:  bufferSize,
	}
}

// AppendCF assumes updates are appended in order
func (s *StateInterpolator) AppendCF(cf int, update *network.GameStateUpdateMessage) {
	s.stateBuffer[cf] = &StateBundle{
		cf:     cf,
		update: update,
	}
	s.cfBuffer = append(s.cfBuffer, cf)
}

func (s *StateInterpolator) Interpolate(cf int) *network.GameStateUpdateMessage {
	if len(s.cfBuffer) < s.bufferSize {
		// not enough updates in the buffer, return early
		return nil
	}

	if cf < s.cfBuffer[0] {
		// command frame is in the past
		return nil
	}

	if cf == s.cfBuffer[0] {
		return s.stateBuffer[cf].update
	}

	// TODO: handle cf that is beyond position 1
	if cf == s.cfBuffer[1] {
		delete(s.stateBuffer, s.cfBuffer[0])
		s.cfBuffer = s.cfBuffer[1:]
		return s.stateBuffer[cf].update
	}

	stateBundle := interpolate(cf, s.stateBuffer[s.cfBuffer[0]], s.stateBuffer[s.cfBuffer[1]])

	return stateBundle.update
}

func interpolate(cf int, start *StateBundle, end *StateBundle) *StateBundle {
	var interpolationRange float64 = float64(cf-start.cf) / float64(end.cf-start.cf)

	stateBundle := &StateBundle{
		cf: cf,
	}

	interpolatedEntities := map[int]network.EntitySnapshot{}
	for id, startSnapshot := range start.update.Entities {
		endSnapshot := end.update.Entities[id]

		interpolatedEntities[id] = network.EntitySnapshot{
			ID:       startSnapshot.ID,
			Type:     startSnapshot.Type,
			Position: endSnapshot.Position.Sub(startSnapshot.Position).Mul(interpolationRange).Add(startSnapshot.Position),
			// TODO: Orientation
		}
	}

	stateBundle.update.Entities = interpolatedEntities
	return stateBundle
}
