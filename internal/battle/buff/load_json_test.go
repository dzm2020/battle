package buff

import (
	"strings"
	"testing"
)

func TestLoadDescriptorsFromJSON(t *testing.T) {
	const raw = `[
	  {"id":1,"maxStacks":2,"policy":0,"durationFrames":5,"effects":[{"kind":0,"armorDeltaPerStack":3}]}
	]`
	reg := NewRegistry()
	if err := LoadDescriptorsFromJSON([]byte(strings.TrimSpace(raw)), reg); err != nil {
		t.Fatal(err)
	}
	d, ok := reg.Get(1)
	if !ok || d.MaxStacks != 2 || len(d.Effects) != 1 || d.Effects[0].ArmorDeltaPerStack != 3 {
		t.Fatalf("descriptor: %+v", d)
	}
}
