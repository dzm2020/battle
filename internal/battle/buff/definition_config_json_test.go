package buff

import (
	"strings"
	"testing"
)

func TestLoadDefinitionConfigFromJSON(t *testing.T) {
	const raw = `[
	  {"id":1,"maxStacks":2,"policy":0,"durationFrames":5,"effects":[{"kind":0,"armorDeltaPerStack":3}]}
	]`
	config := NewDefinitionConfig()
	if err := LoadDefinitionConfigFromJSON([]byte(strings.TrimSpace(raw)), config); err != nil {
		t.Fatal(err)
	}
	d, ok := config.Get(1)
	if !ok || d.MaxStacks != 2 || len(d.Effects) != 1 || d.Effects[0].ArmorDeltaPerStack != 3 {
		t.Fatalf("descriptor: %+v", d)
	}
}
