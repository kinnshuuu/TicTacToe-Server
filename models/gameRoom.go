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

func Make2D[T any](n, m int) [][]T {
	matrix := make([][]T, n)
	rows := make([]T, n*m)
	for i, startRow := 0, 0; i < n; i, startRow = i+1, startRow+m {
		endRow := startRow + m
		matrix[i] = rows[startRow:endRow:endRow]
	}
	return matrix
}

func NewGameRoom(player1 *Client, player2 *Client, roomId int) *GameRoom {
	board := Make2D[int](3, 3)
	return &GameRoom{
		RoomId:    roomId,
		GameState: constants.CONNECTED,
		Player1:   player1,
		Player2:   player2,
		Board:     &Board{board: board},
	}
}

func (g *GameRoom) SendStartPlayingSignal() {
	data1 := MessageToSend{State: constants.CONNECTED, Data: "0", Turn: 1, CommandType: int(constants.COMMAND_TYPE_MOVE)}
	data2 := MessageToSend{State: constants.CONNECTED, Data: "0", Turn: 0, CommandType: int(constants.COMMAND_TYPE_MOVE)}
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
