package skill

import (
	"strings"
	"testing"
)

func TestLoadCatalogConfigFromJSON(t *testing.T) {
	const raw = `[{"id":7,"resource":0,"cost":0,"cooldownFrames":10,"target":0,"castFrames":0,"effects":[]}]`
	catalogConfig := NewCatalogConfig()
	if err := LoadCatalogConfigFromJSON([]byte(strings.TrimSpace(raw)), catalogConfig); err != nil {
		t.Fatal(err)
	}
	sk, ok := catalogConfig.Get(7)
	if !ok || sk.CooldownFrames != 10 {
		t.Fatalf("skill config: %+v", sk)
	}
}
