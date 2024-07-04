package models

import (
	"TicTacToe-Server/constants"
	"strconv"
)

type Board struct {
	board [][]int
}

func (b *Board) CheckForWin(playerId int, msg ReceivedMessage) (int, int) {
	if msg.CommandType != int(constants.COMMAND_TYPE_MOVE) {
		return 0, 0
	}
	tag, _ := strconv.Atoi(msg.Data)
	x, y := getXandY(tag)

	// Check Row
	win := b.CheckSectionForWin(0, y, 1, y, 2, y, playerId)
	if win == 1 {
		return 1, 14
	}

	// Check column
	win = b.CheckSectionForWin(x, 0, x, 1, x, 2, playerId)
	if win == 1 {
		return 1, 14
	}

	// Check Diagonal
	win = b.CheckSectionForWin(0, 0, 1, 1, 2, 2, playerId)
	if win == 1 {
		return 1, 14
	}

	// Check Other Diagonal
	win = b.CheckSectionForWin(0, 2, 1, 1, 2, 0, playerId)
	if win == 1 {
		return 1, 14
	}

	draw := b.CheckForDraw()
	if draw == 1 {
		return 1, 18
	}
	return 0, 0
}

func (b *Board) CheckForDraw() int {
	for i := range b.board {
		for j := range b.board[0] {
			if b.board[i][j] == 0 {
				return 0
			}
		}
	}
	return 1
}

func (b *Board) CheckSectionForWin(x1, y1, x2, y2, x3, y3, playerId int) int {
	if (b.board[x1][y1] == b.board[x2][y2]) && (b.board[x2][y2] == b.board[x3][y3]) && (b.board[x3][y3] == playerId) {
		return 1
	}
	return 0
}

func getXandY(tag int) (int, int) {
	x := tag / 3
	y := tag % 3
	return x, y
}

func (b *Board) PlaceMove(playerId int, msg ReceivedMessage) {
	if msg.CommandType == int(constants.COMMAND_TYPE_MOVE) {
		tag, _ := strconv.Atoi(msg.Data)
		x, y := getXandY(tag)
		b.board[x][y] = playerId
	}
}
