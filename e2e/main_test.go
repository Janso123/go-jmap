//go:build e2e

package e2e

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
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
