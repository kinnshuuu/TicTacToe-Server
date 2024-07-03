package constants

type GameRoomState int

const (
	WAITING GameRoomState = iota + 8
	CONNECTED
	PLAYING
)

type CommandType int

const (
	COMMAND_TYPE_PAUSE CommandType = iota + 60
	COMMAND_TYPE_MOVE
	COMMAND_TYPE_DISCONNECTED
	COMMAND_TYPE_RESULT
)

// const COMMAND_TYPE_PAUSE CommandType = 60
// const COMMAND_TYPE_MOVE CommandType = 61

type PieceType int

const (
	X_PIECE = iota + 1
	O_PIECE
)

const STATE_WIN int = 14
const STATE_LOSE int = 15
const STATE_DRAW int = 18
