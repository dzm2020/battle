package room

import "errors"

var (
	ErrInvalidEntity   = errors.New("entity not in room world or invalid")
	ErrDuplicateEntity = errors.New("entity already assigned to another player slot")
	ErrPlayerNotFound  = errors.New("player not in room")
	ErrRoomExists      = errors.New("room id already exists")
	ErrRoomNotFound    = errors.New("room not found")
	ErrRoomFull        = errors.New("room is full")
	ErrWrongPhase      = errors.New("operation not allowed in current phase")
	ErrDuplicatePlayer = errors.New("player already in room")
	ErrRoomClosed      = errors.New("room is closed")
	ErrNoPlayers       = errors.New("no players in room")

	// PVP：须使用 [room_builder.CreatePVPRoom]；副本配置 [config.DungeonConfig.Type] 须为 [room_builder.DungeonTypePVP]。
	ErrUseCreatePVPRoom     = errors.New("pvp dungeon: use room_builder.CreatePVPRoom with red and blue teams")
	ErrPVPRequiresBothTeams = errors.New("pvp room requires at least one player on each team")
	ErrDungeonNotPVPType    = errors.New("dungeon config type is not PVP")

	ErrNoSpatialGrid = errors.New("room: spatial grid not set")
	ErrGridFull      = errors.New("room: grid has no free cell")

	ErrEntityNotOnGrid = errors.New("room: entity not registered on spatial grid")
)
