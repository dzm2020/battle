package room

import "errors"

var (
	ErrWrongPhase       = errors.New("operation not allowed in current phase")
	ErrRoomClosed       = errors.New("room is closed")
	ErrUseCreatePVPRoom = errors.New("pvp dungeon: use room.CreatePVPRoom with self and enemy players")
	ErrNoSpatialGrid    = errors.New("room: spatial grid not set")
)
