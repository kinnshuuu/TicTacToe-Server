package models

import (
	"TicTacToe-Server/constants"
	"log"
)

type GameRoom struct {
	RoomId    int
	GameState constants.GameRoomState
	Player1   *Client
	Player2   *Client
	Board     *Board
}

func MakeBoard(n, m int) [][]int {
	matrix := make([][]int, n)
	rows := make([]int, n*m)
	for i, startRow := 0, 0; i < n; i, startRow = i+1, startRow+m {
		endRow := startRow + m
		matrix[i] = rows[startRow:endRow:endRow]
	}
	return matrix
}

func NewGameRoom(player1 *Client, player2 *Client, roomId int) *GameRoom {
	board := MakeBoard(3, 3)
	return &GameRoom{
		RoomId:    roomId,
		GameState: constants.CONNECTED,
		Player1:   player1,
		Player2:   player2,
		Board:     &Board{board: board},
	}
}

func (g *GameRoom) SendStartPlayingSignal() {
	data1 := MessageToSend{State: constants.CONNECTED, Data: "0", Turn: 1, CommandType: int(constants.COMMAND_TYPE_CONNECTED)}
	data2 := MessageToSend{State: constants.CONNECTED, Data: "0", Turn: 0, CommandType: int(constants.COMMAND_TYPE_CONNECTED)}
	err := g.Player1.SendWebSocketMessage(data1)
	if err != nil {
		log.Println("Failed to send message from client:", err)
	}
	err = g.Player2.SendWebSocketMessage(data2)
	if err != nil {
		log.Println("Failed to send message from client:", err)
	}
	go g.Player1.handleRoomPlayer(g)
	go g.Player2.handleRoomPlayer(g)
}
