package websocket

import (
	"context"
	"errors"
	"testing"

	"git.sr.ht/~rockorager/go-jmap"
)

func TestDialH2ConnectReturnsUnsupported(t *testing.T) {
	ctx := context.Background()
	client := &jmap.Client{SessionEndpoint: "https://example.invalid/jmap"}
	_, err := DialH2Connect(ctx, client)
	if err == nil {
		t.Fatal("expected error without H2 CONNECT support")
	}
	if !errors.Is(err, ErrH2ConnectUnsupported) {
		t.Fatalf("expected ErrH2ConnectUnsupported, got %v", err)
	}
}
