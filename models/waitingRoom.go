package models

import (
	"log"
	"sync"
)

type WaitingRoom struct {
	Clients []*Client
	Mutex   sync.Mutex
}

func (wr *WaitingRoom) AddPlayer(client *Client) {
	wr.Mutex.Lock()
	defer wr.Mutex.Unlock()

	wr.Clients = append(wr.Clients, client)
	log.Printf("New Player joined with Player ID: %d", client.PlayerId)
	go client.SendPlayerWaitMessage()
}

func NewWaitingRoom() *WaitingRoom {
	return &WaitingRoom{
		Clients: make([]*Client, 0),
		Mutex:   sync.Mutex{},
	}
}
