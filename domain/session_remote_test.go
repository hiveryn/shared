package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRemoteSessionJSONDoesNotExposeCredentials(t *testing.T) {
	b, err := json.Marshal(Session{Machine: "box", SSH: "private-alias", RemoteToken: "private-token", RemotePort: 34567, Connection: "disconnected"})
	if err != nil {
		t.Fatal(err)
	}
	value := string(b)
	for _, secret := range []string{"private-alias", "private-token", "34567", "remote_token", "remote_port"} {
		if strings.Contains(value, secret) {
			t.Fatalf("private remote metadata exposed: %s", value)
		}
	}
	if !strings.Contains(value, `"machine":"box"`) || !strings.Contains(value, `"connection":"disconnected"`) {
		t.Fatal(value)
	}
}
