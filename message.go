package main

import (
	"errors"
)

const (
	actionJoin    = "join"
	actionLeave   = "leave"
	actionJoined  = "joined"
	actionLeft    = "left"
	actionMatched = "matched"
	actionMove    = "move"
	actionMoved   = "moved"
	actionEnd     = "end"
	actionEnded   = "ended"
	dataMoveX     = "x"
	dataMoveY     = "y"
	dataPlayer    = "p"
	dataScore     = "s"
	dataWalls     = "w"
	dataReason    = "r"
	dataWinner    = "winner"
	dataScoreP1   = "score-p1"
	dataScoreP2   = "score-p2"
)

type message struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

type coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

//todo send coordinate as json instead of two different data
func moveData(m *message) (*coordinate, error) {
	if m.Data != nil {
		dataX, okX := m.Data[dataMoveX]
		dataY, okY := m.Data[dataMoveY]
		if okX && okY {
			x, okX := dataX.(float64)
			y, okY := dataY.(float64)
			if okX && okY {
				c := &coordinate{
					X: int(x),
					Y: int(y),
				}
				return c, nil
			}
			return nil, errors.New("move data types in message are incompatible")
		}
		return nil, errors.New("move data in message are missing")
	}
	return nil, errors.New("data in message is nil")
}

func reasonData(m *message) (string, error) {
	if m.Data != nil {
		dataReason, ok := m.Data[dataReason]
		if ok {
			reason, ok := dataReason.(string)
			if ok {
				return reason, nil
			}
			return "", errors.New("reason data types in message are incompatible")
		}
		return "", errors.New("reason data in message are missing")
	}
	return "", errors.New("data in message is nil")
}
