package main

import (
	"errors"
	"log"
)

const defaultWidth = 20
const defaultHeight = 15

type game struct {
	turn    bool
	move    chan *move
	blocks  [][]*block
	player1 *player
	player2 *player
	scoreP1 int
	scoreP2 int
}

func newGame(p1, p2 *player) *game {
	blocksAll := [defaultHeight][]*block{}
	for i := 0; i < defaultHeight; i++ {
		blocksRow := [defaultWidth]*block{}
		for j := 0; j < defaultWidth; j++ {
			blocksRow[j] = &block{
				x: j,
				y: i,
			}
		}
		blocksAll[i] = blocksRow[:]
	}
	blocks := blocksAll[:]
	return &game{
		move:    make(chan *move),
		blocks:  blocks,
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
			b := g.blocks[mv.y][mv.x]
			if b.s != 0 {
				err := errors.New("invalid move")
				log.Println(err)
				continue
			}
			if !mv.p {
				b.s = 1
				g.scoreP1++
			} else {
				b.s = 2
				g.scoreP2++
			}
			data := make(map[string]interface{})
			data[dataMoveX] = mv.x
			data[dataMoveY] = mv.y
			data[dataPlayer] = mv.p
			if !mv.p {
				data[dataScore] = g.scoreP1
			} else {
				data[dataScore] = g.scoreP2
			}
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
