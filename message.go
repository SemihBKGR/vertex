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
	dataMoveX     = "x"
	dataMoveY     = "y"
	dataPlayer    = "p"
	dataScore     = "s"
)

type message struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

func moveData(m *message) (*move, error) {
	if m.Data != nil {
		dataX, okX := m.Data[dataMoveX]
		dataY, okY := m.Data[dataMoveY]
		if okX && okY {
			x, okX := dataX.(float64)
			y, okY := dataY.(float64)
			if okX && okY {
				m := &move{
					x: int(x),
					y: int(y),
				}
				return m, nil
			}
			return nil, errors.New("move data types in message are incompatible")
		}
		return nil, errors.New("move data in message are missing")
	}
	return nil, errors.New("data in message is nil")
}
