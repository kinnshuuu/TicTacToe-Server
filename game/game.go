package game

import (
	"TicTacToe-Server/models"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var playerId = 1
var roomId = 1

type Game struct {
	WaitingRoom *models.WaitingRoom
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
			player1 := g.WaitingRoom.Clients[0]
			g.WaitingRoom.Clients = g.WaitingRoom.Clients[1:]
			player2 := g.WaitingRoom.Clients[0]
			g.WaitingRoom.Clients = g.WaitingRoom.Clients[1:]
			gameRoom := models.NewGameRoom(player1, player2, roomId)
			log.Printf("New Game Room Created with Id: %d and having players with Player ID : %d and %d", roomId, player1.PlayerId, player2.PlayerId)
			roomId++
			go gameRoom.SendStartPlayingSignal()
		}
	}
}

func NewGame() *Game {
	waitingRoom := models.NewWaitingRoom()
	return &Game{
		WaitingRoom: waitingRoom,
		GameRooms:   make([]*models.GameRoom, 0),
	}
}
