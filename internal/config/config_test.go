package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Koshsky/erp-backend/internal/config"
)

// writeConfig writes a minimal config.yaml into a temp dir and returns its path.
func writeConfig(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// TestLoadRejectsInvalidRateLimit guards against the silent-disable/hard-block
// misconfigurations (I9): an enabled rate limit with requests_per_second <= 0
// or burst <= 0 must fail startup with a clear error instead of silently
// disabling the wall or 429ing every request. Setenv is required because the
// test drives the real Load() entrypoint.
func TestLoadRejectsInvalidRateLimit(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://test")
	t.Setenv("JWT_SECRET_KEY", "12345678901234567890123456789012")

	for _, tc := range []struct {
		name string
		yaml string
		want string
	}{
		{
			name: "public limit with zero requests_per_second",
			yaml: "rate_limiting:\n  enabled: true\n  requests_per_second: 0\n  burst: 20\n",
			want: "requests_per_second",
		},
		{
			name: "public limit with zero burst",
			yaml: "rate_limiting:\n  enabled: true\n  requests_per_second: 20\n  burst: 0\n",
			want: "burst",
		},
		{
			name: "user limit with zero requests_per_second",
			yaml: "user_rate_limiting:\n  enabled: true\n  requests_per_second: 0\n  burst: 20\n",
			want: "requests_per_second",
		},
		{
			name: "user limit with zero burst",
			yaml: "user_rate_limiting:\n  enabled: true\n  requests_per_second: 20\n  burst: 0\n",
			want: "burst",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("CONFIG_PATH", writeConfig(t, tc.yaml))
			_, err := config.Load()
			if err == nil {
				t.Fatal("expected an error for the invalid rate limit config")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %q, want it to mention %q", err, tc.want)
			}
		})
	}
}

// TestLoadAllowsDisabledRateLimits checks the disabled path: rate limit blocks
// that are off are not validated, so zero values pass without an error.
func TestLoadAllowsDisabledRateLimits(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://test")
	t.Setenv("JWT_SECRET_KEY", "12345678901234567890123456789012")
	t.Setenv("CONFIG_PATH", writeConfig(t, "rate_limiting:\n  enabled: false\nuser_rate_limiting:\n  enabled: false\n"))

	if _, err := config.Load(); err != nil {
		t.Fatalf("Load() error = %v, want nil for disabled rate limits", err)
	}
}
