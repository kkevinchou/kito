package commandframe

// // AppendCF assumes updates are appended in order
// func (s *CommandFrameHistory) AppendCF(cf int, update *network.GameStateUpdateMessage) {
// 	s.stateBuffer[cf] = &CommandFrame{
// 		cf:     cf,
// 		update: update,
// 	}
// 	s.cfBuffer = append(s.cfBuffer, cf)
// }

// func (s *CommandFrameHistory) Interpolate(cf int) *network.GameStateUpdateMessage {
// 	if len(s.cfBuffer) < s.bufferSize {
// 		if len(s.cfBuffer) > 0 {
// 			// if we're slow in receiving updates we shouldn't return nil
// 			// TODO: have a less aggressive solution than halting the world
// 			return s.stateBuffer[s.cfBuffer[0]].update
// 		}
// 		return nil
// 	}

// 	if cf < s.cfBuffer[0] {
// 		// command frame is in the past
// 		return nil
// 	}

// 	if cf == s.cfBuffer[0] {
// 		return s.stateBuffer[cf].update
// 	}

// 	// TODO: handle cf that is beyond position 1
// 	if cf == s.cfBuffer[1] {
// 		delete(s.stateBuffer, s.cfBuffer[0])
// 		s.cfBuffer = s.cfBuffer[1:]
// 		return s.stateBuffer[cf].update
// 	}

// 	stateBundle := interpolate(cf, s.stateBuffer[s.cfBuffer[0]], s.stateBuffer[s.cfBuffer[1]])

// 	return stateBundle.update
// }

// func interpolate(cf int, start *CommandFrame, end *CommandFrame) *CommandFrame {
// 	var interpolationRange float64 = float64(cf-start.cf) / float64(end.cf-start.cf)

// 	stateBundle := &CommandFrame{
// 		cf: cf,
// 	}

// 	interpolatedEntities := map[int]network.EntitySnapshot{}
// 	for id, startSnapshot := range start.update.Entities {
// 		endSnapshot := end.update.Entities[id]

// 		interpolatedEntities[id] = network.EntitySnapshot{
// 			ID:       startSnapshot.ID,
// 			Type:     startSnapshot.Type,
// 			Position: endSnapshot.Position.Sub(startSnapshot.Position).Mul(interpolationRange).Add(startSnapshot.Position),
// 			// TODO: Orientation
// 		}
// 	}

// 	stateBundle.update.Entities = interpolatedEntities
// 	return stateBundle
// }
