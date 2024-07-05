package models

import (
	"TicTacToe-Server/constants"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn       *websocket.Conn
	PlayerId   int
	PieceType  int
	WriteMutex sync.Mutex
}

func (c *Client) SendPlayerWaitMessage() {
	data := MessageToSend{State: constants.WAITING, Data: "0", Turn: 0, CommandType: int(constants.COMMAND_TYPE_MOVE)}
	err := c.SendWebSocketMessage(data)
	if err != nil {
		log.Printf("Failed to send wait message to client with PlayerId - %d: %v", c.PlayerId, err)
	}
}

func (c *Client) handleRoomPlayer(gameRoom *GameRoom) {
	for {
		var msg ReceivedMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Failed to read message from client with PlayerId - %d: %v", c.PlayerId, err)
			break
		}
		gameRoom.Board.PlaceMove(c.PlayerId, msg)
		c.sendResponseToOtherClients(gameRoom, msg)
		isOver, result := gameRoom.Board.CheckForWin(c.PlayerId, msg)
		if result == constants.STATE_WIN {
			log.Printf("Player with Id %d won", c.PlayerId)
		}
		if isOver == 1 {
			c.BroadCastMessage(gameRoom, result)
		}
	}

	c.RemoveRoom(gameRoom)
}

func (c *Client) BroadCastMessage(gameRoom *GameRoom, result int) {
	data := MessageToSend{State: constants.PLAYING, Data: fmt.Sprint(result), Turn: 0, CommandType: int(constants.COMMAND_TYPE_RESULT)}
	var err error

	if c.PlayerId == gameRoom.Player1.PlayerId {
		err = gameRoom.Player1.SendWebSocketMessage(data)
		if err != nil {
			log.Printf("Failed to send result message to client with PlayerId - %d: %v", gameRoom.Player1.PlayerId, err)
		}
		if result == constants.STATE_WIN {
			data.Data = fmt.Sprint(constants.STATE_LOSE)
			err = gameRoom.Player2.SendWebSocketMessage(data)
			if err != nil {
				log.Printf("Failed to send result message to client with PlayerId - %d: %v", gameRoom.Player2.PlayerId, err)
			}
		}
	} else {
		err = gameRoom.Player2.SendWebSocketMessage(data)
		if err != nil {
			log.Printf("Failed to send result message to client with PlayerId - %d: %v", gameRoom.Player2.PlayerId, err)
		}
		if result == constants.STATE_WIN {
			data.Data = fmt.Sprint(constants.STATE_LOSE)
			err = gameRoom.Player1.SendWebSocketMessage(data)
			if err != nil {
				log.Printf("Failed to send result message to client with PlayerId - %d: %v", gameRoom.Player1.PlayerId, err)
			}
		}
	}
}

func (c *Client) RemoveRoom(gameRoom *GameRoom) {
	data := MessageToSend{State: constants.WAITING, Data: "0", Turn: 0, CommandType: int(constants.COMMAND_TYPE_DISCONNECTED)}

	if gameRoom.Player1 != nil && c.PlayerId == gameRoom.Player1.PlayerId {
		log.Printf("Client %d left the room", c.PlayerId)
		err := gameRoom.Player2.SendWebSocketMessage(data)
		if err != nil {
			log.Printf("Failed to notify client %d about disconnection: %v", gameRoom.Player2.PlayerId, err)
		}
	} else if gameRoom.Player2 != nil && c.PlayerId == gameRoom.Player2.PlayerId {
		log.Printf("Client %d left the room", c.PlayerId)
		err := gameRoom.Player1.SendWebSocketMessage(data)
		if err != nil {
			log.Printf("Failed to notify client %d about disconnection: %v", gameRoom.Player1.PlayerId, err)
		}
	}
	gameRoom.Player1 = nil
	gameRoom.Player2 = nil
}

func (c *Client) sendResponseToOtherClients(gameRoom *GameRoom, msg ReceivedMessage) {
	data := MessageToSend{CommandType: msg.CommandType, State: constants.PLAYING, Data: string(msg.Data), Turn: 1}
	var err error

	if c.PlayerId == gameRoom.Player1.PlayerId {
		err = gameRoom.Player2.SendWebSocketMessage(data)
		if err != nil {
			log.Printf("Failed to send message to client %d: %v", gameRoom.Player2.PlayerId, err)
		}
	} else {
		err = gameRoom.Player1.SendWebSocketMessage(data)
		if err != nil {
			log.Printf("Failed to send message to client %d: %v", gameRoom.Player1.PlayerId, err)
		}
	}
}

func (c *Client) SendWebSocketMessage(data MessageToSend) error {
	c.WriteMutex.Lock()
	defer c.WriteMutex.Unlock()

	err := c.Conn.WriteJSON(data)
	if err != nil {
		return fmt.Errorf("failed to send WebSocket message: %w", err)
	}
	return nil
}
