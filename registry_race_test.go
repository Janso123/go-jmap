package jmap_test

import (
	"sync"
	"testing"

	jsonv2 "encoding/json/v2"

	"github.com/Janso123/go-jmap"
	"github.com/stretchr/testify/require"
)

type raceCap struct{}

func (raceCap) URI() jmap.URI        { return "urn:ietf:params:jmap:testrace" }
func (raceCap) New() jmap.Capability { return &raceCap{} }

func TestRegistryConcurrent(t *testing.T) {
	jmap.RegisterMethod("Race/get", func() jmap.MethodResponse { return &struct{}{} })
	jmap.RegisterCapability(raceCap{})

	const n = 200
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		for i := 0; i < n; i++ {
			jmap.RegisterMethod("Race/get", func() jmap.MethodResponse { return &struct{}{} })
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < n; i++ {
			jmap.RegisterCapability(raceCap{})
		}
	}()
	go func() {
		defer wg.Done()
		raw := []byte(`{"methodResponses":[["Race/get",{},"c0"]],"sessionState":"s"}`)
		for i := 0; i < n; i++ {
			var resp jmap.Response
			require.NoError(t, jsonv2.Unmarshal(raw, &resp))
		}
	}()
	wg.Wait()
}
