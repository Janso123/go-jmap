package jmap_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenericSetRejectsDestroyOnlyTypes(t *testing.T) {
	root, err := os.Getwd()
	require.NoError(t, err)

	for _, pkg := range []string{
		"./sharing/sharenotification/",
		"./calendar/calendareventnotification/",
	} {
		t.Run(pkg, func(t *testing.T) {
			cmd := exec.Command("go", "build", "-tags", "destroyonlyset", pkg)
			cmd.Dir = root
			out, err := cmd.CombinedOutput()
			require.Error(t, err)
			require.Contains(t, string(out), "JMAPCreatable")
		})
	}
}
