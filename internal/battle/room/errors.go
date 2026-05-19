package room

import "errors"

var (
	ErrWrongPhase       = errors.New("operation not allowed in current phase")
	ErrRoomClosed       = errors.New("room is closed")
	ErrInvalidRoomSpec  = errors.New("room: invalid or unknown dungeon in RoomSpec")
	ErrUseCreatePVPRoom = errors.New("pvp dungeon: RoomSpec.Enemy is required")
	ErrNoSpatialGrid    = errors.New("room: spatial grid not set")
)
