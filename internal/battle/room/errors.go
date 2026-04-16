package room

import "errors"

var (
	ErrInvalidEntity   = errors.New("entity is nil")
	ErrRoomExists      = errors.New("room id already exists")
	ErrRoomNotFound    = errors.New("room not found")
	ErrRoomFull        = errors.New("room is full")
	ErrWrongPhase      = errors.New("operation not allowed in current phase")
	ErrDuplicatePlayer = errors.New("player already in room")
	ErrRoomClosed      = errors.New("room is closed")
	ErrNoPlayers       = errors.New("no players in room")
)
