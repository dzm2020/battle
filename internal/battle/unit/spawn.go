package unit

import "battle/ecs"

func SpawnUnitFromConfig(w *ecs.World, unitID string) (ecs.Entity, error) {
	_, _ = unitID, unitID
	return w.CreateEntity(), nil
}
