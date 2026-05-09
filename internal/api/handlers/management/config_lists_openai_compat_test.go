package management

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v6/internal/config"
)

func TestPatchOpenAICompatSupportsAPIKeysShortcut(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	h := &Handler{
		cfg: &config.Config{
			OpenAICompatibility: []config.OpenAICompatibility{
				{
					Name:    "freetheai",
					BaseURL: "https://api.freetheai.xyz/v1",
				},
			},
		},
		configFilePath: writeTestConfigFile(t),
	}

	body := []byte(`{"name":"freetheai","value":{"api-keys":["  fta-key-1  ","","fta-key-2"]}}`)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPatch, "/v0/management/openai-compatibility", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.PatchOpenAICompat(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	got := h.cfg.OpenAICompatibility[0].APIKeys
	if len(got) != 2 || got[0] != "fta-key-1" || got[1] != "fta-key-2" {
		t.Fatalf("api-keys = %#v, want [fta-key-1 fta-key-2]", got)
	}
}

func TestOpenAICompatibilityWithAuthIndexPreservesAPIKeysShortcut(t *testing.T) {
	t.Parallel()

	h := &Handler{
		cfg: &config.Config{
			OpenAICompatibility: []config.OpenAICompatibility{
				{
					Name:    "freetheai",
					BaseURL: "https://api.freetheai.xyz/v1",
					APIKeys: []string{
						"  fta-key-1  ",
						"",
						"fta-key-2",
					},
					Models: []config.OpenAICompatibilityModel{
						{Name: "glm/glm-5.1", Alias: "glm-5.1"},
					},
				},
			},
		},
	}

	got := h.openAICompatibilityWithAuthIndex()
	if len(got) != 1 {
		t.Fatalf("entries = %d, want 1", len(got))
	}
	if got[0].AuthIndex != "" {
		t.Fatalf("auth-index = %q, want empty for api-keys shortcut", got[0].AuthIndex)
	}
	if len(got[0].APIKeys) != 2 || got[0].APIKeys[0] != "fta-key-1" || got[0].APIKeys[1] != "fta-key-2" {
		t.Fatalf("api-keys = %#v, want [fta-key-1 fta-key-2]", got[0].APIKeys)
	}
}
