package room

import "errors"

var (
	ErrWrongPhase       = errors.New("operation not allowed in current phase")
	ErrRoomClosed       = errors.New("room is closed")
	ErrUseCreatePVPRoom = errors.New("pvp dungeon: use room_builder.CreatePVPRoom with red and blue teams")
	ErrNoSpatialGrid    = errors.New("room: spatial grid not set")
)
