package main

import (
	"errors"
	"log"
	"math/rand"
)

const defaultWidth = 20
const defaultHeight = 15
const defaultMinWallCount = 10
const defaultMaxWallCount = 30

type game struct {
	turn    bool
	move    chan *move
	blocks  [][]*block
	player1 *player
	player2 *player
	scoreP1 int
	scoreP2 int
}

func newGame(p1, p2 *player) (*game, []*coordinate) {
	blocks := make([][]*block, defaultHeight)
	for i := 0; i < defaultHeight; i++ {
		blocksRow := make([]*block, defaultWidth)
		for j := 0; j < defaultWidth; j++ {
			b := &block{
				x: j,
				y: i,
			}
			blocksRow[j] = b
		}
		blocks[i] = blocksRow
	}
	coordinates := randomCoordinate(defaultMinWallCount, defaultMaxWallCount, defaultWidth, defaultHeight)
	for _, coordinate := range coordinates {
		blocks[coordinate.Y][coordinate.X].s = -1
	}
	g := &game{
		move:    make(chan *move),
		blocks:  blocks,
		player1: p1,
		player2: p2,
	}
	return g, coordinates
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

func randomCoordinate(minWallCount, maxWallCount, width, height int) []*coordinate {
	count := int(rand.Int31n(int32(maxWallCount-minWallCount))) + minWallCount
	coordinates := make([]*coordinate, count)
	coordinateMap := make(map[*coordinate]interface{})
	for i := 0; i < count; {
		x := int(rand.Int31n(int32(width)))
		y := int(rand.Int31n(int32(height)))
		c := &coordinate{
			X: x,
			Y: y,
		}
		if _, ok := coordinateMap[c]; ok {
			continue
		}
		coordinates[i] = c
		coordinateMap[c] = nil
		i++
	}
	return coordinates
}
