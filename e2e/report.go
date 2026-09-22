//go:build e2e

package e2e

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type kind string

const (
	kindPass   kind = "pass"
	kindFail   kind = "fail"
	kindSkip   kind = "skip"
	kindNotRun kind = "not-run"
)

type invocation struct {
	Scenario string
	RFC      string
	Method   string
	Account  string
	Request  string
	Response string
	Result   kind
}

var scenarioOrder = []string{"Session", "Mail", "Contacts", "Calendar", "Blob", "Quota", "Push"}

var scenarioLead = map[string]string{
	"Session":  "Alice and Bob fetch the JMAP session and call Core/echo.",
	"Mail":     "Alice reads Inbox, creates a message, and sends a second one to Bob.",
	"Contacts": "Alice creates an E2E address book with 30 cards, updates one, and copies a card to Bob.",
	"Calendar": "Alice creates an E2E calendar with 5 events, including one with Bob.",
	"Blob":     "Alice uploads text, attaches it to a message, and copies the blob to Bob.",
	"Quota":    "Alice reads quota.",
	"Push":     "Alice enables push and creates a message until a StateChange arrives.",
}

type journal struct {
	mu         sync.Mutex
	items      []invocation
	harnessErr string
	image      string
}

func (j *journal) add(in invocation) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.items = append(j.items, in)
}

func (j *journal) setHarness(err string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.harnessErr = err
}

func (j *journal) write(path string, started time.Time) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	var b strings.Builder
	counts := collapsedCounts(j.items)
	image := j.image
	if image == "" {
		image = "stalwartlabs/stalwart:v0.16"
	}
	fmt.Fprintf(&b, "# go-jmap e2e\n\n")
	fmt.Fprintf(&b, "- image: %s\n", image)
	fmt.Fprintf(&b, "- start: %s\n", started.UTC().Format(time.RFC3339))
	fmt.Fprintf(&b, "- elapsed: %s\n", time.Since(started).Round(time.Second))
	fmt.Fprintf(&b, "- pass: %d\n- fail: %d\n- skip: %d\n- not-run: %d\n\n",
		counts[kindPass], counts[kindFail], counts[kindSkip], counts[kindNotRun])
	if j.harnessErr != "" {
		fmt.Fprintf(&b, "Harness stopped before scenarios: %s\n", j.harnessErr)
		return os.WriteFile(path, []byte(b.String()), 0o644)
	}
	for _, name := range scenarioOrder {
		var group []invocation
		for _, in := range j.items {
			if in.Scenario == name {
				group = append(group, in)
			}
		}
		if len(group) == 0 {
			continue
		}
		fmt.Fprintf(&b, "## %s\n\n%s\n\n", name, scenarioLead[name])
		type bucket struct {
			method string
			items  []invocation
		}
		var buckets []bucket
		index := map[string]int{}
		for _, in := range group {
			i, ok := index[in.Method]
			if !ok {
				index[in.Method] = len(buckets)
				buckets = append(buckets, bucket{method: in.Method})
				i = len(buckets) - 1
			}
			buckets[i].items = append(buckets[i].items, in)
		}
		for _, bucket := range buckets {
			fmt.Fprintf(&b, "### %s\n\n", bucket.method)
			for n, in := range bucket.items {
				fmt.Fprintf(&b, "#### Invocation %d\n\n", n+1)
				fmt.Fprintf(&b, "- RFC: %s\n- account: %s\n- request: %s\n- response: %s\n- result: %s\n\n",
					in.RFC, in.Account, in.Request, in.Response, in.Result)
			}
		}
	}
	fmt.Fprintf(&b, "## Index\n\n| RFC | Method | Scenario | Result |\n|-----|--------|----------|--------|\n")
	for _, name := range scenarioOrder {
		seen := map[string]kind{}
		var order []string
		rfcOf := map[string]string{}
		for _, in := range j.items {
			if in.Scenario != name {
				continue
			}
			if _, ok := seen[in.Method]; !ok {
				order = append(order, in.Method)
				seen[in.Method] = in.Result
				rfcOf[in.Method] = in.RFC
				continue
			}
			seen[in.Method] = mergeKind(seen[in.Method], in.Result)
		}
		for _, method := range order {
			fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", rfcOf[method], method, name, seen[method])
		}
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func collapsedCounts(items []invocation) map[kind]int {
	counts := map[kind]int{}
	for _, name := range scenarioOrder {
		seen := map[string]kind{}
		for _, in := range items {
			if in.Scenario != name {
				continue
			}
			prev, ok := seen[in.Method]
			if !ok {
				seen[in.Method] = in.Result
				continue
			}
			seen[in.Method] = mergeKind(prev, in.Result)
		}
		for _, result := range seen {
			counts[result]++
		}
	}
	return counts
}

func mergeKind(a, b kind) kind {
	if a == kindFail || b == kindFail {
		return kindFail
	}
	if a == kindPass || b == kindPass {
		return kindPass
	}
	if a == kindSkip && b == kindSkip {
		return kindSkip
	}
	if a == kindNotRun && b == kindNotRun {
		return kindNotRun
	}
	if a == b {
		return a
	}
	return kindPass
}
