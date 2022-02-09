package main

import (
	"log"
)

type queue struct {
	players map[*player]interface{}
	join    chan *player
	leave   chan *player
}

func newQueue() *queue {
	return &queue{
		players: make(map[*player]interface{}),
		join:    make(chan *player),
		leave:   make(chan *player),
	}
}

func (q *queue) startMatching() {
	for {
		select {
		case player := <-q.join:
			if _, ok := q.players[player]; !ok {
				q.players[player] = nil
				m := &message{
					Action: actionJoined,
				}
				player.sendMessage(m)
			} else {
				log.Println("Player is already in queue")
			}
		case player := <-q.leave:
			if _, ok := q.players[player]; ok {
				delete(q.players, player)
				m := &message{
					Action: actionLeft,
				}
				player.sendMessage(m)
			} else {
				log.Println("Player is already not in queue")
			}
		default:
			for len(q.players) > 1 {
				p1 := q.popPlayer()
				p2 := q.popPlayer()
				g, c, imp1, imp2 := newGame(p1, p2)
				p1.game = g
				p2.game = g
				go g.startGame()
				dataP1 := make(map[string]interface{})
				dataP1[dataPlayer] = false
				dataP1[dataWalls] = c
				dataP1[dataInitialMoveP1] = imp1
				dataP1[dataInitialMoveP2] = imp2
				mP1 := &message{
					Action: actionMatched,
					Data:   dataP1,
				}
				p1.sendMessage(mP1)
				dataP2 := make(map[string]interface{})
				dataP2[dataPlayer] = true
				dataP2[dataWalls] = c
				dataP2[dataInitialMoveP1] = imp1
				dataP2[dataInitialMoveP2] = imp2
				mP2 := &message{
					Action: actionMatched,
					Data:   dataP2,
				}
				p2.sendMessage(mP2)
			}
		}
	}
}

func (q *queue) popPlayer() *player {
	for p := range q.players {
		delete(q.players, p)
		return p
	}
	return nil
}
