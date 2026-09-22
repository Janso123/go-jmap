//go:build e2e

package e2e

import (
	"strconv"
	"testing"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/quota"
)

func TestQuota(t *testing.T) {
	sc := &scenario{name: "Quota"}
	if alice == nil || alice.Client == nil {
		t.Fatal("account session missing")
	}
	if !hasCap(alice.Client, quota.URI) {
		skipRest(t, sc, quotaSteps, "missing "+string(quota.URI))
		return
	}

	aliceID, err := alice.Client.PrimaryAccount(quota.URI)
	if err != nil {
		skipRest(t, sc, quotaSteps, err.Error())
		return
	}

	call[*quota.GetResponse](t, sc, step{
		RFC: "RFC 9425", Method: "Quota/get", Account: alice.Name, Request: "ids omitted",
	}, alice.Client, []jmap.URI{quota.URI}, false, &quota.Get{Account: aliceID}, func(resp *quota.GetResponse) (string, error) {
		n := 0
		if resp != nil {
			n = len(resp.List)
		}
		return "len=" + strconv.Itoa(n), nil
	})

	call[*quota.QueryResponse](t, sc, step{
		RFC: "RFC 9425", Method: "Quota/query", Account: alice.Name, Request: "no filter",
	}, alice.Client, []jmap.URI{quota.URI}, false, &quota.Query{Account: aliceID}, func(resp *quota.QueryResponse) (string, error) {
		n := 0
		if resp != nil {
			n = len(resp.IDs)
		}
		return "len=" + strconv.Itoa(n), nil
	})

	if sc.failed {
		t.Fail()
	}
}

var quotaSteps = []step{
	{RFC: "RFC 9425", Method: "Quota/get", Account: "alice", Request: "ids omitted"},
	{RFC: "RFC 9425", Method: "Quota/query", Account: "alice", Request: "no filter"},
}
