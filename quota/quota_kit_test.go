package quota_test

import (
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/quota"
	"github.com/stretchr/testify/require"
)

func TestQuotaGetName(t *testing.T) {
	require.Equal(t, "Quota/get", (&quota.Get{}).Name())
}

func TestQuotaChangesName(t *testing.T) {
	require.Equal(t, "Quota/changes", (&quota.Changes{}).Name())
}

func TestQuotaQueryName(t *testing.T) {
	require.Equal(t, "Quota/query", (&quota.Query{}).Name())
}

func TestQuotaQueryChangesName(t *testing.T) {
	require.Equal(t, "Quota/queryChanges", (&quota.QueryChanges{}).Name())
}

func TestQuotaRequiresURI(t *testing.T) {
	require.Equal(t, []jmap.URI{quota.URI}, (&quota.Get{}).Requires())
}
