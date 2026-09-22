//go:build e2e

package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestReportMarkdown(t *testing.T) {
	j := &journal{}
	j.add(invocation{Scenario: "Poczta", RFC: "RFC 8621", Method: "Email/query", Account: "alice", Request: "inMailbox inbox", Response: "ids=a", Result: kindPass})
	j.add(invocation{Scenario: "Poczta", RFC: "RFC 8621", Method: "Email/query", Account: "bob", Request: "subject=hi", Response: "timeout", Result: kindFail})
	j.add(invocation{Scenario: "Sesja", RFC: "RFC 8620", Method: "Core/echo", Account: "alice", Request: "ping=e2e", Response: "ping=e2e", Result: kindPass})

	dir := t.TempDir()
	path := filepath.Join(dir, "report.md")
	started := time.Date(2026, 9, 22, 7, 0, 0, 0, time.UTC)
	if err := j.write(path, started); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	sesja := strings.Index(text, "## Sesja")
	poczta := strings.Index(text, "## Poczta")
	if sesja < 0 || poczta < 0 || sesja > poczta {
		t.Fatalf("scenario order:\n%s", text)
	}
	if !strings.Contains(text, "pass: 1") || !strings.Contains(text, "fail: 1") {
		t.Fatalf("counts:\n%s", text)
	}
	if !strings.Contains(text, "| RFC 8621 | Email/query | Poczta | fail |") {
		t.Fatalf("index row:\n%s", text)
	}
	if !strings.Contains(text, "### Email/query") || !strings.Contains(text, "#### Wywołanie 1") || !strings.Contains(text, "#### Wywołanie 2") {
		t.Fatalf("invocations:\n%s", text)
	}
}
