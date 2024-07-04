package models

import (
	"TicTacToe-Server/constants"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn      *websocket.Conn
	PlayerId  int
	PieceType int
}

func (c *Client) SendPlayerWaitMessage() {
	data := MessageToSend{State: constants.WAITING, Data: "0", Turn: 0, CommandType: int(constants.COMMAND_TYPE_MOVE)}
	err := c.SendWebSocketMessage(data)
	if err != nil {
		log.Println("Failed to send message from client:", err)
	}
}

func (c *Client) handleRoomPlayer(gameRoom *GameRoom) {
	for {
		var msg ReceivedMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Failed to read message from client:", err)
			break
		}
		gameRoom.Board.PlaceMove(c.PlayerId, msg)
		c.sendResponseToOtherClients(gameRoom, msg)
		isOver, result := gameRoom.Board.CheckForWin(c.PlayerId, msg)
		if isOver == 1 {
			c.BroadCastMessage(gameRoom, result)
		}
	}

	c.RemoveRoom(gameRoom)
}

func (c *Client) BroadCastMessage(gameRoom *GameRoom, result int) {
	data := MessageToSend{State: constants.PLAYING, Data: fmt.Sprint(result), Turn: 0, CommandType: int(constants.COMMAND_TYPE_RESULT)}
	if c.PlayerId == gameRoom.Player1.PlayerId {
		err := gameRoom.Player1.SendWebSocketMessage(data)
		if err != nil {
			log.Println("Failed to send message to client:", err)
		}
		if result == constants.STATE_WIN {
			data.Data = fmt.Sprint(constants.STATE_LOSE)
			err = gameRoom.Player2.SendWebSocketMessage(data)
			if err != nil {
				log.Println("Failed to send message to client:", err)
			}
		}
	} else {
		err := gameRoom.Player2.SendWebSocketMessage(data)
		if err != nil {
			log.Println("Failed to send message to client:", err)
		}
		if result == constants.STATE_WIN {
			data.Data = fmt.Sprint(constants.STATE_LOSE)
			err = gameRoom.Player1.SendWebSocketMessage(data)
			if err != nil {
				log.Println("Failed to send message to client:", err)
			}
		}
	}

}

func (c *Client) RemoveRoom(gameRoom *GameRoom) {
	data := MessageToSend{State: constants.WAITING, Data: "0", Turn: 0, CommandType: int(constants.COMMAND_TYPE_DISCONNECTED)}

	if gameRoom.Player1 != nil && c.PlayerId == gameRoom.Player1.PlayerId {
		log.Println("Hiii")
		err := gameRoom.Player2.SendWebSocketMessage(data)
		if err != nil {
			log.Println("Failed to send message to client:", err)
		}
	} else if gameRoom.Player2 != nil && c.PlayerId == gameRoom.Player2.PlayerId {
		log.Println("Byee")
		err := gameRoom.Player1.SendWebSocketMessage(data)
		if err != nil {
			log.Println("Failed to send message to client:", err)
		}
	}
	gameRoom.Player1 = nil
	gameRoom.Player2 = nil
}

func (c *Client) sendResponseToOtherClients(gameRoom *GameRoom, msg ReceivedMessage) {

	data := MessageToSend{CommandType: msg.CommandType, State: constants.PLAYING, Data: string(msg.Data), Turn: 1}

	if c.PlayerId == gameRoom.Player1.PlayerId {
		err := gameRoom.Player2.SendWebSocketMessage(data)
		if err != nil {
			log.Println("Failed to send message to client:", err)
		}
	} else {
		err := gameRoom.Player1.SendWebSocketMessage(data)
		if err != nil {
			log.Println("Failed to send message to client:", err)
		}
	}
}

func (c *Client) SendWebSocketMessage(data MessageToSend) error {
	err := websocket.WriteJSON(c.Conn, data)
	return err
}
