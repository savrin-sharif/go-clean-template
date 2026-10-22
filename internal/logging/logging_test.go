package logging_test

import (
	"testing"

	"github.com/savrin-sharif/go-clean-template/internal/logging"
)

func TestRedactMap(t *testing.T) {
	t.Parallel()

	in := map[string]string{
		"username":      "alice",
		"password":      "secret",
		"access_token":  "abc",
		"Authorization": "Bearer x",
	}
	out := logging.RedactMap(in)
	if out["username"] != "alice" {
		t.Fatalf("username redacted unexpectedly: %q", out["username"])
	}
	for _, key := range []string{"password", "access_token", "Authorization"} {
		if out[key] != "[REDACTED]" {
			t.Fatalf("%s = %q, want [REDACTED]", key, out[key])
		}
	}
}
