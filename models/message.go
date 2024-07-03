package models

import (
	"TicTacToe-GolangServer/constants"
)

type SendMessage struct {
	CommandType int                     `json:"command_type"`
	State       constants.GameRoomState `json:"state"`
	Data        string                  `json:"data"`
	Turn        int                     `json:"turn"`
}

type ReceivedMessage struct {
	CommandType int    `json:"command_type"`
	Data        string `json:"data"`
}
