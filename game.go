package main

import (
	"errors"
	"log"
)

type game struct {
	turn    bool
	move    chan *move
	blocks  [][]*block
	player1 *player
	player2 *player
}

func newGame(p1, p2 *player) *game {
	return &game{
		move:    make(chan *move),
		player1: p1,
		player2: p2,
	}
}

func (g *game) startGame() {
	for {
		select {
		case mv := <-g.move:
			if mv.p != g.turn {
				err := errors.New("move is not belong to player who has turn")
				log.Println(err)
				continue
			}
			data := make(map[string]interface{})
			data[dataMoveX] = mv.x
			data[dataMoveY] = mv.y
			data[dataMovePlayer] = mv.p
			m := &message{
				Action: actionMoved,
				Data:   data,
			}
			g.player1.sendMessage(m)
			g.player2.sendMessage(m)
			g.turn = !g.turn
		}
	}
}

type block struct {
	s int
	x int
	y int
}

type move struct {
	p bool
	x int
	y int
}
