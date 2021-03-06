package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type player struct {
	game *game
	conn *websocket.Conn
	send chan []byte
}

func (p *player) processMessage(msg []byte) {
	m := &message{}
	err := json.Unmarshal(msg, m)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(m)
	switch m.Action {
	case actionJoin:
		gameQueue.join <- p
	case actionLeave:
		gameQueue.leave <- p
	case actionMove:
		if p.game == nil {
			err := errors.New("player has not any ongoing game")
			log.Println(err)
			return
		}
		c, err := moveData(m)
		if err != nil {
			log.Println(err)
			return
		}
		mv := &move{
			x: c.X,
			y: c.Y,
			p: p.game.player1 != p,
		}
		p.game.move <- mv
	case actionEnd:
		if p.game == nil {
			err := errors.New("player has not any ongoing game")
			log.Println(err)
			return
		}
		r, err := reasonData(m)
		if err != nil {
			log.Println(err)
			return
		}
		err = p.game.sendToEndGame(p, r)
		if err != nil {
			log.Println(err)
		}
	}
}

func (p *player) sendMessage(m *message) {
	msg, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return
	}
	p.send <- msg
}

func (p *player) startMessageReading() {
	defer func() {
		err := p.conn.Close()
		if err != nil {
			log.Printf("error: %v", err)
		}
	}()
	p.conn.SetReadLimit(maxMessageSize)
	p.conn.SetReadDeadline(time.Now().Add(pongWait))
	p.conn.SetPongHandler(func(string) error { p.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		p.processMessage(message)
	}
}

func (p *player) startMessageWriting() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := p.conn.Close()
		if err != nil {
			log.Printf("error: %v", err)
		}
	}()
	for {
		select {
		case message, ok := <-p.send:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				p.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := p.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("error: %v", err)
				return
			}
			w.Write(message)
			n := len(p.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-p.send)
			}

			if err := w.Close(); err != nil {
				log.Printf("error: %v", err)
				return
			}
		case <-ticker.C:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				// conclude player's ongoing game
				if p.game != nil {
					err = p.game.sendToEndGame(p, reasonDisconnect)
					if err != nil {
						log.Println(err)
					}
				}
				log.Printf("error: %v", err)
				return
			}
		}
	}
}
