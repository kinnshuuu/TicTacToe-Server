package game

import (
	"TicTacToe-GolangServer/models"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var playerId = 1

type WaitingRoom struct {
	Clients []*models.Client
	Mutex   sync.Mutex
}

func (wr *WaitingRoom) AddPlayer(client *models.Client) {
	wr.Mutex.Lock()
	defer wr.Mutex.Unlock()

	wr.Clients = append(wr.Clients, client)
	fmt.Println(wr.Clients[len(wr.Clients)-1])
	go client.SendPlayerWaitMessage()
}

type Game struct {
	WaitingRoom *WaitingRoom
	GameRooms   []*models.GameRoom
}

func (g *Game) ServeClient(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}
	client := models.Client{Conn: conn, PlayerId: playerId}
	playerId++
	g.WaitingRoom.AddPlayer(&client)
}

func (g *Game) MatchOnlinePlayer() {
	for {
		if len(g.WaitingRoom.Clients) > 1 {
			player1 := g.WaitingRoom.Clients[len(g.WaitingRoom.Clients)-1]
			g.WaitingRoom.Clients = g.WaitingRoom.Clients[:len(g.WaitingRoom.Clients)-1]
			player2 := g.WaitingRoom.Clients[len(g.WaitingRoom.Clients)-1]
			gameRoom := models.NewGameRoom(player1, player2)
			g.WaitingRoom.Clients = g.WaitingRoom.Clients[:len(g.WaitingRoom.Clients)-1]
			go gameRoom.SendStartPlayingSignal()
		}
	}
}

func NewWaitingRoom() *WaitingRoom {
	return &WaitingRoom{
		Clients: make([]*models.Client, 0),
		Mutex:   sync.Mutex{},
	}
}

func NewGame() *Game {
	waitingRoom := NewWaitingRoom()
	return &Game{
		WaitingRoom: waitingRoom,
		GameRooms:   make([]*models.GameRoom, 0),
	}
}
