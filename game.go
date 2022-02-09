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

const reasonDisconnect = "disconnect"
const reasonResign = "resign"
const reasonComplete = "complete"

type game struct {
	turn    bool
	move    chan *move
	end     chan *end
	blocks  [][]*block
	width   int
	height  int
	player1 *player
	player2 *player
	scoreP1 int
	scoreP2 int
	reason  string
	winner  int
}

func newGame(p1, p2 *player) (*game, []*coordinate, []*coordinate, []*coordinate) {
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
	/*
		coordinates := randomCoordinate(defaultMinWallCount, defaultMaxWallCount, defaultWidth, defaultHeight)
		for _, coordinate := range coordinates {
			blocks[coordinate.Y][coordinate.X].s = -1
		}
	*/

	c1im1 := &coordinate{
		X: defaultWidth/2 + 1,
		Y: defaultHeight / 2,
	}
	c2im1 := &coordinate{
		X: defaultWidth / 2,
		Y: defaultHeight/2 + 1,
	}
	c1im2 := &coordinate{
		X: defaultWidth/2 + 1,
		Y: defaultHeight/2 + 1,
	}
	c2im2 := &coordinate{
		X: defaultWidth / 2,
		Y: defaultHeight / 2,
	}

	blocks[c1im1.Y][c1im1.X].s = 1
	blocks[c2im1.Y][c2im1.X].s = 1
	blocks[c1im2.Y][c1im1.X].s = 2
	blocks[c2im2.Y][c2im1.X].s = 2

	g := &game{
		move:    make(chan *move),
		end:     make(chan *end),
		blocks:  blocks,
		width:   defaultWidth,
		height:  defaultHeight,
		player1: p1,
		player2: p2,
	}
	//return g, coordinates
	return g, make([]*coordinate, 0), []*coordinate{c1im1, c2im1}, []*coordinate{c1im2, c2im2}
}

func (g *game) startGame() {
	log.Printf("game has been started")
	for {
		select {
		case e := <-g.end:
			switch e.reason {
			case reasonDisconnect, reasonResign:
				d := make(map[string]interface{})
				d[dataReason] = e.reason
				d[dataWinner] = !e.player
				m := &message{
					Action: actionEnded,
					Data:   d,
				}
				g.player1.sendMessage(m)
				g.player2.sendMessage(m)
				g.player1.game = nil
				g.player2.game = nil
				if e.player {
					g.winner = 1
				} else {
					g.winner = 2
				}
			case reasonComplete:
				d := make(map[string]interface{})
				d[dataReason] = e.reason
				d[dataWinner] = !e.player
				d[dataScoreP1] = g.scoreP1
				d[dataScoreP2] = g.scoreP2
				d[dataWinner] = g.scoreP1 < g.scoreP2
				m := &message{
					Action: actionEnded,
					Data:   d,
				}
				g.player1.sendMessage(m)
				g.player2.sendMessage(m)
				g.player1.game = nil
				g.player2.game = nil
				if g.scoreP1 > g.scoreP2 {
					g.winner = 1
				} else if g.scoreP1 < g.scoreP2 {
					g.winner = 2
				} else {
					g.winner = 3
				}
			default:
				continue
			}
			return
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

func (g *game) sendToEndGame(p *player, reason string) error {
	if g.winner != 0 {
		return errors.New("game has already been ended")
	}
	switch reason {
	case reasonDisconnect, reasonResign:
		e := &end{
			reason: reason,
			player: g.player1 != p,
		}
		g.end <- e
		return nil
	default:
		return errors.New("unknown reason")
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

type end struct {
	reason string
	player bool
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
