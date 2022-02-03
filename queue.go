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
				player.sendMessage(joinedMessage())
			} else {
				log.Println("Player is already in queue")
			}
		case player := <-q.leave:
			if _, ok := q.players[player]; ok {
				delete(q.players, player)
				player.sendMessage(leftMessage())
			} else {
				log.Println("Player is already not in queue")
			}
		default:
			for len(q.players) > 1 {
				p1 := q.popPlayer()
				p2 := q.popPlayer()
				g := &game{
					player1: p1,
					player2: p2,
				}
				p1.game = g
				p2.game = g
				p1.sendMessage(matchedMessage())
				p2.sendMessage(matchedMessage())
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
