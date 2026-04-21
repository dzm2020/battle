package component

import (
	"testing"

	"battle/ecs"
)

func TestRegisterComponent_CreateComponent(t *testing.T) {
	RegisterComponent("__test_empty", func() ecs.Component {
		t.Helper()
		return &ThreatBook{}
	})
	c, ok := CreateComponent("__test_empty")
	if !ok || c == nil {
		t.Fatal()
	}
	if _, ok := CreateComponent("__no_such_factory__"); ok {
		t.Fatal()
	}
}
