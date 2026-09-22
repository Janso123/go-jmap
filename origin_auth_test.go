package jmap

import (
	"net/http"
	"net/url"
	"testing"
)

func TestOriginOfURLNormalizesDefaultPorts(t *testing.T) {
	t.Parallel()
	httpsDefault, _ := url.Parse("https://mail.example.com:443/session")
	httpsBare, _ := url.Parse("https://mail.example.com/session")
	if originOfURL(httpsDefault) != originOfURL(httpsBare) {
		t.Fatalf("https default port: %q vs %q", originOfURL(httpsDefault), originOfURL(httpsBare))
	}
	if originOfURL(httpsBare) != "https://mail.example.com" {
		t.Fatalf("https origin = %q", originOfURL(httpsBare))
	}
	httpAlt, _ := url.Parse("http://127.0.0.1:8080/session")
	if originOfURL(httpAlt) != "http://127.0.0.1:8080" {
		t.Fatalf("explicit port origin = %q", originOfURL(httpAlt))
	}
}

func TestCheckWebSocketURL(t *testing.T) {
	t.Parallel()
	same := &Client{SessionEndpoint: "https://mail.example.com/jmap/session"}
	if err := same.CheckWebSocketURL("wss://mail.example.com/jmap/ws/", false); err != nil {
		t.Fatalf("same host wss: %v", err)
	}
	if err := same.CheckWebSocketURL("wss://mail.example.com:443/jmap/ws/", false); err != nil {
		t.Fatalf("default port wss: %v", err)
	}
	if err := same.CheckWebSocketURL("wss://other.example.com/jmap/ws/", false); err == nil {
		t.Fatal("foreign origin was allowed")
	}
	if err := same.CheckWebSocketURL("wss://other.example.com/jmap/ws/", true); err != nil {
		t.Fatalf("AllowForeignOrigin: %v", err)
	}
	if err := same.CheckWebSocketURL("ws://example.invalid/jmap/ws/", true); err == nil {
		t.Fatal("non-loopback ws was allowed")
	}
	if err := same.CheckWebSocketURL("ws://127.0.0.1:8080/jmap/ws/", false); err == nil {
		t.Fatal("loopback ws on another origin was allowed")
	}

	loop := &Client{SessionEndpoint: "http://127.0.0.1:8080/jmap/session"}
	if err := loop.CheckWebSocketURL("ws://127.0.0.1:8080/jmap/ws/", false); err != nil {
		t.Fatalf("loopback same origin: %v", err)
	}
	unset := &Client{}
	if err := unset.CheckWebSocketURL("ws://127.0.0.1:9/jmap/ws/", false); err != nil {
		t.Fatalf("loopback without session origin: %v", err)
	}
}

func TestWebSocketResponseOriginOK(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/jmap/ws/", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := &http.Response{Request: req}
	if err := WebSocketResponseOriginOK("ws://127.0.0.1:8080/jmap/ws/", resp); err != nil {
		t.Fatalf("ws handshake origin: %v", err)
	}

	other, err := http.NewRequest(http.MethodGet, "http://evil.example/jmap/ws/", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := WebSocketResponseOriginOK("ws://127.0.0.1:8080/jmap/ws/", &http.Response{Request: other}); err == nil {
		t.Fatal("redirected handshake origin was accepted")
	}
}
