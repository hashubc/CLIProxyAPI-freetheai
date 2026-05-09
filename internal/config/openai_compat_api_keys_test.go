package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenAICompatibilityAPIKeysShortcutLoadsAndSaves(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join(t.TempDir(), "config.yaml")
	input := `openai-compatibility:
  - name: "freetheai"
    base-url: "https://api.freetheai.xyz/v1"
    api-keys:
      - "  fta-key-1  "
      - ""
      - "fta-key-2"
    models:
      - name: "glm/glm-5.1"
        alias: "glm-5.1"
`
	if errWrite := os.WriteFile(configPath, []byte(input), 0o600); errWrite != nil {
		t.Fatalf("write config: %v", errWrite)
	}

	cfg, err := LoadConfigOptional(configPath, false)
	if err != nil {
		t.Fatalf("LoadConfigOptional() error = %v", err)
	}
	if len(cfg.OpenAICompatibility) != 1 {
		t.Fatalf("openai-compatibility entries = %d, want 1", len(cfg.OpenAICompatibility))
	}
	got := cfg.OpenAICompatibility[0].APIKeys
	if len(got) != 2 || got[0] != "fta-key-1" || got[1] != "fta-key-2" {
		t.Fatalf("api-keys = %#v, want [fta-key-1 fta-key-2]", got)
	}

	if err = SaveConfigPreserveComments(configPath, cfg); err != nil {
		t.Fatalf("SaveConfigPreserveComments() error = %v", err)
	}
	saved, errRead := os.ReadFile(configPath)
	if errRead != nil {
		t.Fatalf("read saved config: %v", errRead)
	}
	text := string(saved)
	if !strings.Contains(text, "api-keys:") {
		t.Fatalf("saved config removed api-keys:\n%s", text)
	}
	if !strings.Contains(text, "fta-key-1") || !strings.Contains(text, "fta-key-2") {
		t.Fatalf("saved config missing API keys:\n%s", text)
	}
}
