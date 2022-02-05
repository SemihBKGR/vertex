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
	width   int
	height  int
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
		width:   defaultWidth,
		height:  defaultHeight,
		player1: p1,
		player2: p2,
	}
	return g, coordinates
}

func (g *game) startGame() {
	for {
		select {
		case mv := <-g.move:
			if err := g.validMovement(mv); err != nil {
				log.Println(err)
				continue
			}
			b := g.blocks[mv.y][mv.x]
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

func (g *game) validMovement(mv *move) error {
	if mv.p != g.turn {
		return errors.New("move is not belong to player who has turn")
	}
	b := g.blocks[mv.y][mv.x]
	if b.s != 0 {
		return errors.New("move over the non-empty block")
	}
	if !g.reachedBlock(mv) {
		return errors.New("move over the unreached block")
	}
	return nil
}

func (g *game) reachedBlock(mv *move) bool {
	if (!mv.p && g.scoreP1 == 0) || (mv.p && g.scoreP2 == 0) {
		return true
	}
	y := mv.y - 1
	for y >= 0 {
		s := g.blocks[y][mv.x].s
		if (s == 1 && !mv.p) || (s == 2 && mv.p) {
			return true
		}
		if s == 0 {
			y--
			continue
		}
		break
	}
	y = mv.y + 1
	for y < g.height {
		s := g.blocks[y][mv.x].s
		if (s == 1 && !mv.p) || (s == 2 && mv.p) {
			return true
		}
		if s == 0 {
			y++
			continue
		}
		break
	}
	x := mv.x - 1
	for x >= 0 {
		s := g.blocks[mv.y][x].s
		if (s == 1 && !mv.p) || (s == 2 && mv.p) {
			return true
		}
		if s == 0 {
			x--
			continue
		}
		break
	}
	x = mv.x + 1
	for x < g.width {
		s := g.blocks[mv.y][x].s
		if (s == 1 && !mv.p) || (s == 2 && mv.p) {
			return true
		}
		if s == 0 {
			x++
			continue
		}
		break
	}
	return false
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
