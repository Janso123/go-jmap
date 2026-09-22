//go:build e2e

package e2e

import (
	"net/http"
	"testing"
	"time"
)

func TestHarnessHealth(t *testing.T) {
	c := &http.Client{Timeout: 5 * time.Second}
	resp, err := c.Get(baseURL + "/healthz/ready")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %s", resp.Status)
	}
}
