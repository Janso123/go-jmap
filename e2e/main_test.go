//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/core/push/websocket"

	_ "github.com/Janso123/go-jmap/core"
)

const (
	baseURL     = "http://127.0.0.1:18080"
	composeProj = "go-jmap-e2e"
	reportPath  = "e2e/report.md"
)

var (
	rep     = &journal{image: "stalwartlabs/stalwart:v0.16"}
	started = time.Now().UTC()
)

func TestMain(m *testing.M) {
	// go test runs with the package directory as cwd. The compose file
	// and report path are relative to the module root.
	if err := os.Chdir(moduleRoot()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	code := 1
	defer func() {
		_ = rep.write(reportPath, started)
		fmt.Fprintf(os.Stderr, "e2e report: %s\n", reportPath)
		_ = compose("down", "-v")
		os.Exit(code)
	}()
	if _, err := exec.LookPath("docker"); err != nil {
		rep.setHarness("docker compose is required for e2e tests")
		fmt.Fprintln(os.Stderr, "docker compose is required for e2e tests")
		return
	}
	if err := compose("up", "-d", "--wait"); err != nil {
		rep.setHarness(err.Error())
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if err := waitReady(60 * time.Second); err != nil {
		rep.setHarness(err.Error())
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if err := provision(); err != nil {
		rep.setHarness(err.Error())
		fmt.Fprintln(os.Stderr, err)
		return
	}
	var err error
	alice, err = login("alice", "alice@example.org", "alice-e2e")
	if err != nil {
		rep.setHarness(err.Error())
		fmt.Fprintln(os.Stderr, err)
		return
	}
	bob, err = login("bob", "bob@example.org", "bob-e2e")
	if err != nil {
		rep.setHarness(err.Error())
		fmt.Fprintln(os.Stderr, err)
		return
	}
	code = m.Run()
}

func moduleRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	if filepath.Base(wd) == "e2e" {
		if _, err := os.Stat(filepath.Join(wd, "docker-compose.yml")); err == nil {
			return filepath.Dir(wd)
		}
	}
	return wd
}

func compose(args ...string) error {
	cmd := exec.Command("docker", append([]string{"compose", "-p", composeProj, "-f", "e2e/docker-compose.yml"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker compose %s: %w", args[0], err)
	}
	return nil
}

func waitReady(d time.Duration) error {
	deadline := time.Now().Add(d)
	c := &http.Client{Timeout: 2 * time.Second}
	var last error
	for time.Now().Before(deadline) {
		resp, err := c.Get(baseURL + "/healthz/ready")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
			last = fmt.Errorf("health %s", resp.Status)
		} else {
			last = err
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("healthz/ready: %w", last)
}

func provision() error {
	// alice-e2e is below the default zxcvbn score and bob-e2e is 7
	// characters. The Basic passwords stay those values; only the server
	// policy is lowered so create will accept them.
	if _, err := cli("update", "Authentication",
		"--field", "passwordMinStrength=zero",
		"--field", "passwordMinLength=1",
	); err != nil {
		return err
	}
	if _, err := cli("create", "Action", "--field", "@type=ReloadSettings"); err != nil {
		return err
	}
	if _, err := cli("create", "Domain", "--field", "name=example.org", "--field", "isEnabled=true"); err != nil {
		return err
	}
	out, err := cli("query", "Domain", "--where", "name=example.org", "--fields", "id")
	if err != nil {
		return err
	}
	id := lastToken(out)
	if id == "" || id == "id" {
		return fmt.Errorf("domain id missing in %q", out)
	}
	for _, name := range []string{"alice", "bob"} {
		if _, err := cli("create", "Account/User",
			"--field", "name="+name,
			"--field", "domainId="+id,
			"--field", `credentials={"0":{"@type":"Password","secret":"`+name+`-e2e"}}`,
			"--field", `roles={"@type":"User"}`,
			"--field", `permissions={"@type":"Inherit"}`,
			"--field", `encryptionAtRest={"@type":"Disabled"}`,
			"--field", `aliases={}`,
			"--field", `memberGroupIds={}`,
			"--field", `quotas={}`,
		); err != nil {
			return err
		}
	}
	return nil
}

func cli(args ...string) (string, error) {
	cmd := exec.Command("docker", append([]string{
		"compose", "-p", composeProj, "-f", "e2e/docker-compose.yml",
		"run", "--rm", "--no-deps", "cli",
	}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("stalwart-cli %s: %w\n%s", args[0], err, out)
	}
	return string(out), nil
}

func lastToken(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	if len(lines) == 0 {
		return ""
	}
	fields := strings.Fields(lines[len(lines)-1])
	if len(fields) == 0 {
		return ""
	}
	return fields[len(fields)-1]
}

func login(name, email, password string) (*account, error) {
	c := jmap.NewClient(baseURL+"/.well-known/jmap", jmap.WithBasic(email, password), jmap.WithTimeout(30*time.Second))
	if err := c.Authenticate(context.Background()); err != nil {
		return nil, fmt.Errorf("%s session: %w", name, err)
	}
	for _, raw := range []string{c.Session.APIURL, c.Session.UploadURL, c.Session.DownloadURL, c.Session.EventSourceURL} {
		if err := assertPublicURL(raw); err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
	}
	if cap, ok := c.Session.Capabilities[websocket.URI].(*websocket.WebSocket); ok {
		if err := assertPublicURL(cap.URL); err != nil {
			return nil, fmt.Errorf("%s websocket: %w", name, err)
		}
	}
	return &account{Name: name, Email: email, Password: password, Client: c}, nil
}

func assertPublicURL(raw string) error {
	if raw == "" {
		return nil
	}
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if u.Host != "127.0.0.1:18080" {
		return fmt.Errorf("%s host is %s", raw, u.Host)
	}
	if u.Scheme != "http" && u.Scheme != "ws" {
		return fmt.Errorf("%s scheme is %s", raw, u.Scheme)
	}
	return nil
}
