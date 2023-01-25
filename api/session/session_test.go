package session_test

import (
	"github.com/branislavlazic/midnight/api/session"
	"testing"
)

func FuzzVerifySessionID(f *testing.F) {
	f.Add("secR3t")
	f.Fuzz(func(t *testing.T, secret string) {
		ID, err := session.GenerateSessionID(secret)
		if err != nil {
			t.Fatalf("failed to generate session id for secret %s", secret)
		}
		ok := session.VerifySessionID(ID, secret)
		if !ok {
			t.Fatalf("failed to verify session id for secret %s", secret)
		}
	})
}
